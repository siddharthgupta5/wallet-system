# Wallet Management REST API

A simple RESTful backend service in Go to manage user wallets and basic transactions between them.

## Features

- Create users
- Create wallets for users
- Check wallet balance
- Transfer funds between wallets
- List all transactions for a wallet

## Technologies

- Go 1.21
- Gin web framework
- PostgreSQL database
- GORM ORM
- Docker

## Getting Started

### Prerequisites

- Docker and Docker Compose installed
- Go 1.21+ (for local development)

### Running with Docker

1. Clone the repository:
   ```bash
   git clone https://github.com/siddharthgupta5/wallet-api.git
   cd wallet-api
   ```

2. Start the services:
   ```bash
   docker-compose up --build
   ```
   The API will be available at http://localhost:8080

### Running Locally

1. Install Go 1.21+
2. Set up PostgreSQL database and update the `.env` file with your credentials
3. Run the application:
   ```bash
   go run cmd/main.go
   ```

## API Endpoints

### Users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/:id` - Get user details

### Wallets
- `POST /api/v1/wallets` - Create a new wallet
- `GET /api/v1/wallets/:id` - Get wallet details
- `GET /api/v1/wallets/:id/balance` - Get wallet balance
- `GET /api/v1/wallets/user/:id` - Get all wallets for a user

### Transactions
- `POST /api/v1/transactions` - Create a new transaction (credit/debit)
- `POST /api/v1/transactions/transfer` - Transfer funds between wallets
- `GET /api/v1/transactions/wallet/:id` - Get all transactions for a wallet

## Testing

To run tests:
```bash
go test -v ./...
```

## Project Setup

1. Initialize Git repository:
   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   ```

2. Create a GitHub repository and push:
   ```bash
   git remote add origin https://github.com/siddharthgupta5/wallet-api.git
   git push -u origin main
   ```

## API Testing Examples

You can use tools like Postman or cURL to test the endpoints. Here are some example requests:

### Create a user
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

### Create a wallet for the user
```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"currency":"USD"}'
```

### Check wallet balance
```bash
curl http://localhost:8080/api/v1/wallets/1/balance
```

### Transfer funds
```bash
curl -X POST http://localhost:8080/api/v1/transactions/transfer \
  -H "Content-Type: application/json" \
  -d '{"source_wallet_id":1,"destination_wallet_id":2,"amount":50,"description":"Test transfer"}'
```

### List transactions
```bash
curl http://localhost:8080/api/v1/transactions/wallet/1
```

## Implementation Details

This implementation includes:
- RESTful API design
- Proper database modeling with relationships
- Input validation
- Error handling
- Unit tests
- Dockerization
- Environment variables for configuration
- Clean project structure
- Comprehensive documentation

## License

This project is licensed under the MIT License.
