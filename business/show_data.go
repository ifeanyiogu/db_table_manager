package business

import (
    "database/sql"
    "strings"
    "errors"
    "fmt"
)


func ShowData(db *sql.DB, id int64, t, u string)(Row, error){
    name := strings.TrimSpace(t)
    user := strings.TrimSpace(u)
    
    re := Re()
    if !ValidIdent(t, re){
        return Row{}, fmt.Errorf("invalid table name")
    }
    
    if !ValidIdent(user, re){
        return Row{}, fmt.Errorf("invalid user name")
    }
    
    columns, err := FetchColumns(db, user+"_"+name)
    if err != nil{
       return Row{}, err 
    }
    
    row := Row{}
    
    colType := make(map[string]string)
    for _, c := range columns {
        colType[c.Name] = c.DataType
    }
    
    collst := make([]string, 0)
    
    for _, col := range columns{
        collst = append(collst, QuoteIdent(col.Name))
    }
    
    colstr := strings.Join(collst, ", ")
    query := fmt.Sprintf(`SELECT %s FROM %s WHERE t_id = $1`, colstr, QuoteIdent(user+"_"+name))
    
    ptr := make([]any, len(columns))
    vals := make([]any, len(columns))
        
    for i, _ := range vals{
        ptr[i] = &vals[i]
    }
    
    err = db.QueryRow(query, id).Scan(ptr...)
    if err != nil{
        if errors.Is(err, sql.ErrNoRows){
            return Row{}, fmt.Errorf("No Data for id %d in table %s", id, name)
        }
        return Row{}, err
    }

    for i, v := range vals{
        var data DataR
        if columns[i].Name == "t_id"{
            data.ColumnName = "id"
        }else{
            data.ColumnName = columns[i].Name
        }
        data.Type = fmt.Sprintf("%T", v)
        if v == nil{
            data.Value = "NULL"
        }else{
        data.Value = parseValue(v, colType[columns[i].Name])
        }
        row.Data = append(row.Data, data)
    }
    
    return row, nil
}
