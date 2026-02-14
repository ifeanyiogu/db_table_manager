package business

import(
    "database/sql"
    "fmt"
)

type ColumnData struct{
    TableName string
    ColumnName string
    DataType string
    Nullable string
    Default any
}

func ListTables(db *sql.DB, name string) ([]Table, error){
    arg := name+"%"
    
    query := fmt.Sprintf(`SELECT t.table_name, c.column_name, c.data_type, c.is_nullable, c.column_default FROM information_schema.tables t LEFT JOIN information_schema.columns c ON t.table_name = c.table_name WHERE t.table_schema = 'public' AND t.table_type = 'BASE TABLE' AND t.table_name LIKE $1`)
    rows, err := db.Query(query, arg)
    if err != nil{
        return nil, err
    }
    defer rows.Close()
    
    tableMap := make(map[string][]ColumnData)
    
    for rows.Next(){
        var tbn sql.NullString
        var column ColumnData
        if err := rows.Scan(&tbn, &column.ColumnName, &column.DataType, &column.Nullable, &column.Default); err != nil{
            return nil, err
        }
        
        if tbn.Valid{
            tableMap[tbn.String] = append(tableMap[tbn.String], column)
        }
    }
    tables := make([]Table, 0)
    for key, val := range tableMap{
        table := Table{
            Name: key,
        }
        for _, col := range val{
            colum := Column{
                Name: col.ColumnName,
                DataType: col.DataType,
                Nullable: col.Nullable,
                Default: col.Default,
            }
            table.Columns = append(table.Columns, colum)
        }
        tables = append(tables, table)
    }
    return tables, nil
    
}