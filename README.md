# Bhagad-House-Booking
[Bhagad House Booking](https://bhagadhouse.com) is a private app created in golang with [gioui](https://gioui.org/) and other libraries. 
It's a private app but it's code are open sourced and MIT-Licensed. You are free to use it.

## Prerequisites
1. [golang](https://go.dev/)
2. [gioui](https://gioui.org/)
3. [Docker or Docker Desktop](https://www.docker.com/)
4. [golang-migrate](https://pkg.go.dev/github.com/golang-migrate/migrate/v4@v4.15.2)
5. [Postgresql 14 or later](https://www.postgresql.org/download/)

## For Local Development

### Database
1. Make sure your postgresql database is running at 5432.
2. Create a database named ```bhagad_house_booking``` locally (DB_URL=postgresql://root:secret@localhost:5432/bhagad_house_booking?sslmode=disable)
and necessary tables.
```#!console
postgres15 createdb --username=root --owner=root bhagad_house_booking
cd common && make migrateup
```
3. Create a user (Optional, only if admin role is desired).
```#!console
cd common && go run cmd/main.go
```
The above command will create a user named ```Owner Admin``` and email ```admin@bhagadhouse.com``` and password ```12345678```.<br>
You will need to add a role named ```Admin``` in the users table of bhagad_house_booking

### Backend
1. Run the backend app locally on port 8001
```#!console
export PORT=8001 && cd backend && go run cmd/server/main.go
```

### Frontend
1. Build the static files for web
```#!console
cd frontend && gogio -target js -o output/wasm ./cmd/main.go
```
2. Copy generated wasm file to static directory
```#!console
cd frontend && cp output/wasm/main.wasm cmd/static/dist/main.wasm
```
3. Copy environment variables for frontend(client) app
```#!console
echo "window.API_URL = 'http://localhost:8001';" >> frontend/cmd/static/dist/wasm.js
```
3. Serve the static files
```#!console
export window.STATIC_FOLDER = 'dist'
export window.INNER_PORT = '8002'
cd frontend/cmd/static && go run main.go
```

## Development using docker-compose
1. Build and run docker images for backend,frontend and postgresql.
```
docker compose up
```
2. Refer to [Local Development](#for-local-development)

## Compile for windows on Arch Linux
```#!console
 CGO_ENABLED=1 GOOS=windows CC=x86_64-w64-mingw32-gcc go build -ldflags -H=windowsgui -o output/windows/booking.exe cmd/main.go
```
Note: [-ldflags -H=windowsgui](https://stackoverflow.com/questions/23250505/how-do-i-create-an-executable-from-golang-that-doesnt-open-a-console-window-whe)


## Development 
### Dependencies ###
* The frontend needs to be served as static app
* The frontend also needs [env vars](frontend/config.go)
* The backend in turn needs postgres database to be running.

### Steps 

### Production ###
flyctl deploy -a bhagad-house-booking-frontend --build-arg API_URL=https://bhagad-house-booking-backend.fly.dev -e GIN_MODE=release -c frontend.fly.toml

flyctl deploy -a bhagad-house-booking-backend -e GIN_MODE=release -c backend.fly.toml
