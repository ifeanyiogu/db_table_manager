package db

import(
    "database/sql"
)
func Connect() *sql.DB{
    db := Config()
    
    query := `CREATE TABLE IF NOT EXISTS users (id BIGSERIAL NOT NULL PRIMARY KEY, username TEXT NOT NULL UNIQUE, password TEXT NOT NULL, created_at TIMESTAMPTZ DEFAULT NOW())`
    
    _, err := db.Exec(query)
    if err != nil{
        db.Close()
        panic(err)
    }
    
    return db
} 