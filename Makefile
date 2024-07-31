include .env



.PHONY: build
build:
	@templ generate 
	@go build -o bin/api 
	
	
	
	

run: build
	@./bin/api
test:
	@go test -v ./...



templWatch:
	@templ generate --watch --proxy="http://localhost:8000" --open-browser=false


templ:
	@templ generate


status:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_DSN) goose -dir=./db/migrations status
	@echo "Database Status"


migrate:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_DSN) goose -dir=./db/migrations up
	@echo "Database Migrated"


migrate-down:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_DSN) goose -dir=./db/migrations down
	@echo "Database Migrated Down"

drop:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_DSN) goose -dir=./db/migrations drop
	@echo "Database Dropped"


tailwind:
	@npx tailwindcss -i ./static/css/input.css -o ./static/css/styles.css --watch





tailwind-build:
	@npx tailwindcss -i ./views/css/input.css -o ./views/css/style.min.css --minify