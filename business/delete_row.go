package business


import (
    "database/sql"
    "strings"
    "fmt"
)

func DeleteRow(db *sql.DB,t string, u string, id int64)error{
    table := strings.TrimSpace(t)
    user := strings.TrimSpace(u)
    if table == "" || user == ""{
        return fmt.Errorf("Table name is missing")
    }
    
    re := Re()
    if !ValidIdent(table, re) || !ValidIdent(user, re){
        return fmt.Errorf("invalid table name")
    }
    
    if id < 1{
        return fmt.Errorf("Invalid Row id")
    }
    
    table_name := QuoteIdent(user+"_"+table)
    query := fmt.Sprintf(`DELETE FROM %s WHERE t_id = $1`, table_name)
    
    tx, err := db.Begin()
    if err != nil{
        return err
    }
    
    defer tx.Rollback()
    
    res, err := tx.Exec(query, id)
    if err != nil{
        return err
    }
    
    rA, err := res.RowsAffected()
    if rA != 1{
        if rA > 1{
            return fmt.Errorf("Error: Something went wrong")
        }
        return fmt.Errorf("Data with Row ID: %d not found")
    }
    
    if err := tx.Commit(); err != nil{
        return fmt.Errorf("Error: operation Uncompleted")
    }
    return nil
}