define GetFromConfig
$(shell node -p "require('./configuration/config.json').$(1)")
endef

.PHONY: all

all: 
	@echo "Run (make setup) for first run or (make run) to run the server"

down: 
	@echo "Migrating Down"
	./goose -dir ./db/migrations/ postgres "$(call GetFromConfig,database.dbstring)" down

down_all:
	@echo "Migrating Down to Empty"
	./goose -dir ./db/migrations/ postgres "$(call GetFromConfig,database.dbstring)" down-to 0
up:
	@echo "Migrating Up"
	./goose -dir ./db/migrations/ postgres "$(call GetFromConfig,database.dbstring)" up

build_goose:
	@echo "Building goose binary"
	go build -o goose db/cmd/main.go

database_version:
	./goose -dir ./db/migrations/ postgres "$(call GetFromConfig,database.dbstring)" version

build_main:
	@echo "Building main.go"
	go build main.go

golangci-lint:
	golangci-lint run

run:
	go run main.go

setup: build_goose up run

refresh_database: down_all setup

run_mail_catcher:
	@echo "Go to http://localhost:1080/ on your browser to quit or monitor"
	mailcatcher --ip 0.0.0.0

run_controller_test:
	@echo "Running all unit tests for controllers" 
	go test -count=1 ./tests/controllers/... 

run_services_test:
	@echo "Running all unit tests for services"
	go test -count=1 ./tests/services/...

run_controller_test_verbose:
	@echo "Running all unit tests for controllers with verbose" 
	go test -count=1 ./tests/controllers/... -v

run_services_test_verbose:
	@echo "Running all unit tests for services with verbose"
	go test -count=1 ./tests/services/... -v

run_unit_tests_verbose: run_services_test_verbose run_controller_test_verbose 

run_unit_tests: run_services_test run_controller_test

run_benchmark:
	@echo "Running Benchmark"
	go test  ./tests/benchmark/ -bench=.
