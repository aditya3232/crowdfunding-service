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
	CorsMiddleware          *middleware.CorsMiddleware
	Oauth2Controller        *http.Oauth2Controller
	UserController          *http.UserController
	CampaignController      *http.CampaignController
	CampaignImageController *http.CampaignImageController
	TransactionController   *http.TransactionController
}

func (c *RouteConfig) Setup() {
	c.App.Use(c.CorsMiddleware.CorsMiddleware())
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

	GuestGroup.Get("/campaigns", c.CampaignController.List)
	GuestGroup.Get("/campaigns/:campaignId", c.CampaignController.Get)

	// endpoint ini yg akan digunakan midtrans mengirim notification status pembayaran
	// dipengaturan midtrans, meunuju ke settings, lalu payment, lalu edit notifcation URL
	GuestGroup.Post("/transactions/notification", c.TransactionController.CreateTransactionNotification)
}

func (c *RouteConfig) SetupAuthRoute() {
	AuthGroup := c.App.Group("/api")
	AuthGroup.Use(c.Oauth2Middleware.Oauth2Middleware)

	AuthGroup.Get("/users", c.UserController.List)
	AuthGroup.Get("/users/me", c.UserController.CurrentUser)
	AuthGroup.Put("/users/:userId", c.UserController.Update)
	AuthGroup.Put("/users/avatar/upload", c.UserController.UpdateAvatar)
	AuthGroup.Get("/users/:userId", c.UserController.Get)

	AuthGroup.Post("/campaigns", c.CampaignController.CreateCampaign)
	AuthGroup.Put("/campaigns/:campaignId", c.CampaignController.Update)
	AuthGroup.Post("/campaigns/image/upload", c.CampaignImageController.Create)

	AuthGroup.Get("/transactions", c.TransactionController.GetTransactionsByUserID) // ambil semua transaksi dari user yg login
	AuthGroup.Post("/transactions", c.TransactionController.CreateTransaction)
	AuthGroup.Get("/transactions/campaigns/:campaignId", c.TransactionController.GetTransactionsByCampaignID) // ambil semua campaign dari pemiliknya (filter user_id wajib diisi)

}
