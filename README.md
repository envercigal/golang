# Golang example

## 🧾 Overview

This is a Go-based application that can be easily run using Docker Compose. It includes a complete test suite written in Go.

#### Not: The driver ID was added to make it easier to identify the driver during testing.
#### Not: “During Circuit Breaker testing, in the half-open state, a ‘too many requests’ error was returned, and in the open state, a ‘system error’ was returned as the status code.”
## 🛠 Prerequisites

Before getting started, make sure you have the following installed:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## 🚀 Running the Application

To start the application using Docker Compose, run the following command in the root directory of the project:

```bash
   docker compose up
```


To run all tests in the project, use:
```bash
  go test ./...
```

[Download Postman Collection](./docs/driver_location.postman_collection.json)