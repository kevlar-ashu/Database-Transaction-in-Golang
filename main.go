package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user-name"
	password = "your-password"
	dbname   = "database-name"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Create a new connection to our database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// `tx` is an instance of `*sql.Tx` through which we can execute our queries

	// Here, the query is executed on the transaction instance, and not applied to the database yet
	_, err = tx.ExecContext(ctx, "INSERT INTO pets (name, species) VALUES ('Fido', 'dog'), ('Albert', 'cat')")
	if err != nil {
		// Incase we find any error in the query execution, rollback the transaction
		tx.Rollback()
		return
	}

	// The next query is handled similarly
	_, err = tx.ExecContext(ctx, "INSERT INTO food (name, quantity) VALUES ('Dog Biscuit', 3), ('Cat Food', 5)")
	if err != nil {
		tx.Rollback()
		return
	}

	// Run a query to get a count of all cats
	row := tx.QueryRow("SELECT count(*) FROM pets WHERE species='cat'")
	var catCount int
	// Store the count in the `catCount` variable
	err = row.Scan(&catCount)
	if err != nil {
		tx.Rollback()
		return
	}

	// Now update the food table, increasing the quantity of cat food by 10x the number of cats
	_, err = tx.ExecContext(ctx, "UPDATE food SET quantity=quantity+$1 WHERE name='Cat Food'", 10*catCount)
	if err != nil {
		tx.Rollback()
		return
	}

	// Finally, if no errors are recieved from the queries, commit the transaction
	// this applies the above changes to our database
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}
