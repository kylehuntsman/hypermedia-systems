.PHONY: run dev build prepare clean

# Runs the application
run: public/css/style.css
	go run .

# Runs the application with hot reloading
dev: bin/main
	go run github.com/cosmtrek/air@v1.47.0

## Builds the application
build: bin/main

# Runs any preparation steps
prepare: public/css/style.css

# Removes any generated files and directories
clean:
	rm -rf bin tmp public/css/style.css

bin/main: *.go contact/*.go public/css/style.css
	go build -o bin/main .

public/css/style.css: templates/*.html
	tailwindcss -o public/css/style.css