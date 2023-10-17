.PHONY: run dev build prepare clean

# Runs the application
run: templates/*_templ.go
	go run main.go

# Runs the application with hot reloading
dev: bin/main
	go run github.com/cosmtrek/air@v1.47.0

## Builds the application
build: bin/main

# Runs any preparation steps, in this case generating the templated files
prepare: templates/*_templ.go

# Removes any generated files and directories
clean:
	rm -rf bin tmp templates/*_templ.go

bin/main: templates/*_templ.go
	go build -o bin/main main.go

templates/%_templ.go: templates/%.templ
	go run github.com/a-h/templ/cmd/templ@v0.2.408 generate