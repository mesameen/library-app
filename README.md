# library-app

## About
handles library operations

## Third party libraries

1) [config ](github.com/kelseyhightower/envconfig) to load the env variables and binds to the struct
2) [log](go.uber.org/zap) for logging
3) [requests](https://github.com/gin-gonic/gin) for handling http requests
4) [testing](github.com/stretchr/testify/assert) using assert package for testing

## Test and Run

`make run`: to up and run the application

`make test`: runs all test cases and show the result in html

## Config

`StoreType` - Defines type of store going to use in this app values: local, postgres etc. Based on this value Store can be created during the bootup.

## Requests

### GetBook

#### Request

```
curl --location --request GET 'localhost:3000/book/book_10'
```

#### Responses

##### Success
```json
{
    "title": "Book_1",
    "available_copies": 4
}
```
##### Failure
`404 - Not Found`
```json
{
    "error": "book with title 'book_10' isn't presents. not found"
}
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

#### Responses

##### Success
```json
{
    "id": 4,
    "name_of_borrower": "sandeep",
    "title": "book_1",
    "loan_date": 1739612322,
    "return_date": 1740217122
}
```
##### Failure
`404 - Not Found`
```json
{
    "error": "book with title 'book_10' isn't presents. not found"
}
```
`400 - Bad Request`
```json
{
    "error": "NameOfBorrower or Title are missed in the request"
}
```

### ExtendLoan

#### Request

```
curl --location --request POST 'localhost:3000/extend/1'
```

#### Responses

##### Success
```json
{
    "message": "loan got extended to 3 weeks",
    "return_date": 1742031488
}
```
##### Failure
`404 - Not Found`
```json
{
    "error": "loan 12 isn't presents not found"
}
```

### ReturnBook

#### Request

```
curl --location --request POST 'localhost:3000/return/1'
```

#### Responses

##### Success
```json
{
    "message": "book returned"
}
```

##### Failure
`404 - Not Found`
```json
{
    "error": "loan 12 isn't presents not found"
}
```
