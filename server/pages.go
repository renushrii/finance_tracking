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

	tagsAmount, err := s.spends.TagsAmount(userID, from, to)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "application/text")

	return c.Render("home", fiber.Map{
		"Spends":     spends,
		"TagsAmount": tagsAmount,
	})
}

func (s *Server) addSpend(c *fiber.Ctx) error {
	return c.Render("add_spend", fiber.Map{})
}

func (s *Server) updatespend(c *fiber.Ctx) error {
	return c.Render("update_spend", fiber.Map{})
}
