package main

import (
    "fmt"
    "log"
    "encoding/json"
    "database/sql" 
    _ "github.com/go-sql-driver/mysql"
) 
func showAllRanks(jsonString []byte, db *sql.DB) string{ 
    var name string
    var score int
    var replayData string
    var time1 string
    fmt.Println(jsonString)
    rows, err := db.Query("SELECT name, score, replay_data, time FROM view_ranking")
    if err != nil {
        log.Fatal("ksoo error", err)
    }
    defer rows.Close()

    containerJson := make(map[string]interface{}, 0)
    containerJson["result"] = 0
    rankData := make([]map[string]interface{}, 0)
    for rows.Next() {
        err := rows.Scan(&name, &score, &replayData, &time1)
        if err != nil {
            // log.Fatal(err)
            time1 = ""
        }
        var m = map[string]interface{}{
            "name":name,
            "score":score,
            "replay_data":replayData,
            "time": time1,
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
func updateUser(jsonString []byte, db *sql.DB) string{
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
func addRank(jsonString []byte, db *sql.DB) string{
 
    // to do: 아이디 있는지 확인 해야함.
    type Command struct {
        Cmd string
        Id string
        Score  int
        Replay_data string
    }
    command := Command{}
    if err := json.Unmarshal(jsonString, &command); err != nil {
        // do error check
        fmt.Println(err)
    }
    result, err := db.Exec("INSERT INTO ranks(score, replay_data, users_id) VALUES (?, ?, ?)", command.Score, command.Replay_data, command.Id)
    if err != nil {
        log.Fatal(err)
    }
 
    // sql.Result.RowsAffected() 체크
    n, err := result.RowsAffected()
    if n == 1 {
        fmt.Println("1 row inserted.")
    }
    fmt.Printf("%v\n", command)
    return "addRank"
}
