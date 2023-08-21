package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/renushrii/finance-tracking/repository"
	"golang.org/x/crypto/bcrypt"
)

// we can rename this to auth.go instead of login.go
// because we are doing both login and signup

// todo: use secret key from config.yaml
const secretKey = "1234"

func (s *Server) Logout(c *fiber.Ctx) error {
	// get token cookie
	token := c.Cookies("token")
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now(),
	})

	return c.Redirect("/p/home")
}

func (s *Server) Login(c *fiber.Ctx) error {
	email, password := c.FormValue("email"), c.FormValue("password")
	token, err := s.login(email, password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).Redirect("/p/auth")
	}

	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: token,
	})

	return c.Redirect("/p/home")
}

func getUserID(c *fiber.Ctx) (int, error) {
	user := c.Locals("user")
	if user == nil {
		return 0, errors.New("user not found")
	}

	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	return int(claims["userId"].(float64)), nil
}

func (s *Server) login(email, password string) (string, error) {
	user, err := s.users.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if !compareHash(user.Password, password) {
		return "", errors.New("email or password is incorrect")
	}

	claims := jwt.MapClaims{
		"email":  email,
		"userId": user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secretKey))
	return t, err
}

func (s *Server) Signin(c *fiber.Ctx) error {
	firstname, lastname, email, password := c.FormValue("firstname"), c.FormValue("lastname"), c.FormValue("email"), c.FormValue("password")

	user := repository.User{
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		Password:  password,
	}

	// validate if the user has entered all the required fields
	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	hashPassword, err := hash(user.Password)
	if err != nil {
		if err != nil {
			return err
		}
		log.Println(err)
		return c.Redirect("/p/auth")
	}

	// we want to store the hashed password in the database
	user.Password = hashPassword

	_, err = s.users.FindByEmail(user.Email)
	if err == nil {
		return c.Redirect("/p/auth")
	}

	err = s.users.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Redirect("/p/auth")
	}

	token, err := s.login(email, password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Redirect("/p/auth")
	}

	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: token,
	})

	return c.Redirect("/p/home")
}

func hash(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func compareHash(hashed, s string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(s)) == nil
}
