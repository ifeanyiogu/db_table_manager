package handlers

import(
    "database/sql"
    
    "github.com/gofiber/fiber/v2"
    
    "go-tb/users"
)

func CreateUserHandler(db *sql.DB)fiber.Handler{
    return func(c *fiber.Ctx) error{
        var user users.User
        err := c.BodyParser(&user)
        if err != nil{
            return c.Status(400).JSON(fiber.Map{"error":err.Error()})
        }
        res, err := user.Register(db, c.Context())
        if err != nil{
            return c.Status(400).JSON(fiber.Map{"error":err.Error()})
        }
        return c.Status(201).JSON(fiber.Map{"token": res.Token, "name":res.Username})
    }
}