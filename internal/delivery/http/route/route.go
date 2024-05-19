package route

import (
	"crowdfunding-service/internal/delivery/http"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RouteConfig struct {
	App            *fiber.App
	Log            *logrus.Logger
	UserController *http.UserController
}

func (c *RouteConfig) Setup() {
	c.App.Use(c.recoverPanic)
	c.SetupGuestRoute()
}

func (c *RouteConfig) recoverPanic(ctx *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic occurred: %v", r)
			c.Log.WithError(err).Error("Panic occured")
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
	}()

	return ctx.Next()
}

func (c *RouteConfig) SetupGuestRoute() {
	GuestGroup := c.App.Group("/api")

	GuestGroup.Get("/users", c.UserController.List)
	GuestGroup.Post("/users", c.UserController.RegisterUser)
	GuestGroup.Put("/users/:userId", c.UserController.Update)
	GuestGroup.Get("/users/:userId", c.UserController.Get)
	GuestGroup.Delete("/users/:userId", c.UserController.Delete)
}
