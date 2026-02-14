package users

import(
    "database/sql"
    "strings"
    "fmt"
    "errors"
    "context"
    
    "go-tb/business"
)


func Login(u *User, db *sql.DB, ctx context.Context)(ResponseUser, error){
    username := strings.TrimSpace(u.Username)
    password := strings.TrimSpace(u.Password)
    
    if username == "" || password == ""{
        return ResponseUser{}, fmt.Errorf("username and password are required")
    }
    
    re := business.Re()
    
    if !business.ValidIdent(username, re){
        return ResponseUser{}, fmt.Errorf("invalid username")
    }
    query := `SELECT password, id FROM users WHERE username = $1`
    var hashed_password string
    var id int64
    err := db.QueryRowContext(ctx, query, username).Scan(&hashed_password, &id)
    if err != nil{
        if errors.Is(err, sql.ErrNoRows){
            return ResponseUser{}, fmt.Errorf("Invalid username and password", username)
        }
        fmt.Println(err)
        return ResponseUser{}, fmt.Errorf("something went wrong")
    }
    
    err = CompareHashPassword([]byte(hashed_password), []byte(password))
    if err != nil{
        fmt.Println(err)
        return ResponseUser{}, fmt.Errorf("Invalid credentials")
    }
    
    
    token, err := GenerateToken(id)
    if err != nil{
        return ResponseUser{}, err
    }
    
    
    return ResponseUser{Username:username, Token: token, ID:id}, nil
}
