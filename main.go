package main
 
import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "bytes"
    "github.com/gorilla/mux"
    "database/sql" 
    "strings"
    _ "github.com/go-sql-driver/mysql"
) 

 

func jdodge(w http.ResponseWriter, r *http.Request) {
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
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/jdodge/service", jdodge)
    log.Fatal(http.ListenAndServe(":8080", router))
}
