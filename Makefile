define GetFromConfig
$(shell node -p "require('./configuration/config.json').$(1)")
endef

.PHONY: all

all: 
	@echo "Run (make setup) for first run or (make run) to run the server"

down: 
	@echo "Migrating Down"
	./goose -dir ./db/migrations/ postgres "$(call GetFromConfig,database.dbstring)" down

up:
	@echo "Migrating Up"
	./goose -dir ./db/migrations/ postgres "$(call GetFromConfig,database.dbstring)" up

build_goose:
	@echo "Building goose binary"
	go build -o goose db/cmd/main.go

build_main:
	@echo "Building main.go"
	go build main.go

golangci-lint:
	golangci-lint run

run:
	go run main.go

setup: build_goose up run

refresh_database: down setup

run_mail_catcher:
	@echo "Go to http://localhost:1080/ on your browser to quit or monitor"
	mailcatcher --ip 0.0.0.0
