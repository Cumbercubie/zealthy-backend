package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	repositories "github.com/zealthy/helpdesk-backend/internal/adapters"
	"github.com/zealthy/helpdesk-backend/internal/adapters/handlers"
	"github.com/zealthy/helpdesk-backend/internal/core/services"
	sql "github.com/zealthy/helpdesk-backend/internal/infra/database"
	"github.com/zealthy/helpdesk-backend/internal/infra/server"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		// Handle error loading .env file
		log.Fatalf("Error loading .env file: %v\n", err)
		os.Exit(1)
	}
	dbpool, derr := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	if derr != nil {
		log.Fatalf("Unable to create connection pool: %v\n", derr)
		os.Exit(1)
	}

	defer dbpool.Close()

	dbHandler := sql.NewDBHandler(dbpool)
	err := sql.CreateAndSetSchema(dbHandler, "zealthy")
	err = sql.CreateTicketsTable(dbHandler)

	if err != nil {
		log.Println("Error creating tickets table", err.Error())
	}
	ticketRepo := repositories.NewTicketRepository(dbHandler)
	ticketService := services.NewTicketService(ticketRepo)
	ticketHandler := handlers.NewTicketHandler(ticketService)

	r := server.NewServer()

	ticketHandler.RegisterTicketRoute(r)

	interruptChannel := make(chan os.Signal, 1)

	go func() {
		r.Run(":3000")
	}()

	<-interruptChannel

	os.Exit(0)
}
