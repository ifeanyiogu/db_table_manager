package business

import (
    "database/sql"
    "strings"
    "fmt"
)
func DeleteTable(db *sql.DB, t string, u string)error{
    re := Re()
    table := strings.TrimSpace(t)
    user := strings.TrimSpace(u)
    
    if table == "" || user == ""{
        return fmt.Errorf("username and table are required")
    }
    
    if !ValidIdent(table, re) || !ValidIdent(user, re){
        return fmt.Errorf("Invalid table_name")
    }
    
    table_name := QuoteIdent(user+"_"+table)
    
    query := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, table_name)
    
    tx, err := db.Begin()
    if err != nil{
        return err
    }
    defer tx.Rollback()
    
    _, err = tx.Exec(query)
    if err != nil{
        return err
    }
    
    if err := tx.Commit(); err != nil{
        return err
    }
    return nil
    
}