package business

import(
    "database/sql"
    "strings"
    "fmt"
)

func UpdateRow(db *sql.DB, r Rows, i int64)(int64, error){
    table := strings.TrimSpace(r.TableName)
    user := strings.TrimSpace(r.UserName)
    
    re := Re()
    
    if !ValidIdent(table, re){
        return 0, fmt.Errorf("invalid table name")
    }
    
    if !ValidIdent(user, re){
        return 0, fmt.Errorf("invalid user name")
    }
    
    columns, err := FetchColumns(db, user+"_"+table)
    if err != nil{
        return 0, err
    }
    
    for _, col := range columns{
        if col.Nullable != "YES"{
            if col.Name != "t_id" && col.Name != "created_at"{
                for _, data := range r.Data{
                    if data.ColumnName == col.Name{
                        if strings.TrimSpace(data.Value) == ""{
                            return 0, fmt.Errorf("empty field for %s", data.ColumnName)
                        }
                    }
                }
               
            }
            
        }
        
        
    }
    for _, dat := range r.Data{
        avalaible := false
        for _, colu := range columns{
            if colu.Name == dat.ColumnName{
                avalaible = true
                if !ValidType(colu.DataType, dat.Value){
                    return 0, fmt.Errorf("invalid datatype %s of column %s ",dat.Value, dat.ColumnName)
                }
            }
                
        }
        if dat.ColumnName == "t_id" || dat.ColumnName == "created_at"{
            return 0, fmt.Errorf("this field %s is not required", dat.ColumnName)
        }
        
        if !avalaible{
            return 0, fmt.Errorf("this field %s is out of range", dat.ColumnName)
        }
    }
    
    collst := make([]string, 0)
    
    for _, col := range columns{
        for _, dat := range r.Data{
            if col.Name == dat.ColumnName{
                collst = append(collst, fmt.Sprintf("%s=%s",QuoteIdent(col.Name), QuoteLiteral(dat.Value)))
            }
        }
    }
    
    colstr := strings.Join(collst, ", ")
    
    query := fmt.Sprintf(`UPDATE %s SET %s WHERE t_id=$1`, QuoteIdent(user+"_"+table), colstr)
    
    tx, err := db.Begin()
    if err != nil{
        return 0, err
    }
    
    
    res, err := tx.Exec(query, i)
    if err != nil{
        tx.Rollback()
        return 0, err
    }
    
    rowsAffected, err := res.RowsAffected()
    if err != nil{
        tx.Rollback()
        return 0, err
    }
    
    if rowsAffected < 1{
        tx.Rollback()
        return 0, fmt.Errorf("row data with id %d not found",i)
    }
    
    if rowsAffected > 1{
        tx.Rollback()
        return 0, fmt.Errorf("operation unsuccessful!")
    }
    tx.Commit()
    return i, nil
}