# financial-chat

Simple browser-based chat application using Go

This application allows several users to talk in a chatroom and also to get stock quotes from an API using the specific command `/stock=stock_code`

## Commands
- How to run:
    - To run you must have docker installed, run `docker-compose up -d`
    - With make: `make start`
    - Whitout make: `docker-compose up -d` and `go run .`
    - Access `http://localhost:8081`
- How to stop:
    - `make stop`
- How to run tests:
    - `make test`

// TODO
- implement bot
