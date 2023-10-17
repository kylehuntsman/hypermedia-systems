# Hypermedia Systems
Following the web application described in the https://hypermedia.systems book using Golang.

> For web developers frustrated with the complexity of modern practice, those looking to brush up on web fundamentals, web development shops looking to bring their apps to mobile, and any workaday programmer looking for an introduction to hypermedia and REST.

## Technology Stack
- [Golang](https://go.dev/) as the programming language
- [Fiber](https://docs.gofiber.io/) as the web framework
- [Templ](https://templ.guide/) as the server-side templating engine
- [Air](https://github.com/cosmtrek/air) for live reloading
- [HTMX](https://htmx.org/) for the client side

## Getting Started

### Prerequisites
- [Golang](https://golang.org/doc/install) v1.20 or higher

### Running
Running the application via the `make` or `make run` command will run the application, generating any missing template files in the process.
```bash
$ make
```

## Development

### Running
When running the application in development mode, the project uses the [Air](https://github.com/cosmtrek/air) library for live reloading and is configured to restart the server when any `.go` or `.templ` files are changed. The Air reload configuration can be found in the [`.air.toml`](./.air.toml) file. Air will restart the server by running the `make build` command and then running the resulting binary.

Run `make dev` to start the application for development.
```bash
$ make dev
```

### Building
Running the `make build` command will build the application, generating any missing template files in the process. The resulting binary will be placed in the [bin](./bin/) directory.
```bash
$ make build
```

### Templating
The project uses [Templ](https://templ.guide/) as the server-side templating engine. The templates are in the [templates](./templates/) directory and are named with the `.templ` extension. Templ will generate the corresponding `*_templ.go` files in the [templates](./templates/) directory during the prepare step of the build process.

Running the `make prepare` command will generate any missing template files using Templ. The generation command is run automatically when running or building the application via `make`, `make run`, `make dev`, or `make build`.
```bash
$ make prepare
```

### Cleaning
Running the `make clean` command will remove any generated files, including the `bin` directory and all generated template files.
```bash
$ make clean
```