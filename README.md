# IP2Country Service

IP2Country Service is a Go-based REST API that takes an IP address as input and returns the corresponding country and city information. The service includes rate limiting and can be easily extended to support different IP-to-country databases.

## Features

- **REST API**: Exposes a simple HTTP GET endpoint to query country and city by IP address.
- **Rate Limiting**: Configurable rate limiting to control the number of requests per second.
- **Extendable**: Designed to support multiple IP-to-country databases.
- **Configuration**: Reads configuration from environment variables.
- **Tests**: Includes unit tests for key components.
- **Docker Support**: Provides a Dockerfile for containerized deployment.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.20 or later)
- [Docker](https://www.docker.com/) (optional, for containerized deployment)

### Installation

1. Clone the repository:

   ```
   git clone https://github.com/ran3396/ip2country.git
   cd ip2country
   ```

2. Build the project:

   ```
   go build -o ip2country
   ```

### Configuration

The service is configured using environment variables:

- `PORT`: The port on which the service will listen (default: `8080`).
- `RATE_LIMIT`: The maximum number of requests per second (default: `100`).
- `IP_DB_PATH`: The path to the IP-to-country CSV database file.

Example configuration:

```
export PORT=8080
export RATE_LIMIT=100
export IP_DB_PATH=./testdata/ipdb.csv
```

### Running the Service

1. Set the environment variables:

   ```
   export PORT=8080
   export RATE_LIMIT=100
   export IP_DB_PATH=./testdata/ipdb.csv
   ```

2. Run the service:

   ```
   ./ip2country
   ```

The service will be available at `http://localhost:8080`.

### API Usage

**Endpoint**: `/api/v1/find-country`

**Method**: `GET`

**Query Parameters**:
- `ip`: The IP address to query.

**Example Request**:

```
curl "http://localhost:8080/api/v1/find-country?ip=2.22.233.255"
```

**Example Response**:

```
{
    "country": "United Kingdom",
    "city": "London"
}
```

### Running Tests

To run the tests, use the following command:

```
go test ./...
```

### Docker Deployment

1. Build the Docker image:

   ```
   docker build -t ip2country .
   ```

2. Run the Docker container:

   ```
   docker run -p 8080:8080 -e PORT=8080 -e RATE_LIMIT=100 -e IP_DB_PATH=/path/to/ipdb.csv ip2country
   ```

### Project Structure
```
ip2country/
├── Dockerfile
├── README.md
├── go.mod
├── go.sum
├── main.go
├── middleware/
│ └── rate_limiter.go
├── config/
│ └── config.go
├── handlers/
│ ├── ip_handler.go
│ └── ip_handler_test.go
├── ipdb/
│ └── ip_database.go
├── utils/
│ └── response.go
└── testdata/
└── ipdb.csv
```
