## Table of Contents

- [Prerequisites](#Prerequisites)
- [Installation and Run Application](#installation)

## Prerequisites
- Download and install [Go 1.21 or higher](https://go.dev/doc/install) (if you want to running it without docker)
- Install [Golang Migrate](https://github.com/golang-migrate/migrate) (if you want to running it without docker)
- Install Postgresql [Postgres](https://www.postgresql.org/download/)  (if you want to running it without docker)
- Download and Install [Docker](https://www.docker.com/products/docker-desktop/)
- Download and Install [Nodejs](https://nodejs.org/en/download/current) (if you want to running it without docker)

## Installation

### Running without docker
1. Clone the repository
   ```bash
   git clone https://github.com/fitraaditama7/order.git

2. Copy and fill value on .env
    ```bash
   cp .env.example .env
   
3. Running Application
    ```bash
   docker compose up -d --build
   
3. Open browser and go to `localhost:8000`

   
