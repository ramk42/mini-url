# Mini URL

Mini URL is a URL shortening service written in Go. It provides a simple API to shorten URLs and redirect users to the original URLs.

## Features

- Shorten long URLs
- Redirect to original URLs
- Set expiration dates for shortened URLs
- Track click counts
- Health check endpoint

## Requirements

- Go 1.23 or later
- Docker
- Docker Compose

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/ramk42/mini-url.git
    cd mini-url
    ```

2. Build and run the application using Docker Compose:
    ```sh
    make up
    ```

## Configuration
Only for development environment, you can use the following command to run the application :

```sh
- `DB_HOST`: Database host (default: `localhost`)
- `DB_PORT`: Database port (default: `5432`)
- `DB_USER`: Database user (default: `postgres`)
- `DB_PASSWORD`: Database password (default: `postgres`)
- `DB_NAME`: Database name (default: `mini_url`)
- `DB_SSLMODE`: Database SSL mode (default: `disable`) 
- `BASE_URL`: Base URL for the shortened URLs (default: `http://localhost:8080`)
```

## API Endpoints

### Shorten URL

- **URL**: `/shorten`
- **Method**: `POST`
- **Request Body**:
    ```json
    {
        "long_url": "https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html",
        "expiration_days": 30
    }
    ```
- **Response**:
    ```json
    {
        "short_url": "http://localhost:8080/abc123"
    }
    ```

### Resolve URL

- **URL**: `/[slug]`
- **Method**: `GET`
- **Response**: Redirects to the original URL

### Health Check

- **URL**: `/health`
- **Method**: `GET`
- **Response**: `OK`

## Development

### first time setup

To setup the project for the first time, use the following command:

You have to create a `.env` file in the root directory of the project with the following content:

```sh
echo .env.example > .env
```

Run the postgres database using the following command:

```sh
make up_database
```

execute the following command to create the tables in the database:

```sql
CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL UNIQUE,
    slug VARCHAR(10) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    clicks INT DEFAULT 0 CHECK(clicks >= 0)
);

CREATE INDEX IF NOT EXISTS idx_slug ON urls(slug);
```
Note : The database will be created automatically when the application is run for the first time.

### Running Tests

To run unit tests, use the following command:
```sh
make go-test
```
#### test scripts
You can test the two api endpoints using the following scripts:

```sh
godotenv -f .env make go-run
./test-scripts/shorten.sh
./test-scripts/resolve.sh
./test-scripts/vegeta-test.sh # for performance testing
```
### running the application without Docker
install the `godotenv` package using the following command:

```sh
go get github.com/joho/godotenv
```

To run the application without Docker with hot reloading, use the following command:

```sh
godotenv -f .env make go-run
```

### Linting

To run the linter, use the following command:
```sh
make code-quality
```

## Acknowledgements

- [Go](https://golang.org/)
- [Chi](https://github.com/go-chi/chi)
- [Zerolog](https://github.com/x/zerolog)
- [Purell](https://github.com/PuerkitoBio/purell)
- [Docker](https://www.docker.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [Pgconn](github.com/jackc/pgx/v5/pgconn)


Happy reading! 😇
