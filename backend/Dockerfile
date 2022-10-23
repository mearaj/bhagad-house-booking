#Build stage
FROM golang:1.19-alpine3.16 AS buildstage
WORKDIR /app
COPY . .
RUN go build cmd/server/main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=buildstage /app/main .

EXPOSE 8080
CMD ["/app/main"]