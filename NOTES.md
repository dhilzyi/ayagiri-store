## Overview
**POS RESTAURANT SYSTEM**

HIGH DIAGRAM SYSTEM FLOW

Kitchen will open SSE to backend inital and listening to the new orders.

Client fetch products data from backend -> Make an order POST to backend + Open an Stream for signalling to that orderID -> Backend receives the order send the data through the all kitchen channel
-> Kitchen will send a signal when the order is complete -> Backend receive a signal for that orderID -> Backend will send a stream to that client who hold the orderID.

- Client
    - Every table has tablet that customer can order f&b
    - It will send the order to server backend

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

- Golang Package
    - github.com/joho/godotenv
    - goose
    - sqlc


### Commands
psql "postgres://postgres:postgres@localhost:5432/ayagiri?sslmode=disable"

goose -dir ./sql/schema postgres "postgres://postgres:postgres@localhost:5432/ayagiri" up

pg_dump "postgres://postgres:postgres@localhost:5432/ayagiri?sslmode=disable" -d ayagiri --data-only > seed.sql

- Rebuild after code changed
docker compose build

- Run all compose up
docker compose up
- Reset compose. deleting containers and network. Use `-v ` for deleting with its volume
docker compose down

- Run compose detach only db
docker compose up -d db
- Run migration compose
docker compose run --rm migrate
- Inspect yaml
docker compose config
- Stop compose docker
docker compose stop
- Get inside shell of container
docker exec -it <container-name> sh

- dump database
sudo -u postgres pg_dump ayagiri > backup.sql
- dump data only database
sudo -u postgres pg_dump -a <temp_restore> > <data_only.sql>
- restore database
sudo -u postgres psql -d ayagiri < backup.sql

- Connect to psql
sudo -u postgres psql
