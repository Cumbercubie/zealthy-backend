# Zealthy Ticket Support Back-end System

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)

This repository is back-end web API for Zealthy's ticket support system.

## Clone the project

```
$ git clone https://github.com/Cumbercubie/zealthy-backend.git
$ cd zealthy-backend
```

## Export the .env variable

```
$ export DATABASE_URL=<database_url>
```

## Start the server
```
$ cd zealthy-backend
$ go run main.go
```

## Run tests
```
$ go test ./internal/test
```
Available APIs:

* `GET/v1/api/tickets`: List all tickets with optional filter by `status` params
* `GET/v1/api/ticket/:ticket_id`: Get ticket with specified `ticket_id`
* `POST/v1/api/ticket`: Create ticket
* `PUT/v1/api/ticket/:id`: Update a ticket with specified `id`

## `GET/v1/api/tickets`

params:

* status: NEW | IN_PROGRESS | RESOLVED

## `GET/v1/api/ticket/:ticket_id`

params:

* ticket_id: ticket_id field of desired ticket

## `POST/v1/api/ticket`

payload:

* name:
* description:
* issuer_email:

### NOTE:

* __Status will be new upon creation__
* __Response will be empty upon creation__

## `PUT/v1/api/ticket/:id`

### NOTE:

* Will not update ticket if status is RESOLVED