package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/berkanturkmen/VatanSoft/database"
	"github.com/berkanturkmen/VatanSoft/handlers"
	"github.com/berkanturkmen/VatanSoft/middleware"
	"github.com/berkanturkmen/VatanSoft/cache"
)

func main() {
	app := fiber.New()

	database.ConnectDB()

	cache.ConnectRedis()

	api := app.Group("/api")

	api.Post("/register", handlers.Register)
	api.Post("/send-reset-code", handlers.GetResetCode)
	api.Post("/reset-password", handlers.ResetPassword)
	api.Post("/login", handlers.Login)

	secureApi := api.Group("/secure")
	secureApi.Use(middleware.SecureMiddleware())
	secureApi.Get("/employee-current", handlers.GetCurrentEmployee)
	secureApi.Post("/logout", handlers.Logout)

	secureApi.Get("/employees", middleware.PermissionMiddleware("Employee"), handlers.ListEmployees)
	secureApi.Post("/employees", middleware.PermissionMiddleware("Owner"), handlers.CreateEmployee)
	secureApi.Get("/employees/:id", middleware.PermissionMiddleware("Employee"), handlers.GetEmployee)
	secureApi.Put("/employees/:id", middleware.PermissionMiddleware("Owner"), handlers.UpdateEmployee)
	secureApi.Delete("/employees/:id", middleware.PermissionMiddleware("Owner"), handlers.DeleteEmployee)

	secureApi.Get("/polyclinics", middleware.PermissionMiddleware("Employee"), handlers.GetAllPolyclinics)
	secureApi.Get("/polyclinics/mine", middleware.PermissionMiddleware("Employee"), handlers.GetPolyclinicsByHospital)
	secureApi.Post("/attach-polyclinic/:polyclinicID", middleware.PermissionMiddleware("Owner"), handlers.AttachPolyclinic)
	secureApi.Post("/detach-polyclinic/:polyclinicID", middleware.PermissionMiddleware("Owner"), handlers.DetachPolyclinic)

	secureApi.Get("/personnels", middleware.PermissionMiddleware("Employee"), handlers.ListPersonnels)
	secureApi.Post("/personnels", middleware.PermissionMiddleware("Owner"), handlers.CreatePersonnel)
	secureApi.Get("/personnels/:id", middleware.PermissionMiddleware("Employee"), handlers.GetPersonnel)
	secureApi.Put("/personnels/:id", middleware.PermissionMiddleware("Owner"), handlers.UpdatePersonnel)
	secureApi.Delete("/personnels/:id", middleware.PermissionMiddleware("Owner"), handlers.DeletePersonnel)

	secureApi.Get("/cities", middleware.PermissionMiddleware("Employee"), handlers.GetCities)
	secureApi.Get("/cities/:cityID", middleware.PermissionMiddleware("Employee"), handlers.GetRegionsByCity)

	secureApi.Get("/jobs", middleware.PermissionMiddleware("Employee"), handlers.GetJobs)
	secureApi.Get("/jobs/:jobID", middleware.PermissionMiddleware("Employee"), handlers.GetTitlesByJob)

	app.Listen(":8080")
}
