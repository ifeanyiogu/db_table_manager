package middleware


import(
    "database/sql"
    "strings"
    
    "github.com/gofiber/fiber/v2"
    
    "go-tb/users"
)

func AuthUser(db *sql.DB)fiber.Handler{
    return func(c *fiber.Ctx)error{
        token, err := ScrapeHeader(strings.TrimSpace(c.Get("Authorization")))
        if err != nil{
            return c.Status(400).JSON(fiber.Map{"error":err.Error()})
        }
        id, err := ParseToken(token)
        if err != nil{
            return c.Status(400).JSON(fiber.Map{"error":"login failed"})
        }
        
        user := users.User{
            ID: id,
            Token: token,
        }
        
        c.Locals("user", user)
        return c.Next()
    }
}