package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "usfonapr6fnyadxz:Wm4Sp3BCi01PPIwyMM7U@(b09440an79d4csniwsfn-mysql.services.clever-cloud.com:3306)/b09440an79d4csniwsfn")

	if err != nil {
		log.Panic(err)
	}

	r, err := db.Exec("CREATE DATABASE users")

	if err != nil {
		log.Panic(err)
	}

	fmt.Print(r)

}

// helpful-quanta-378421:southamerica-east1:diego-segutierrez4
