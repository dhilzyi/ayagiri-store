## Description
This is a mock of POS Store, customers will order using tablet at their table. There's also shop side interface for daily orders and admin interface for managing basic CRUD to the database.

## Tech Stacks
### Front End
I'm aiming simple interface throughout the project. No framework and fancy css. Just some basic vanilla JS
- HTML5
- CSS
- Javascript Vanilla
### Back End
- Postgresql for database
- Golang 
    - sqlc :generating sql query automatically to golang program
    - goose :making schema version is easier
    - goplayground/validator :for validating

## Quick Start
I plan to create docker in the future. for now you can start the program with bash script start.sh

- Rename .env.example to .env
- Set up necessary variable to .env file provided
```bash
$ chmod +x ./start.sh
$ ./start.sh
```

## Structure Project
- Server
```
    ├── cmd
    │   └── server
    │       └── main.go - entry file. adding the handler to the endpoints
    ├── go.mod
    ├── go.sum
    ├── internal
    │   ├── api
    │   │   ├── handler_categories.go
    │   │   ├── handler.go
    │   │   ├── handler_orders.go
    │   │   ├── handler_products.go
    │   │   ├── helpers.go
    │   │   ├── helper_test.go
    │   │   └── type.go
    │   ├── database
    │   │   ├── categories.sql.go
    │   │   ├── copyfrom.go
    │   │   ├── db.go
    │   │   ├── models.go
    │   │   ├── orders.sql.go
    │   │   └── products.sql.go
    │   ├── domain
    │   │   └── types.go
    │   └── orders
    │       ├── orders.go
    │       └── types.go
    ├── README.md
    ├── sql
    │   ├── queries
    │   │   ├── categories.sql
    │   │   ├── orders.sql
    │   │   └── products.sql
    │   └── schema
    │       ├── 001_create_categories.sql
    │       ├── 002_create_products.sql
    │       └── 003_create_order.sql
    ├── sqlc.yaml
    ├── start.sh
```
