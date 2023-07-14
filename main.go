package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "tenant"
)

func main() {
	// Establish database connection
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create master schema if it doesn't exist
	if _, err := db.Exec("CREATE SCHEMA IF NOT EXISTS master;"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS master.tenants(
		id INT NOT NULL,
		name VARCHAR(200)
	)`); err != nil {
		log.Fatal(err)
	}

	// Register a new tenant and create their schema
	tenantID := 1
	tenantName := "example_tenant"

	if _, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", tenantName)); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(fmt.Sprintf("INSERT INTO master.tenants (id, name) VALUES (%d, '%s');", tenantID, tenantName)); err != nil {
		log.Fatal(err)
	}

	// create table in $tenantName schema
	if _, err := db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.example (id INT NOT NULL, name VARCHAR(200))", tenantName)); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(fmt.Sprintf("INSERT INTO %s.example (id, name) VALUES (%d, '%s');", tenantName, 1, "test insert")); err != nil {
		log.Fatal(err)
	}

	// Switch to the tenant's schema
	if _, err := db.Exec(fmt.Sprintf("SET search_path TO %s;", tenantName)); err != nil {
		log.Fatal(err)
	}

	// Example query
	rows, err := db.Query("SELECT * FROM example;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		log.Printf("id: %d, name: %s", id, name)
	}

	// Process query results
	// ...

	// Switch back to the master schema if needed
	if _, err := db.Exec("SET search_path TO master;"); err != nil {
		log.Fatal(err)
	}
}
