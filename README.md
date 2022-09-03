# bhagad-house-booking

## Prerequisites

1. [golang](https://go.dev/)
2. [gioui](https://gioui.org/)

## Development
To run the app locally
```#!console
go run cmd/main.go
```
To build the app locally
```#!console
go build cmd/main.go
```

## Compile for windows on Arch Linux
```#!console
GOOS=windows go build -o output/windows/booking.exe cmd/main.go
```