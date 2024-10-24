# Packs Calculator Service

This application calculates the optimal pack sizes for a given order quantity.
It also allows adding and deleting pack sizes.

## How to Run

To run the application, use Docker Compose. This command builds the API, sets up a Postgres database, runs migrations, and executes tests in separate containers:

```bash
docker-compose up --build
```

Once the application is running, you can access the static HTML at:

**http://localhost:8080/static/**

This page allows you to interact with the API.

## Available Endpoints

Here are the available API endpoints and their respective parameters:

1. **POST `/api/v1/packs`**
    - Adds a new pack size.
    - **Body**: `{ "size": <pack_size> }`
    - Example:
      ```json
      {
        "size": 1000
      }
      ```

2. **DELETE `/api/v1/packs`**
    - Deletes an existing pack size.
    - **Body**: `{ "size": <pack_size> }`
    - Example:
      ```json
      {
        "size": 1000
      }
      ```

3. **GET `/api/v1/calculate?quantity=<order_quantity>`**
    - Calculates the minimum number of packs needed for the given order quantity.
    - **Query parameter**: `quantity` (order quantity)
    - Example:
      ```
      /api/v1/calculate?quantity=5000
      ```