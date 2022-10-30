#Build stage
FROM golang:1.19-alpine3.16 AS buildstage
WORKDIR /app
COPY backend ./backend
COPY common ./common
COPY frontend ./frontend
RUN cd backend && go mod tidy
RUN cd common && go mod tidy
RUN cd frontend && go mod tidy
RUN cd frontend && go install gioui.org/cmd/gogio@latest
RUN cd frontend && gogio -target js -o dist cmd/main.go
RUN cd backend && go build cmd/server/main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=buildstage /app/backend/main .
COPY --from=buildstage /app/frontend/dist ./dist
COPY --from=buildstage /app/migrate ./migrate
COPY common/.env .
COPY start.sh .
COPY wait-for.sh .
COPY common/db/migration ./migration

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]