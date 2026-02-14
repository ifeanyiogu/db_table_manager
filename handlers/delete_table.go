package handlers 


import (
    "database/sql"
    
    "github.com/gofiber/fiber/v2"
    
    "go-tb/business"
    "go-tb/users"
)

func DeleteTableHandler(db *sql.DB)fiber.Handler{
    return func(c *fiber.Ctx)error{
        user := c.Locals("user").(users.User)
        if err := user.GetUser(db); err != nil{
            return c.Status(400).JSON(fiber.Map{"error": err.Error()})
        }
        
        table := c.Params("table")
        if err := business.DeleteTable(db, table, user.Username); err != nil{
            return c.Status(400).JSON(fiber.Map{"error": err.Error()})
        }
        
        return c.Status(200).JSON(fiber.Map{"message":"deleted"})
    }
}