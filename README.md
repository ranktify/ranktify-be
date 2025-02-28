# Ranktify Backend

## Installation

1. **Install Go**: Download and install Go by following the instructions at [https://go.dev/doc/install](https://go.dev/doc/install).

2. **Install Dependencies**:
```bash
go mod tidy
```

## Running the Application

1.  **Run the Application**:
     Before running the application, we need to set if we are in local development or in production, to do that run one of the following `export` :

```bash
# local development
export APP_ENV=local
# production mode
export APP_ENV=prod
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