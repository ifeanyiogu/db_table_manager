package handlers


import(
    "database/sql"
    "strings"
    "strconv"
    
    "github.com/gofiber/fiber/v2"
    
    "go-tb/business"
    "go-tb/users"
)
func UpdateTableHandler(db *sql.DB) fiber.Handler{
    return func(c *fiber.Ctx) error{
        var rows business.Rows 
        if err := c.BodyParser(&rows); err != nil{
            return c.Status(400).JSON(fiber.Map{"error":err.Error()})
        }
        user := c.Locals("user").(users.User)
        if err := user.GetUser(db);err != nil{
            return c.Status(404).JSON(fiber.Map{"error":err.Error()})
        }
        rows.UserName = user.Username
        idstr := strings.TrimSpace(c.Params("id"))
        
        idint, err := strconv.ParseInt(idstr, 10, 64)
        if err != nil{
            return c.Status(400).JSON(fiber.Map{"error":"invalid row id"})
        }
        
        id, err := business.UpdateRow(db, rows, idint)
        if err != nil{
            return c.Status(400).JSON(fiber.Map{"error":err.Error()})
        }
        
        return c.Status(201).JSON(fiber.Map{"Message": "updated", "id":id})
    }
}