# financial-chat

A simple chat application built with Go, Gorilla Websockets, and a simple html/js.

# Features
- Real-time chat messaging
- Multiple chat rooms
- Server-side stock quotes
    - Use the command `/stock=stock_code` (eg: `/stock=aapl.us`)

# Prerequisites
To run the application, you need the following:
- Go 1.20 or higher
- Docker

## Usage
1. Start docker containers:
    - `docker-compose up -d`
2. Run
- With make:
    - `make start`
- Whitout make: 
    - `go run .`

3. Open a web browser and go to `http://localhost:8081`
4. Enter a username and a room name, and start chatting!

- How to stop:
    - `make stop`
    - `docker-compose down`
- How to run tests:
    - `make test`

