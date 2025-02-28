# Ranktify Backend

## Installation

1. **Install Go**: Download and install Go by following the instructions at [https://go.dev/doc/install](https://go.dev/doc/install).

2. **Install Dependencies**:
```bash
go mod tidy
```

## Running the Application

1.  **Run the Application**:
     Before running the application, make sure to create .env file with the folliwing format:
```python
DB_HOST=localhost
DB_PORT=9090
DB_USER=ranktifyUser
DB_PASSWORD=concalma
DB_NAME=ranktify
DB_SSLMODE=disable # For local development: disabled, production: required
```
And finall run the app:

```bash
#running the application
go run cmd/main.go
```

## Docker Setup

1. Enter the 'local' directory to find the docker-compose.yaml and the init.sql files

2. Install docker, follow the instructions below

     - Linux (docker engine): https://docs.docker.com/engine/install/
     - Windows (docker desktop): https://docs.docker.com/desktop/setup/install/windows-install/


### Linux setup

Make sure to have docker compose. In the local directory run:
```bash
docker compose up -d
```

Open datagrip and enter the following credentials in the data source properties: 

     Host: localhost
     Port: 9090
     User: ranktifyUser
     Password: concalma
     Database: ranktify

### Windows setup

**->When someone is able to do this, please update this<-**