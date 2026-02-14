package business

import (
    "regexp"
    "database/sql"
    "strings"
    "fmt"
    "context"
)

func Re() *regexp.Regexp{
    return regexp.MustCompile(`^[a-z]+(?:_[a-z]+)*$`)
}

type Table struct{
    Name string `json:"name"`
    Columns []Column `json:"columns"`
}

type Column struct{
    Name string `json:"name"`
    DataType string `json:"data_type"`
    Nullable string `json:"nullable"`
    Default any `json:"default"`
}

func (t *Table)CreateTable(db *sql.DB, ctx context.Context, user string)error{
    re := Re()
    name := strings.TrimSpace(t.Name)
    if name == ""{
        return fmt.Errorf("Empty table name")
    }
    
    if len(name) > 30{
        return fmt.Errorf("Table name too long")
    }
    
    if !ValidIdent(name, re){
        return fmt.Errorf("Invalid table name Format")
    }
    
    name = user +"_"+name
    columns := t.Columns
    if len(columns) < 1{
        return fmt.Errorf("Empty table Columns")
    }
    
    var count int64
    if err := db.QueryRow(`SELECT count(table_name) FROM information_schema.tables WHERE table_schema = 'public' AND table_name LIKE $1`, user+"_%").Scan(&count); err != nil{
        if err != sql.ErrNoRows{
            return err
        }
    }
    fmt.Println(count)
    if count >= 2{
        return fmt.Errorf("You've Reached Your limit: %d tables of 2 allowed", count)
    }
    formated_columns := make([]string, 0)
    
    formated_columns = append(formated_columns, `t_id BIGSERIAL NOT NULL PRIMARY KEY`)
    
    for _, column := range columns{
        col_name := strings.TrimSpace(column.Name)
        if col_name == ""{
            return fmt.Errorf("Empty column name")
        }
        
        if col_name == "t_id" || col_name == "created_at"{
            return fmt.Errorf("column name %s is a reserved name", col_name)
        }
        
        if len(col_name) > 20{
            return fmt.Errorf("column name %s too long", col_name)
        }
        
        if !ValidIdent(col_name, re){
            return fmt.Errorf("Invalid column name format %s", col_name)
        }
        
        data_type := strings.TrimSpace(column.DataType)
        if data_type == ""{
            return fmt.Errorf("no data_type provided for column name %s", col_name)
        }
        ok, main_type := CheckType(data_type)
        if !ok{
            return fmt.Errorf("data_type %s provided for column name %s is Invalid",data_type,  col_name)
        }
        
        col_lst := make([]string,0)
        col_lst = append(col_lst, QuoteIdent(col_name))
        col_lst = append(col_lst, main_type)
        if strings.TrimSpace(column.Nullable) == "NO"{
            col_lst = append(col_lst, "NOT NULL")
        }
        
        col_str := strings.Join(col_lst, " ")
        formated_columns = append(formated_columns, col_str)
        
    }
    formated_columns = append(formated_columns, `created_at TIMESTAMPTZ DEFAULT NOW()`)
    formated_columns_str := strings.Join(formated_columns, ", ")
    
    query_str := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (%s)`, QuoteIdent(name), formated_columns_str)
    
    _, err := db.ExecContext(ctx, query_str)
    return err
}