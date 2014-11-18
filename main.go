package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	// "newsreader/crawler/data/console"
	"newsreader/crawler/workers"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing db path")
		os.Exit(0)
	}
	db, _ := sql.Open("mysql", os.Args[1])
	defer db.Close()
	cnbeta := workers.NewCnbeta(workers.NewNewsReadWriter(db), workers.NewCommentReadWriter(db))
	// cnbeta := workers.NewCnbeta(new(console.ConsoleReadWriter), new(console.ConsoleReadWriter))
	err := cnbeta.Start()
	if err != nil {
		fmt.Println(err)
	}
}
