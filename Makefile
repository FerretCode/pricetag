include .env

## help: print this help message
.PHONY: help
help:
	@echo "Usage:"
	@sed -n "s/^##//p" ${MAKEFILE_LIST} | column -t -s ":" |  sed -e "s/^/ /"

# confirmation dialog helper
.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

## css: build daemon for TailwindCSS
.PHONY: css
css:
	tailwindcss -i ./ui/input.css -o ./ui/static/main.css --watch

## css-minify: build CSS for production
.PHONY: css-minify
css-minify:
	tailwindcss -i ./ui/input.css -o ./ui/static/main.css --minify


## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo "Tidying and verifying module dependencies..."
	go mod tidy
	go mod verify
	@echo "Formatting code..."
	go fmt ./...
	@echo "Vetting code..."
	go vet ./...
	staticcheck ./...
	@echo "Running tests..."
	go test -race -vet=off ./...
	
## build: build the cmd/pricetag application for production
.PHONY: build
build: css-minify
	@echo "Building cmd/pricetag..."
	go build -ldflags="-s" -o=./bin/pricetag ./cmd/pricetag

## run: run the cmd/pricetag application
.PHONY: run
run:
	go run ./cmd/pricetag -port=${PORT} -dev \
		-db-dsn=${DATABASE_URL}

## dev: run cmd/pricetag and CSS daemon
.PHONY: dev
dev:
	${MAKE} -j2 css run
