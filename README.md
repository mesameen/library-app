# library-app

## About
handles library operations

## Third party libraries

1) [config ](github.com/kelseyhightower/envconfig) to load the env variables and binds to the struct
2) [log](go.uber.org/zap) for logging
3) [requests](https://github.com/gin-gonic/gin) for handling http requests
4) [testing](github.com/stretchr/testify/assert) using assert package for testing

## Config

`StoreType` - Defines type of store going to use in this app values: local, postgres etc. Based on this value Store can be created during the bootup.
