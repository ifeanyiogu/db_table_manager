package handlers

import(
    "database/sql"
    
    "github.com/gofiber/fiber/v2"
    
    "go-tb/business"
    "go-tb/users"
)

func CreateTableHandler(db *sql.DB)fiber.Handler{
    return func(c *fiber.Ctx) error{
        var table business.Table 
        
        err := c.BodyParser(&table)
        if err != nil{
            return c.Status(400).JSON(fiber.Map{"error":err.Error()})
        }
        ctx := c.Context()
        
        user := c.Locals("user").(users.User)
        if err := user.GetUser(db); err != nil{
            return c.Status(404).JSON(fiber.Map{"error":err.Error()})
        }
        if err := table.CreateTable(db, ctx, user.Username); err != nil{
            return c.Status(404).JSON(fiber.Map{"error":err.Error()})
        }
        
        return c.Status(201).JSON(fiber.Map{"success":"created"})
    }
}