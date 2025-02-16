FROM golang:1.22-alpine3.18 AS build
RUN apk --no-cache add ca-certificates
WORKDIR /build
COPY . .
RUN go mod tidy
RUN GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o app .

FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /build
COPY --from=build /build/app .
EXPOSE 3000
CMD [ "./app" ]
