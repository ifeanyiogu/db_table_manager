package business

import(
    "regexp"
    "strings"
    "database/sql"
    "strconv"
    "fmt"
    "time"
    
    "github.com/shopspring/decimal"
)


func ValidIdent(s string, re *regexp.Regexp)bool{
    if len(s) > 20{
        return false
    }
    return re.MatchString(s)
}

func CheckType(s string)(bool, string){
    types := map[string]string{
        "int": "INTEGER",
        "bigint": "BIGINT",
        "text":"TEXT",
        "decimal": "NUMERIC",
        "bool":"BOOLEAN",
        "datetime":"TIMESTAMPTZ",
    }
    for key, val := range types{
        if key == s{
            return true, val
        }
    }
    return false, ""
}

func QuoteIdent(s string)string{
    return `"`+strings.ReplaceAll(s, `"`, `""`)+`"`
}

func QuoteLiteral(s string)string{
    return `'`+strings.ReplaceAll(s, `'`, `''`)+`'`
}


func FetchColumns(db *sql.DB, t string)([]Column, error){
    query := `SELECT column_name, data_type, is_nullable, column_default FROM information_schema.columns WHERE table_name = $1 ORDER BY ordinal_position`
    
    rows, err := db.Query(query, t)
    if err != nil{
        return nil, err
    }
    defer rows.Close()
    
    columns := make([]Column, 0)
    for rows.Next(){
        column := Column{}
        if err := rows.Scan(&column.Name, &column.DataType, &column.Nullable, &column.Default); err != nil{
            return nil, err
        }
        
        columns = append(columns, column)
    }
    
    return columns, nil
    
}

func ValidType(t string, v string)bool{
    switch t{
        case "integer":
            _, err := strconv.Atoi(v)
            if err != nil{
                return false
            }
        case "boolean":
           _, err := strconv.ParseBool(v)
            if err != nil{
                return false
            }
        case "numeric":
            _, err := decimal.NewFromString(v)
            if err != nil{
                return false
            }
        case "bigint":
            _, err := strconv.Atoi(v)
            if err != nil{
                return false
            }
        case "text":
            return true
            
        case "timestamptz", "timestamp with time zone", "timestamp":
            _, err := time.Parse(time.RFC3339, v)
            if err != nil{
                return false
            }
        default:
            return false
    }
    return true
}


func parseValue(v any, c string)string{
    switch c{
        case "integer", "bigint":
            switch val := v.(type){
                case int, int64, int32:
                    return fmt.Sprintf("%v",val)
            }
        case "numeric":
            switch val := v.(type){
                case float32,float64:
                    return fmt.Sprintf("%f",val)
                case []byte:
                    return string(val)
            }
        case "boolean":
            if val, ok := v.(bool);ok{
                return fmt.Sprintf("%t", val)
            }
        case "text":
            if val, ok := v.(string);ok{
                return string(val)
            }
        default:
            switch val := v.(type){
                default:
                    return fmt.Sprintf("%v",val)
            }
    }
    return ""
}