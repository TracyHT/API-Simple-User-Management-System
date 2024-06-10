
# User Management API

This is a simple User Management API built with Go and the Gin web framework, using SQLite as the database. The API allows you to create, read, update, and delete user records, and also list users with optional filtering and sorting.

## Features

- Create a new user
- Get user details by ID
- Update user details by ID
- Delete user by ID
- List users with optional filtering and sorting

## Endpoints

### Create a User

```http
POST /users
```

#### Request Body

```json
{
  "username": "john_doe",
  "firstname": "John",
  "lastname": "Doe",
  "email": "john@example.com",
  "avatar": "http://example.com/avatar.jpg",
  "phone": "+1234567890",
  "date_of_birth": "1990-01-01",
  "address_country": "USA",
  "address_city": "New York",
  "address_street_name": "Broadway",
  "address_street_address": "123"
}
```

### Get a User by ID

```http
GET /users/:id
```

### Update a User by ID

```http
PUT /users/:id
```

#### Request Body

Provide only the fields you want to update. Example:

```json
{
  "firstname": "John Updated",
  "email": "john_updated@example.com"
}
```

### Delete a User by ID

```http
DELETE /users/:id
```

### List Users

```http
GET /users
```

#### Query Parameters

- `username`: Filter by username
- `email`: Filter by email
- `firstname`: Filter by first name
- `lastname`: Filter by last name
- `address_country`: Filter by country
- `sort_by`: Field to sort by (default: `id`)
- `sort_order`: Sorting order (`asc` or `desc`, default: `asc`)

Example:

```http
GET /users?address_country=Vietnam&sort_by=username&sort_order=asc
```

## Running the Application

1. Clone the repository and open in IDE:

```sh
git clone
```

2. Install dependencies:

```sh
go mod tidy
```

3. Run the application:

```sh
go run main.go
```

The application will start on `http://localhost:8080`.

## Dependencies

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [SQLite Driver](https://github.com/mattn/go-sqlite3)

## Using Random Data API

If any user information field is empty, the backend will generate random user information by calling the [random-data-api.com](https://random-data-api.com/) services.

## Using Hoppscotch to Test

Hoppscotch is a web-based API testing tool. Follow the steps below to test the API endpoints:

1. Go to [Hoppscotch](https://hoppscotch.io/).
2. Set the request method and URL for the desired endpoint.
3. Add any necessary headers (e.g., `Content-Type: application/json`).
4. For POST and PUT requests, provide the request body in JSON format.
5. Send the request and view the response.
