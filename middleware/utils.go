package middleware

import(
    "errors"
    "strings"
    
    "github.com/golang-jwt/jwt/v5"
    
    "go-tb/users"
)

var SecreteKey []byte = users.SecreteKey

func ParseToken(s string)(int64, error){
    token, err := jwt.Parse(s, func(token *jwt.Token) (any, error){
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
            return nil, errors.New("Invalid Token Method")
        }
        return SecreteKey, nil
        
    })
    if err != nil{
        return 0, err
    }
    
    if !token.Valid{
        return 0, errors.New("Invalid token")
    }
    
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok{
        return 0, errors.New("Invalid claims")
    }
    
    id, ok := claims["user"].(float64)
    if !ok{
        return 0, errors.New("user id missing")
    }
    
    return int64(id), nil
}


func ScrapeHeader(s string)(string, error){
    parts := strings.Fields(s)
    if len(parts) != 2{
        return "",  errors.New("invalid authorization header")
    }
    if parts[0] != "Bearer"{
        return "",  errors.New("invalid authorization Type")
    }
    return parts[1], nil
}