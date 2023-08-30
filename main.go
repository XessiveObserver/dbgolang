// file: main.go
package main

import (
	"database/sql"
	"fmt"
	"log"

	// we have to import the driver, but don't use it in our code
	// so we use the `_` symbol
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Bird struct {
	species, description string
}

func main() {
	// The `sql.Open` function opens a new `*sql.DB` instance. We specify the driver name
	// and the URI for our database. Here, we're using a Postgres URI
	conn := "postgres://postgres:postgres@localhost:5432/birding"
	db, err := sql.Open("pgx", conn)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	// To verify the connection to our database instance, we can call the `Ping`
	// method. If no error is returned, we can assume a successful connection
	if err := db.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}
	fmt.Println("database is reachable")

	// Returns multiple rows from the database
	rows, err := db.Query("SELECT bird, description FROM birds LIMIT 10")
	if err != nil {
		log.Fatalf("Couldn't execute query: %v", err)
	}
	// Instance of birds
	birds := []Bird{}

	// Iterate over returned rows
	for rows.Next() {
		bird := Bird{} // create instance of bird
		if err := rows.Scan(&bird.species, &bird.description); err != nil {
			log.Fatalf("Couldn't scan row: %v", err)
		}

		// append current bird instance to birds slice
		birds = append(birds, bird)
	}

	// print length of all birds plus displaying available birds
	fmt.Printf("found %d birds: %+v\n", len(birds), birds)

	// Inserting sample data
	newBird := Bird{
		species:     "rooster",
		description: "wakes you up in the morning",
	}

	// Exec command returns result instead of rows
	result, err := db.Exec("INSERT into birds (bird, description) VALUES ($1, $2)", newBird.species, newBird.description)
	if err != nil {
		log.Fatalf("could not insert row: %v", err)
	}

	// RowsAffected returns total number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("could not get affected rows: %v", err)
	}

	// log how many rows are affected
	fmt.Printf("inserted %v rows", rowsAffected)
}
