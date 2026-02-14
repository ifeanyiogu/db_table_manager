package handlers

import(
    "database/sql"
    "strings"
    "strconv"
    
    "github.com/gofiber/fiber/v2"
    
    "go-tb/users"
    "go-tb/business"
)

func ShowDataHandler(db *sql.DB) fiber.Handler{
    return func(c *fiber.Ctx)error{
        user := c.Locals("user").(users.User)
        if err := user.GetUser(db); err != nil{
            return c.Status(404).JSON(fiber.Map{"error":err.Error()})
        }
        table := c.Params("table")
        idstr := strings.TrimSpace(c.Params("id"))
        id, err := strconv.ParseInt(idstr, 10, 64)
        if err != nil{
            return c.Status(400).JSON(fiber.Map{"error":"Invalid Row ID"})
        }
        
        data, err := business.ShowData(db, id, table, user.Username)
        if err != nil{
            return c.Status(400).JSON(fiber.Map{"error":err.Error()})
        }
        
        return c.Status(200).JSON(data)
    }
}