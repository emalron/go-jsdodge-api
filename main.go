package main
 
import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "database/sql" 
    "reflect"
    _ "github.com/go-sql-driver/mysql"
)


var db *sql.DB
 
func updateUser(jsonString []byte) string{
    type UpdateUserCommand struct {
        Cmd string
        Id string
        Name  string
    }

    updateUserCommand := UpdateUserCommand{}
    if err := json.Unmarshal(jsonString, &updateUserCommand); err != nil {
        // do error check
        fmt.Println(err)
    }
    fmt.Printf("%v\n", updateUserCommand)
    return "updateUser"
}
 
func showAllRanks(jsonString []byte) string{ 
    var name string
    var score int
    var replayData string
    fmt.Println(jsonString)
    rows, err := db.Query("SELECT name, score, replay_data FROM view_ranking")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    containerJson := make(map[string]interface{}, 0)
    containerJson["result"] = 0
    rankData := make([]map[string]interface{}, 0)
    for rows.Next() {
        err := rows.Scan(&name, &score, &replayData)
        if err != nil {
            log.Fatal(err)
        }
        var m = map[string]interface{}{
            "name":name,
            "score":score,
            "replay_data":replayData,
        }
        rankData = append(rankData, m)
    }
    containerJson["message"] = rankData
    jsonBytes, err := json.Marshal(containerJson)
    if err != nil {
        panic(err)
    }

    // JSON 바이트를 문자열로 변경
    retString := string(jsonBytes)
    fmt.Println(retString)
    return retString
}
func addRank(jsonString []byte) string{
 
    type AddRankCommand struct {
        Cmd string
        Id string
        Score  int
    }
    rankCommand := AddRankCommand{}
    if err := json.Unmarshal(jsonString, &rankCommand); err != nil {
        // do error check
        fmt.Println(err)
    }
    fmt.Printf("%v\n", rankCommand)
    return "addRank"
}
func homeLink(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome home!")
}
func jdodge(w http.ResponseWriter, r *http.Request) {
    // fmt.Printf("%v", *r)
    // decoder := json.NewDecoder(r.Body)
    var data1 map[string]interface{}
    _ = json.NewDecoder(r.Body).Decode(&data1)
    fmt.Printf("%v\n", data1)

    m := map[string] func([]byte)string{
        "showAllRanks": showAllRanks,
        "addRank": addRank,
        "alterUser": updateUser,
    }
    jsonString, _ := json.Marshal(m)
    response := m[data1["cmd"].(string)](jsonString)

    // fmt.Fprintf(w, "Welcome home!")
    fmt.Fprintf(w, response)
}

func main() {
    db, _ = sql.Open("mysql", "admin:aflabb213@tcp(rds-mysql-jdodge.cartgynyelnn.ap-northeast-2.rds.amazonaws.com:3306)/jdodge")
    fmt.Println(reflect.TypeOf(db))


    // if err != nil {
    //     log.Fatal(err)
    // }
    defer db.Close()

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", homeLink)
    router.HandleFunc("/jdodge/service", jdodge)
    log.Fatal(http.ListenAndServe(":8080", router))
}
