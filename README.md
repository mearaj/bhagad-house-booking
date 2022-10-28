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
 CGO_ENABLED=1 GOOS=windows CC=x86_64-w64-mingw32-gcc go build -ldflags -H=windowsgui -o output/windows/booking.exe cmd/main.go
```
Note: [-ldflags -H=windowsgui](https://stackoverflow.com/questions/23250505/how-do-i-create-an-executable-from-golang-that-doesnt-open-a-console-window-whe)