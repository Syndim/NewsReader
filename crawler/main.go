package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "nc/workers"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Missing db path")
        os.Exit(0)
    }
    db, _ := sql.Open("mysql", os.Args[1])
    defer db.Close()
    cnbeta := workers.NewCnbeta()
    err := cnbeta.Start(db)
    if err != nil {
        fmt.Println(err)
    }
}
