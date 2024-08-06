package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/berkanturkmen/VatanSoft/database"
	"github.com/berkanturkmen/VatanSoft/models"
	"time"
)

func SecureMiddleware() fiber.Handler {
	return func(c * fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		tokenString := ""
		if len(authHeader) > 7 && authHeader[: 7] == "Bearer " {
			tokenString = authHeader[7: ]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		var session models.Session
		if err := database.DB.Where("token = ?", tokenString).First( & session).Error;err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		if time.Now().After(session.ExpirationTime) {
			if err := database.DB.Delete( & session).Error;
			err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
					"Message": "",
					"Success": false,
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		var employee models.Employee
		if err := database.DB.First( & employee, session.EmployeeID).Error;err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		c.Locals("EmployeeID", session.EmployeeID)
		c.Locals("HospitalID", employee.HospitalID)
		c.Locals("Permission", string(employee.Permission))

		return c.Next()
	}
}

func PermissionMiddleware(required string) fiber.Handler {
	return func(c * fiber.Ctx) error {
		permission, ok := c.Locals("Permission").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		if (required == "Employee" && (permission == "Owner" || permission == "Employee")) || (required == "Owner" && permission == "Owner") {
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}
}
