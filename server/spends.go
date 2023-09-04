package server

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/renushrii/finance-tracking/repository"
)

func (s *Server) GetSpends(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	from := c.Query("from")
	to := c.Query("to")
	dateFormat := "2006-01-02"

	fromTime, err := time.Parse(dateFormat, from)
	if err != nil {
		return err
	}

	toTime, err := time.Parse(dateFormat, to)
	if err != nil {
		return err
	}

	spends, err := s.spends.Get(userID, fromTime, toTime)
	if err != nil {
		return err
	}

	return c.JSON(spends)
}

func (s *Server) CreateSpend(c *fiber.Ctx) error {
	description := c.FormValue("description")
	tag := c.FormValue("tag")
	amountstr := c.FormValue("amount")

	userId, err := getUserID(c)
	if err != nil {
		return err
	}

	amount, err := strconv.Atoi(amountstr)
	if err != nil {
		return err
	}

	spend := repository.Spend{
		Description: description,
		Tag:         tag,
		Amount:      amount,
		UserID:      userId,
	}

	_, err = s.spends.Create(spend)
	if err != nil {
		return err
	}
	return c.Redirect("/p/home?msg=spend created successfully")

}

func (s *Server) GetSpendByTag(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	tag := c.Params("tag")
	spends, err := s.spends.GetByTag(userID, tag)
	if err != nil {
		return err
	}
	return c.JSON(spends)
}

func (s *Server) UpdateSpendById(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	id := c.Params("id")
	ID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	var spend repository.Spend
	if err := c.BodyParser(&spend); err != nil {
		return err
	}

	spend.Id = ID
	spend.UserID = userID

	if err := s.spends.Update(spend); err != nil {
		return err
	}
	return c.Redirect("/p/home?msg=spend updated successfully")
}

func (s *Server) DeleteSpendById(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	id := c.Params("id")
	ID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if err := s.spends.Delete(userID, ID); err != nil {
		return err
	}
	return c.Redirect("/p/home?msg=spend deleted successfully")
}

func (s *Server) PercentageByTag(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	tag := c.Params("tag")
	spends, err := s.spends.GetByTag(userID, tag)
	if err != nil {
		return err
	}

	totalAmount := 0
	for _, spend := range spends {
		totalAmount += spend.Amount
	}

	var calculatedPercentage float64
	if totalAmount > 0 {
		calculatedPercentage = float64(totalAmount) / float64(len(spends))
	}

	return c.JSON(fiber.Map{
		"tag":        tag,
		"percentage": calculatedPercentage,
	})
}
