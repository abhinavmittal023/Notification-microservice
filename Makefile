define GetFromConfig
$(shell node -p "require('./configuration/config.json').$(1)")
endef

all: 
	@echo "Run (make setup) for first run or (make run) to run the server"

down: 
	@echo "Migrating Down"
	./goose -dir ./db/migrations/ postgres "user=$(call GetFromConfig,database.user) password=$(call GetFromConfig,database.password) dbname=$(call GetFromConfig,database.dbname) sslmode=$(call GetFromConfig,database.sslmode)" down

up:
	@echo "Migrating Up"
	./goose -dir ./db/migrations/ postgres "user=$(call GetFromConfig,database.user) password=$(call GetFromConfig,database.password) dbname=$(call GetFromConfig,database.dbname) sslmode=$(call GetFromConfig,database.sslmode)" up

build_goose:
	@echo "Building goose binary"
	go build -o goose db/cmd/main.go

build_main:
	@echo "Building main.go"
	go build main.go

run:
	go run main.go

setup: build_goose up run

refresh_database: down setup
