package server

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// keeping all the pages in one place
func (s *Server) Auth(c *fiber.Ctx) error {
	return c.Render("auth", fiber.Map{})
}

func (s *Server) Home(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		fmt.Print(err)
		return c.Status(fiber.StatusInternalServerError).Redirect("/p/auth")
	}

	to := time.Now()
	from := to.AddDate(0, -1, 0)

	spends, err := s.spends.Get(userID, from, to)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "application/text")
	return c.Render("home", fiber.Map{
		"Spends": spends,
		// "SpendsByTag": <....>, []Struct{Tag, Amount}
	})
}

func (s *Server) addSpend(c *fiber.Ctx) error {
	return c.Render("add_spend", fiber.Map{})
}

func (s *Server) updatespend(c *fiber.Ctx) error {
	return c.Render("update_spend", fiber.Map{})
}

// todo: create page for /p/spends/create page

// todo: create page for /p/spends/update page
// todo: create html template for /p/spends/update page

/*
/p/auth - login/signin page
	login - /auth/login - POST with FormValues (on success redirect to /p/home)
	signin - /auth/signin - POST with FormValues (on success redirect to /p/home)

/p/home - home page (from & to are set to now and now-1month)
	chart
	table by tag
	all spends table  - (edit button) --> /p/spends/update?id={spendID}

	button: create spend --> hyperlink to /p/spends/create    <a href="/p/spends/create"> <button> </a>

/p/spends/create - create spend page
	this page has a form with input fields
		- description
		- amount
		- date
		- tag
	submit button - POST /api/spends/create - (on success redirect to /p/home)

/p/spends/update?id=spendID - update spend page
	this page has a form with input fields
		- description       (pre-filled)
		- amount
		- date
		- tag

		<form action="/api/spends/update/{{ .ID }}" method="post">
			<input name="description" type="textarea" value="{{ .Description }}">
			<input name="amount" type="number" value="{{ .Amount }}">
	submit button - POST /api/spends/update/{id} - (on success redirect to /p/home)
*/
