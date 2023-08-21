package server

import (
	"fmt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/renushrii/finance-tracking/repository"
)

type Server struct {
	spends repository.Spends
	users  repository.Users
	app    *fiber.App

	// todo: add secret key here
}

// New returns a new Server instance
func New(spends repository.Spends, users repository.Users) *Server {
	engine := html.New("./templates", ".html")

	s := &Server{app: fiber.New(
		fiber.Config{
			Views: engine,
		},
	), spends: spends, users: users}
	s.register()

	return s
}

func (s *Server) register() {
	app := s.app

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} Body: ${green}${body}${white}\n",
	}))

	app.Use("/p/home", jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(secretKey)},
		TokenLookup: "cookie:token",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/p/auth")
		},
	}))

	app.Use("/api", jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(secretKey)},
		TokenLookup: "cookie:token",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/p/auth")
		},
	}))

	app.Get("/p/spends/create", jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(secretKey)},
		TokenLookup: "cookie:token",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/p/auth")
		},
	}), s.addSpend)

	// static files
	app.Static("/s", "./static")

	// authentication
	app.Post("/auth/login", s.Login)
	app.Post("/auth/signin", s.Signin)
	app.Get("/auth/logout", s.Logout)
	// app.Post("/auth/logout", s.Logout)

	// pages
	app.Get("/p/auth", s.Auth)
	app.Get("/p/home", s.Home)

	// /uri    ---> function (mapping)

	// apis
	app.Get("/api/spends/get", s.GetSpends)
	app.Post("/api/spends/create", s.CreateSpend)

	app.Get("/api/spends/get/:tag", s.GetSpendByTag)
	app.Get("/api/spends/percentage/:tag", s.PercentageByTag)

	app.Post("/api/spends/update/:id", s.UpdateSpendById)
	app.Get("/api/spends/delete/:id", s.DeleteSpendById)
}

// Start starts the server
func (s *Server) Start(host string, port int) error {
	return s.app.Listen(fmt.Sprintf("%s:%d", host, port))
}
