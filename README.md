# library-app

## About
handles e-library operations

## Third party libraries

1) [config ](github.com/kelseyhightower/envconfig) to load the env variables and binds to the struct
2) [log](go.uber.org/zap) for logging
3) [requests](https://github.com/gin-gonic/gin) for handling http requests
4) [testing](github.com/stretchr/testify/assert) using assert package for testing
5) [postgres](https://github.com/jackc/pgx) to store data in to postgres DB
6) [swag](https://github.com/swaggo/swag) for swagger documentation

## Config

`StoreType` - Defines type of store going to use to run the app supported values: `local` (default) and `postgres`.

## Test and Run

`make run`: to up and run the application in local system

`make test`: runs all test cases and show the result in html

`make dockerdeploy`: up and run as docker container

`make helminstall`: deploys the application in kubernetes environment

# Swagger
[swagger](http://localhost:3000/swagger/index.html) renders swagger doc.

## Requests

### GetAllBooks

#### Request

```
curl -X 'GET' \
  'http://localhost:3000/api/v1/book' \
  -H 'accept: application/json'
```

### GetBook

#### Request

```
curl --location --request GET 'localhost:3000/api/v1/book/book_10'
```

### BorrowBook

#### Request

```
curl --location 'localhost:3000/borrow' \
--header 'Content-Type: application/json' \
--data '{
    "title": "book_1",
    "name_of_borrower": "sandeep"
}'
```

### ExtendLoan

#### Request

```
curl --location --request POST 'localhost:3000/api/v1/extend/1'
```

### ReturnBook

#### Request

```
curl --location --request POST 'localhost:3000/api/v1/return/1'
```