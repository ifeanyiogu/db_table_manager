package db

import(
    "database/sql"
    "fmt"
    "os"
    
    _"github.com/lib/pq"
)

func Config() *sql.DB{
    db_setting := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`, os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
    db, err := sql.Open("postgres", db_setting)
    if err != nil{
        db.Close()
        panic(err)
    }
    
    if err = db.Ping(); err != nil{
        db.Close()
        panic(err)
    }
    
    return db
}