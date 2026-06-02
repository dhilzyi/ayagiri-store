## Description
This is a mock of POS Store, customers will order using tablet at their table. There's also shop side interface for daily orders and admin interface for managing basic CRUD to the database.

## Tech Stacks
### Front End
I'm aiming simple interface throughout the project. No framework and fancy css.
- HTML5
- CSS
- Javascript Vanilla
### Back End
- Postgresql for database
    - It use pgx driver
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
    │   ├── api - handler function and initial
    │   ├── database - auto generated from .sql by sqlc program
    │   ├── domain
    │   └── orders - orders manager module 
    ├── sql
    │   ├── queries - raw sql queries for each table
    │   │   ├── categories.sql
    │   │   ├── orders.sql
    │   │   └── products.sql
    │   └── schema - schema migration version
    │       ├── 001_create_categories.sql
    │       ├── 002_create_products.sql
    │       └── 003_create_order.sql
    ├── start.sh - simple script to run server
```
- Web
```
└── web
    ├── customer - customer side clients
    │   ├── css
    │   ├── index.html
    │   └── js
    └── kitchen - kitchen side clients
        ├── css
        │   └── components
        ├── index.html
        ├── js
        └── styles.css
```
