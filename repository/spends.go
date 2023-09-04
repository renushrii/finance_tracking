package repository

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrSpendDoesNotBelongToUser = errors.New("spend does not belong to user")

type Spends interface {
	// Get spends which were created between from and to
	Get(userID int, from, to time.Time) ([]Spend, error)
	// Create a new spend, for a user_id
	Create(spend Spend) (int, error)
	// GetByTag returns spends which have the given tag
	GetByTag(userID int, tag string) ([]Spend, error)

	TagsAmount(userId int, from, to time.Time) ([]TagsAmount, error)

	Delete(userID, id int) error
	Update(spend Spend) error
}

type TagsAmount struct {
	Tag    string `json:"tag"`
	Amount int    `json:"amount"`
}

type Spend struct {
	Id          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	Amount      int       `db:"amount" json:"amount"`
	Description string    `db:"description" json:"description"`
	Tag         string    `db:"tag" json:"tag"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type spends struct {
	db *sqlx.DB
}

func NewSpends(db *sqlx.DB) Spends {
	return &spends{
		db: db,
	}
}

func (s *spends) Get(userID int, from, to time.Time) ([]Spend, error) {
	var spends []Spend
	query := `
		SELECT * FROM spends
		WHERE created_at >= ? AND created_at <= ? AND user_id = ?
	`

	err := s.db.Select(&spends, query, from, to, userID)
	return spends, err
}

// JWT, MySQL, APIs, Database Design, Golang, fiber framework, ...

func (s *spends) Create(spend Spend) (int, error) {
	query := `
		INSERT INTO spends (user_id, amount, description, tag, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, spend.UserID, spend.Amount, spend.Description, spend.Tag, time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func (s *spends) GetByTag(userID int, tag string) ([]Spend, error) {
	var spends []Spend
	query := `
		SELECT * FROM spends
		WHERE tag = ? AND user_id = ?
	`
	err := s.db.Select(&spends, query, tag, userID)
	return spends, err
}

func (s *spends) TagsAmount(userID int, from, to time.Time) ([]TagsAmount, error) {
	var spends []TagsAmount
	query := `
		SELECT tag, sum(amount) AS amount FROM spends
		WHERE user_id = ?
		AND created_at >= ? AND created_at <= ?
		GROUP BY tag
	`
	err := s.db.Select(&spends, query, userID, from, to)
	return spends, err
}

func (s *spends) Delete(userID, id int) error {
	if !s.belongsTo(userID, id) {
		return ErrSpendDoesNotBelongToUser
	}

	query := `
		DELETE FROM spends
		WHERE id = ?
	`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *spends) belongsTo(userID int, spendID int) bool {
	query := `
		SELECT COUNT(*) FROM spends
		WHERE id = ? AND user_id = ?
	`
	var count int
	s.db.Get(&count, query, spendID, userID)
	return count > 0
}

func (s *spends) Update(spend Spend) error {
	if !s.belongsTo(spend.UserID, spend.Id) {
		return ErrSpendDoesNotBelongToUser
	}

	query := `UPDATE spends SET amount = :amount, description = :description, tag = :tag WHERE id = :id AND user_id = :user_id`
	_, err := s.db.NamedExec(query, spend)
	return err
}
