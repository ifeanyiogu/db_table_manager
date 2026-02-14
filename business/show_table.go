package business

import(
    "strings"
    "fmt"
    "database/sql"
)
type TableData struct{
    Name string `json:"name"`
    Columns []string `json:"columns"`
    Rows []Row `json:"rows"`
}

type Row struct{
    Data []DataR `json:"data"`
}

type DataR struct{
    ColumnName string `json"name"`
    Value any `json:"value"`
    Type string `json:"type"`
}

func ShowTable(db *sql.DB, t, u string)(TableData, error){
    name := strings.TrimSpace(t)
    user := strings.TrimSpace(u)
    
    re := Re()
    if !ValidIdent(t, re){
        return TableData{}, fmt.Errorf("invalid table name")
    }
    
    if !ValidIdent(user, re){
        return TableData{}, fmt.Errorf("invalid user name")
    }
    
    columns, err := FetchColumns(db, user+"_"+name)
    if err != nil{
       return TableData{}, err 
    }
    
    table_data := TableData{
        Name: name,
    }
    
    colType := make(map[string]string)
    for _, c := range columns {
        colType[c.Name] = c.DataType
    }
    
    collst := make([]string, 0)
    
    for _, col := range columns{
        if col.Name == "t_id"{
            table_data.Columns = append(table_data.Columns, "id")
        }else{
            table_data.Columns = append(table_data.Columns, col.Name)
        }
        collst = append(collst, QuoteIdent(col.Name))
    }
    
    colstr := strings.Join(collst, ", ")
    query := fmt.Sprintf(`SELECT %s FROM %s ORDER BY t_id`, colstr, QuoteIdent(user+"_"+name))
    
    rows, err := db.Query(query)
    if err != nil{
        return TableData{}, err
    }
    defer rows.Close()
    
    for rows.Next(){
        ptr := make([]any, len(columns))
        vals := make([]any, len(columns))
        
        for i, _ := range vals{
            ptr[i] = &vals[i]
        }
        
        if err := rows.Scan(ptr...); err != nil{
            return TableData{}, err
        }
        
        var row Row
        
        for i, v := range vals{
            var data DataR
            data.ColumnName = table_data.Columns[i]
            data.Type = fmt.Sprintf("%T", v)
            if v == nil{
                data.Value = "NULL"
            }else{
            data.Value = parseValue(v, colType[table_data.Columns[i]])
            }
            row.Data = append(row.Data, data)
        }
        
        table_data.Rows = append(table_data.Rows, row)
    }
    return table_data, nil
}
