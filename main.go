package main

import (
    "fmt"
    "os"
    
    "github.com/joho/godotenv"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    
    "go-tb/handlers"
    "go-tb/middleware"
    "go-tb/db"
)

func main(){
    if err := godotenv.Load(); err != nil{
        panic("Unable to env file")
    }
    db := db.Connect()
    defer db.Close()
    app := fiber.New()
    
    corsConfig := cors.Config{
        AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
        AllowHeaders: "Authorization, Content-Type",
        AllowCredentials: true,
        AllowMethods: "POST, GET, PUT, DELETE",
    }
    
    app.Use(cors.New(corsConfig))
    
    app.Use(func (c *fiber.Ctx)error{
        fmt.Println(c.Method(), c.Path())
        return c.Next()
    })
    
    app.Post("/users/reg", handlers.CreateUserHandler(db))
    app.Post("/users/login", handlers.LoginUserHandler(db)) 
    
    app.Post("/tables/create", middleware.AuthUser(db), handlers.CreateTableHandler(db))
    app.Post("/tables/drop/:table", middleware.AuthUser(db), handlers.DeleteTableHandler(db))
    
    app.Post("/rows", middleware.AuthUser(db), handlers.CreateRowHandler(db))
    app.Put("/rows/:id", middleware.AuthUser(db), handlers.UpdateTableHandler(db))
    app.Delete("/rows/:table/:id", middleware.AuthUser(db), handlers.DeleteRowHandler(db))
    
    
    app.Get("/tables", middleware.AuthUser(db), handlers.ListTableHandler(db))
    app.Get("/tables/show/:table", middleware.AuthUser(db), handlers.ShowTableHandler(db))
    app.Get("/columns/get/:table", middleware.AuthUser(db), handlers.ShowColumnsHandler(db))
    app.Get("/data/:table/:id", middleware.AuthUser(db), handlers.ShowDataHandler(db))
    
    port, ok := os.LookupEnv("PORT")
    if !ok{
        panic("invalud port")
    }
    
    app.Listen(":"+port)
}