# financial-chat

A simple chat application built with Go with Gorilla Websockets and HTML/JS.

# Features
- Real-time chat messaging
- Multiple chat rooms (create a new window and connect using a different room name)
- Store messages in cache and load old messages when a new user enter the room
- Bot commands
    - Type `/help` in chat
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

- Change environment variables in `.env` file
- How to stop:
    - `make stop`
    - `docker-compose down`

## Testing
Tests are using mocks to inject service dependencies
Mocks were created with Mockery

- How to run tests:
    - `make test`
- Creating mocks
    - mockery --dir=./infra/cache --all --output=./infra/cache/mocks --outpkg=mocks
    - mockery --dir=./infra/amqp --all --output=./infra/amqp/mocks --outpkg=mocks
    - mockery --dir=./app/websocket --all --output=./app/websocket/mocks --outpkg=mocks
    - mockery --dir=./app/consumer --all --output=./app/consumer/mocks --outpkg=mocks

## Project structure
- "app": config file and handlers/services divided by domain 
    - websocket (chat)
    - consumer (queue handler)
- "infra": external interfaces with implementations (redis and rabbitmq)
- "public": frontend files