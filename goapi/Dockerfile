##
## Build
##
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY cmd /app/cmd

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /go-api ./cmd


##
## Deploy
##
FROM alpine:latest
# RUN apk --no-cache add ca-certificates
COPY --from=build /go-api ./
RUN chmod +x ./go-api
ENTRYPOINT ["./go-api"]
EXPOSE 1323