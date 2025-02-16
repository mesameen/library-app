# library-app

## About
handles library operations

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

`make run`: to up and run the application

`make test`: runs all test cases and show the result in html


# Swagger
[swagger](http://localhost:3000/swagger/index.html) doc is running here
