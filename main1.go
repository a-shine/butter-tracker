package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "docker"
	dbname   = "world"
)

func ConnectToDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createSchema(db *sql.DB) {
	const drop = `DROP TABLE IF EXISTS nodes;`
	_, err := db.Exec(drop)

	if err != nil {
		log.Fatal(err)
	}

	const create string = `
  		CREATE TABLE IF NOT EXISTS nodes (
  		addr VARCHAR NOT NULL PRIMARY KEY,
  		peers VARCHAR NOT NULL,
  		groups VARCHAR,
  		last_seen TIMESTAMP NOT NULL DEFAULT NOW()
  		);`

	_, err = db.Exec(create)

}

func main() {
	db := ConnectToDB()
	defer db.Close()
	createSchema(db)

	for {
		time.Sleep(time.Second * 10)
		//	Remove all the nodes that haven't been seen in a while
		const delete string = `
		DELETE FROM nodes
		WHERE last_seen < NOW() - INTERVAL '10 seconds';`
		_, err := db.Exec(delete)
		if err != nil {
			log.Fatal(err)
		}
	}
}
