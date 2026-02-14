package handlers

import(
    "database/sql"
    
    "github.com/gofiber/fiber/v2"
    
    "go-tb/users"
    "go-tb/business"
)

func ShowColumnsHandler(db *sql.DB) fiber.Handler{
    return func(c *fiber.Ctx)error{
        user := c.Locals("user").(users.User)
        if err := user.GetUser(db); err != nil{
            return c.Status(404).JSON(fiber.Map{"error":err.Error()})
        }
        re := business.Re()
        table := c.Params("table")
        
        if !business.ValidIdent(table, re){
            return c.Status(404).JSON(fiber.Map{"error":"invalid table name"})
        }
        
        data, err := business.FetchColumns(db, user.Username+"_"+table)
        if err != nil{
            return c.Status(400).JSON(fiber.Map{"error":err.Error()})
        }
        
        return c.Status(200).JSON(fiber.Map{"data": data})
    }
}