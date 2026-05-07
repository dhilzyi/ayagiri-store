## Overview
**POS RESTAURANT SYSTEM**

- Client
    - Every table has tablet that customer can order f&b
    - It will send the order through websocket to server backend

- Kitchen Client
    - Chef will look at screen order to see the orderes
    - It will being send in real time. Valueing time
    - 

- Server
    - Integration betweeen Client and Kitchen Client
    - Handling and store state

- Schemas Data
    - Orders
        - List<ProductID>
        - Timestamp

    - Table
        - TableID
        - Orders

    - Product
        - ProductID
        - ProductName
        - Price

## Stacks
- POSTGRES

- Golang Module
    - github.com/joho/godotenv
    - goose
    - sqlc

- Javascript and HTML vanilla

### Commands
psql "postgres://postgres:postgres@localhost:5432/ayagiri?sslmode=disable"

goose -dir ./sql/schema postgres "postgres://postgres:postgres@localhost:5432/ayagiri" up

