package sql

import (
	"context"
	"fmt"
	"log"
)

const (
	CHECK_TABLE_EXIST_QUERY = "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1)"
)

func CreateAndSetSchema(db *DBHandler, schema string) error {
	dbConn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return err
	}

	defer dbConn.Release()
	// Create schema if it doesn't exist
	_, err = dbConn.Exec(context.Background(), fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, schema))
	if err != nil {
		return fmt.Errorf("error creating schema: %w", err)
	}
	// Set search path to use the new schema
	_, err = dbConn.Exec(context.Background(), fmt.Sprintf(`SET search_path TO %s`, schema))
	if err != nil {
		return fmt.Errorf("error setting search path: %w", err)
	}

	log.Printf("Schema '%s' is ready and set as default", schema)
	return nil
}

func CreateTicketsTable(db *DBHandler) error {
	// Check if the table exists
	dbConn, err := db.getDbConnection()

	defer dbConn.Release()
	var exists bool

	_, err = dbConn.Exec(context.Background(), CHECK_TABLE_EXIST_QUERY, "tickets")
	if err != nil {
		return fmt.Errorf("error checking if tickets table exists: %w", err)
	}

	if !exists {
		// Create the table
		_, err = dbConn.Exec(context.Background(), `
            CREATE TABLE tickets (
				id SERIAL PRIMARY KEY,
                ticket_id UUID NOT NULL,
                name VARCHAR(255) NOT NULL,
                issuer_email VARCHAR(255) NOT NULL,
                description VARCHAR(255) NOT NULL,
                status VARCHAR(50) NOT NULL,
                response VARCHAR(255),
				note VARCHAR(255),
				created_at TIMESTAMP NOT NULL
				updated_at TIMESTAMP NOT NULL,
            )
        `)
		if err != nil {
			return fmt.Errorf("error creating tickets table: %w", err)
		}
		fmt.Println("Tickets table created successfully")
	} else {
		fmt.Println("Tickets table already exists")
	}

	return nil
}
