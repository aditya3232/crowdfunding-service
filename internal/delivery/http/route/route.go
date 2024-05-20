package route

import (
	"crowdfunding-service/internal/delivery/http"
	"crowdfunding-service/internal/delivery/http/middleware"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RouteConfig struct {
	App                     *fiber.App
	Log                     *logrus.Logger
	Oauth2Middleware        *middleware.Oauth2Middleware
	Oauth2Controller        *http.Oauth2Controller
	UserController          *http.UserController
	CampaignController      *http.CampaignController
	CampaignImageController *http.CampaignImageController
}

func (c *RouteConfig) Setup() {
	c.App.Use(c.recoverPanic)
	c.SetupGuestRoute()
	c.SetupAuthRoute()
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

	GuestGroup.Post("/users", c.UserController.RegisterUser)

	GuestGroup.Get("/google/login", c.Oauth2Controller.GoogleLogin)
	GuestGroup.Get("/google/callback", c.Oauth2Controller.GoogleCallback)
	GuestGroup.Post("/google/revoke", c.Oauth2Controller.RevokeToken)
	GuestGroup.Post("/google/refresh", c.Oauth2Controller.RefreshToken)
}

func (c *RouteConfig) SetupAuthRoute() {
	AuthGroup := c.App.Group("/api")
	AuthGroup.Use(c.Oauth2Middleware.Oauth2Middleware)

	AuthGroup.Get("/users", c.UserController.List)
	AuthGroup.Get("/users/me", c.UserController.CurrentUser)
	AuthGroup.Put("/users/:userId", c.UserController.Update)
	AuthGroup.Put("/users/upload/avatar", c.UserController.UpdateAvatar)
	AuthGroup.Get("/users/:userId", c.UserController.Get)
	AuthGroup.Delete("/users/:userId", c.UserController.Delete)

	AuthGroup.Post("/campaigns", c.CampaignController.CreateCampaign)
	AuthGroup.Get("/campaigns", c.CampaignController.List)
	AuthGroup.Get("/campaigns/:campaignId", c.CampaignController.Get)
	AuthGroup.Put("/campaigns/:campaignId", c.CampaignController.Update)

	AuthGroup.Post("/campaigns/upload/image", c.CampaignImageController.Create)

}
