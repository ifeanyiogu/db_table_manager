package business

import(
    "database/sql"
    "strings"
    "fmt"
)

type Rows struct{
    UserName string `json:"username"`
    TableName string `json:"table_name"`
    Data []Data `json:"data"`
}

type Data struct{
    ColumnName string `json:"name"`
    Value string `json:"value"`
}

func CreateRow(db *sql.DB, r Rows)(int64, error){
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
                found := false
                for _, data := range r.Data{
                    if data.ColumnName == col.Name{
                        if strings.TrimSpace(data.Value) == ""{
                            return 0, fmt.Errorf("empty field for %s", data.ColumnName)
                        }
                        found = true
                    }
                }
                if !found{
                    return 0, fmt.Errorf("non nullable field for %s", col.Name)
                }
            }
            
        }
        
        
    }
    for _, dat := range r.Data{
        avalaible := false
        for _, colu := range columns{
            if colu.Name == dat.ColumnName{
                avalaible = true
            }
        }
        if dat.ColumnName == "t_id" || dat.ColumnName == "created_at"{
            return 0, fmt.Errorf("this field %s is not required", dat.ColumnName)
        }
        
        if !avalaible{
            return 0, fmt.Errorf("this field %s is out of range", dat.ColumnName)
        }
    }
    
    for _, dat := range r.Data{
        for _, col := range columns{
            if dat.ColumnName == col.Name{
                if !ValidType(col.DataType, dat.Value){
                    return 0, fmt.Errorf("invalid datatype %s of column %s ",dat.Value, dat.ColumnName)
                }
            }
        }
    }
    
    collst := make([]string, 0)
    vallst := make([]string, 0)
    
    for _, col := range columns{
        for _, dat := range r.Data{
            if col.Name == dat.ColumnName{
                collst = append(collst, QuoteIdent(col.Name))
                vallst = append(vallst, QuoteLiteral(dat.Value))
            }
        }
    }
    
    colstr := strings.Join(collst, ", ")
    valstr := strings.Join(vallst, ", ")
    
    query := fmt.Sprintf(`INSERT INTO %s(%s) VALUES (%s) RETURNING t_id`, QuoteIdent(user+"_"+table), colstr, valstr)
    
    var id int64 
    if err := db.QueryRow(query).Scan(&id); err != nil{
        return 0, err
    }
    
    return id, nil
}