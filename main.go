// Butter tracker package enables the tracking of nodes deployed on a Butter network. It is
// mostly used for metrics and testing.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
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

// createSchema for postgres database of Butter nodes
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

// This information could be used to generate a live graph of the network and its PCG overlay by
// integrating the database/TRacker server with websockets

func main() {
	db := ConnectToDB()
	defer db.Close()
	createSchema(db)

	// Nodes periodically ping the Tracker telling the server  they are still alive and hence to
	// maintain and update their entry in the databse
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
