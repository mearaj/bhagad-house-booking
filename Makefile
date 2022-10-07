# https://github.com/mearaj/bhagad-house-booking

REPO_NAME=.

POSTGRES_USER=root
POSTGRES_PASSWORD=secret

print:


dep-tools:
	# https://github.com/kyleconroy/sqlc
	brew install sqlc

	# https://github.com/golang-migrate/migrate
	brew install golang-migrate

	# gio cmd
	go install gioui.org/cmd/gogio@latest

mod-fix:
	# copy the giowidgets in into correct file system location. See the go-mod !
	rm -rf giowidgets
	git clone git@github.com:mearaj/giowidgets.git
	@echo giowidgets >> .gitignore

	rm -rf $(REPO_NAME)/frontend/ui/view/third-party
	mkdir -p $(REPO_NAME)/frontend/ui/view/third-party
	cp -r ./giowidgets $(REPO_NAME)/frontend/ui/view/third-party/giowidgets

mod-upgrade:
	# go mod update
	go install github.com/oligot/go-mod-upgrade@latest
	cd ./giowidgets && go-mod-upgrade
	cd ./giowidgets && go mod tidy
	cd ./$(REPO_NAME) && go-mod-upgrade
	cd ./$(REPO_NAME) && go mod tidy

FRONT_OUT=$(REPO_NAME)/frontend/output
FRONT_OUT_ABS=$(PWD)/$(REPO_NAME)/frontend/output
front-build-mac:
	# mac
	cd $(REPO_NAME)/frontend/cmd && go build -o $(FRONT_OUT_ABS)/mac/booking .
front-build-win:
	# mac
	cd $(REPO_NAME)/frontend/cmd && CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o $(FRONT_OUT_ABS)/win/booking.exe .
front-build-wasm:
	# wasm
	#fails because frontend and backend are intertwined...
	cd $(REPO_NAME)/frontend/cmd && gogio -target js -o $(FRONT_OUT_ABS)/wasm .
front-build-delete:
	cd $(REPO_NAME)/frontend/cmd && rm -f $(FRONT_OUT_ABS)
front-run-mac:
	cd $(FRONT_OUT_ABS)/mac && ./booking
front-run-win:
	cd $(FRONT_OUT_ABS)/mac && ./booking.exe
front-run-wasm:
	go run .

# postgresql://root:secret@localhost:5432/bhagad_house_booking?sslmode=disable"
# MAKE sure postres NOT running on desktop, using the same port !!
back-gen:
	cd $(REPO_NAME)/common/ && $(MAKE) sqlc
back-db-run:
	cd $(REPO_NAME)/common/ && $(MAKE) postgres
back-db-create:
	cd $(REPO_NAME)/common/ && $(MAKE) createdb
back-db-drop:
	cd $(REPO_NAME)/common/ && $(MAKE) dropdb
back-db-migrate-up:
	cd $(REPO_NAME)/common/ && $(MAKE) migrateup
back-db-migrate-down:
	cd $(REPO_NAME)/common/ && $(MAKE) migratedown
back-db-test:
	# fills db with data..
	cd $(REPO_NAME)/common/ && $(MAKE) test


### pgweb
# Web DB inspector

pgweb-run:
	# db web gui. light and easy
	# https://github.com/sosedoff/pgweb
	go install github.com/sosedoff/pgweb@latest
	pgweb --url postgresql://root:secret@localhost:5432/bhagad_house_booking?sslmode=disable


### dblab
# TUI db inspector

dblab-run:
	# still very beta. Use pgweb.
	go install github.com/danvergara/dblab@latest
	dblab --url postgresql://root:secret@localhost:5432/bhagad_house_booking?sslmode=disable