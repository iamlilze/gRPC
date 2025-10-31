# gRPC Authorization Service

A simple authentication and authorization service built with gRPC in Go.

## Features

- User registration with password hashing (bcrypt)
- User login with token generation
- Token validation
- User logout with token revocation
- In-memory storage for users and tokens

## Project Structure

```
.
├── api/
│   └── auth/           # Generated protobuf code
├── cmd/
│   └── server/         # Server entry point
├── internal/
│   ├── service/        # Business logic
│   └── storage/        # Data storage layer
└── proto/              # Protocol buffer definitions
```

## Prerequisites

- Go 1.24 or later
- Protocol Buffers compiler (protoc)
- gRPC Go plugins

## Installation

1. Clone the repository:
```bash
git clone https://github.com/iamlilze/gRPC.git
cd gRPC
```

2. Install dependencies:
```bash
go mod download
```

## Building

Build the server:
```bash
go build -o bin/server ./cmd/server/main.go
```

## Running

Start the gRPC server:
```bash
./bin/server
```

The server will start on port 50051 by default. You can change this by setting the `PORT` environment variable:
```bash
PORT=8080 ./bin/server
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

## API

The service provides the following gRPC methods:

### Register
Creates a new user account.

**Request:**
```protobuf
message RegisterRequest {
  string username = 1;
  string email = 2;
  string password = 3;
}
```

**Response:**
```protobuf
message RegisterResponse {
  string user_id = 1;
  string message = 2;
  bool success = 3;
}
```

### Login
Authenticates a user and returns a token.

**Request:**
```protobuf
message LoginRequest {
  string username = 1;
  string password = 2;
}
```

**Response:**
```protobuf
message LoginResponse {
  string token = 1;
  string user_id = 2;
  bool success = 3;
  string message = 4;
}
```

### ValidateToken
Checks if a token is valid.

**Request:**
```protobuf
message ValidateTokenRequest {
  string token = 1;
}
```

**Response:**
```protobuf
message ValidateTokenResponse {
  bool valid = 1;
  string user_id = 2;
  string username = 3;
}
```

### Logout
Invalidates a user's token.

**Request:**
```protobuf
message LogoutRequest {
  string token = 1;
}
```

**Response:**
```protobuf
message LogoutResponse {
  bool success = 1;
  string message = 2;
}
```

## Using grpcurl

You can test the service using [grpcurl](https://github.com/fullstorydev/grpcurl):

### Register a user:
```bash
grpcurl -plaintext -d '{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}' localhost:50051 auth.AuthService/Register
```

### Login:
```bash
grpcurl -plaintext -d '{
  "username": "testuser",
  "password": "password123"
}' localhost:50051 auth.AuthService/Login
```

### Validate token:
```bash
grpcurl -plaintext -d '{
  "token": "YOUR_TOKEN_HERE"
}' localhost:50051 auth.AuthService/ValidateToken
```

### Logout:
```bash
grpcurl -plaintext -d '{
  "token": "YOUR_TOKEN_HERE"
}' localhost:50051 auth.AuthService/Logout
```

## Development

### Regenerate protobuf code:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/auth.proto
mv proto/auth.pb.go proto/auth_grpc.pb.go api/auth/
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.