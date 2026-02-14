package users

import(
    "fmt"
    "time"
    "os"
    
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "github.com/joho/godotenv"
)



func LoadEnv(key string)string{
    err := godotenv.Load()
    if err != nil{
        panic("env failed")
    }
    
    val, exists := os.LookupEnv(key)
    if !exists{
        panic("Could not find key")
    }
    return val
}
var SecreteKey []byte = []byte(LoadEnv("SECRET_KEY"))

func HashPassword(s string)(string, error){
    hash,err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
    if err != nil{
        fmt.Println(err)
        return "", fmt.Errorf("Something went wrong")
    }
    
    return string(hash), nil
}


func GenerateToken(id int64)(string, error){
    token:= jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user": id,
        "iat": time.Now().Unix(),
        "exp": time.Now().Add(time.Hour*24).Unix(),
    })
    token_string, err := token.SignedString(SecreteKey)
    if err != nil{
        fmt.Println(err)
        return "", fmt.Errorf("Something went wrong")
    }
    
    return token_string, nil
}

func CompareHashPassword(hashed []byte, pass []byte)error{
    err := bcrypt.CompareHashAndPassword(hashed, pass)
    return err
}