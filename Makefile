.PHONY: build templ tailwind clean run docker-build docker-run

IMAGE ?= stylesheets

TEMPL := $(shell go env GOPATH)/bin/templ
TAILWIND := ./node_modules/.bin/tailwindcss

# Generate templ files, build Tailwind, then compile Go binary
build: templ tailwind
	go build -o ./bin/stylesheets .

# Run templ code generator
templ:
	$(TEMPL) generate

# Build Tailwind CSS (minified for production)
tailwind:
	$(TAILWIND) -i static/css/input.css -o static/css/output.css --minify

# Watch mode — run each in a separate terminal during development
watch-templ:
	$(TEMPL) generate --watch

watch-tailwind:
	$(TAILWIND) -i static/css/input.css -o static/css/output.css --watch

# Build and run the server
run: build
	./bin/stylesheets

# Build Docker image
docker-build:
	docker build -t $(IMAGE) .

# Run the Docker image locally (builds first if needed)
docker-run: docker-build
	docker run --rm -p 8080:8080 $(IMAGE)

# Clean generated artifacts
clean:
	rm -rf ./bin
	rm -f static/css/output.css
	find . -name "*_templ.go" -delete
