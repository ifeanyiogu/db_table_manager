package handlers


import(
    "database/sql"
    
    "github.com/gofiber/fiber/v2"
    
    "go-tb/business"
    "go-tb/users"
)
func CreateRowHandler(db *sql.DB) fiber.Handler{
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
        
        id, err := business.CreateRow(db, rows)
        if err != nil{
            return c.Status(400).JSON(fiber.Map{"error":err.Error()})
        }
        
        return c.Status(201).JSON(fiber.Map{"Message": "created", "id":id})
    }
}