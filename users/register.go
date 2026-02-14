
package users

import(
    "database/sql"
    "strings"
    "fmt"
    "errors"
    "context"
    
    "github.com/lib/pq"
    
    "go-tb/business"
)

type User struct{
    ID int64 `json:"id"`
    Token string `json"token"`
    Username string `json:"username"`
    Password string `json:"password"`
}
type ResponseUser struct{
    Username string `json:"username"`
    Token string `json:"token"`
    ID int64 `json:"id"`
}

func (u *User)Register(db *sql.DB, ctx context.Context)(ResponseUser, error){
    username := strings.TrimSpace(u.Username)
    password := strings.TrimSpace(u.Password)
    
    if username == "" || password == ""{
        return ResponseUser{}, fmt.Errorf("username and password are required")
    }
    
    re := business.Re()
    
    if !business.ValidIdent(username, re){
        return ResponseUser{}, fmt.Errorf("invalid username")
    }
    
    hashed_password, err := HashPassword(password)
    if err != nil{
        fmt.Println(err)
        return ResponseUser{}, fmt.Errorf("something went wrong try again")
    }
    
    
    query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
    var id int64
    err = db.QueryRowContext(ctx, query, username, hashed_password).Scan(&id)
    if err != nil{
        var pgErr *pq.Error
        if errors.As(err, &pgErr) && pgErr.Code== "23505"{
            return ResponseUser{}, fmt.Errorf("User with username %s already exist", username)
        }
        fmt.Println(err)
        return ResponseUser{}, fmt.Errorf("something went wrong")
    }
    
    token, err := GenerateToken(id)
    if err != nil{
        return ResponseUser{}, err
    }
    
    
    return ResponseUser{Username:username, Token: token, ID:id}, nil
}

func (u *User)GetUser(db *sql.DB) error{
    query := `SELECT username FROM users WHERE id = $1`
    err := db.QueryRow(query, u.ID).Scan(&u.Username)
    if err != nil{
        if errors.Is(err, sql.ErrNoRows){
            return fmt.Errorf("user with ID %d not found", u.ID)
        }
        return err
    }
    return nil
}