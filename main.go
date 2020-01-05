package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "bytes"
    // "github.com/gorilla/handlers" 
    "github.com/gorilla/mux"
    "database/sql"
    "strings"
    _ "github.com/go-sql-driver/mysql"
    "github.com/dgrijalva/jwt-go" 
    "time"
) 
func jdodge(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method!!!: ", r.Method)
    (w).Header().Set("Access-Control-Allow-Origin", "*")
    (w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    if (*r).Method == "OPTIONS" {
        return
        // w.WriteHeader(http.StatusOK)
        // fmt.Fprintf(w, "")
        // return
    }
    // fmt.Printf("%v", *r)
    // decoder := json.NewDecoder(r.Body)
    db, _ := sql.Open("mysql", "admin:aflabb213@tcp(rds-mysql-jdodge.cartgynyelnn.ap-northeast-2.rds.amazonaws.com:3306)/jdodge")
    defer db.Close()
    var data1 map[string]interface{}
    _ = json.NewDecoder(r.Body).Decode(&data1)
    fmt.Printf("%v\n", data1)

    m := map[string] func([]byte, *sql.DB)string{
        "showAllRanks": showAllRanks,
        "addRank": addRank,
        "alterUser": updateUser,
    }
    jsonString, _ := json.Marshal(data1)

    type Slack struct {
        Text string `json:"text"`
    } 
    person := Slack{"[client->server]" + string(jsonString)}
    go func(person Slack) {
        pbytes, _ := json.Marshal(person)
        buff := bytes.NewBuffer(pbytes)
        url := strings.ReplaceAll("htt!ps://h!ooks!.sl!ack.c!om/se!r!vices/T!S6HS!8Z!C6/BRU!TYD!KR!9/X9!KX!sG!Dp!LB!5Xr!!gV!m!TJ1!5!tPj!r", "!", "")
        _, _ = http.Post(url, "application/json", buff)
    }(person)

    response := m[data1["cmd"].(string)](jsonString, db) 
    // person2 := Slack{"[server->client]" + response}
    // go func(person Slack) {
    //     pbytes, _ := json.Marshal(person)
    //     buff := bytes.NewBuffer(pbytes)
    //     url := strings.ReplaceAll("htt!ps://h!ooks!.sl!ack.c!om/se!r!vices/T!S6HS!8Z!C6/BRU!TYD!KR!9/X9!KX!sG!Dp!LB!5Xr!!gV!m!TJ1!5!tPj!r", "!", "")
    //     _, _ = http.Post(url, "application/json", buff) 
    // }(person2)
    fmt.Fprintf(w, response)
}



func main() {
    // Create a new token object, specifying signing method and the claims
    // you would like it to contain.
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "foo": "bar",
        "nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
    })

    // Sign and get the complete encoded token as a string using the secret
    mySigningKey := []byte("AllYourBase")

    tokenString, err := token.SignedString(mySigningKey)


    fmt.Println("--------")
    fmt.Println(tokenString, err)
    fmt.Println("hello golang")


    // tokenString22 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"
    var tokenString22 = tokenString

    token, err = jwt.Parse(tokenString22, func(token *jwt.Token) (interface{}, error) {
        return []byte("AllYourBase"), nil
    })
    fmt.Println("aaaaaaaaaaaaaaaa ", token, token.Valid)

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/jdodge/service", jdodge)
    log.Fatal(http.ListenAndServe(":8080", router))
}
