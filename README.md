# imdb-rest-apis

#### Prerequisite

1. [Install Golang](https://golang.org/doc/install)
2. Setup GOPATH [Link1](https://golang.org/doc/code.html#GOPATH) and [Link2](https://github.com/golang/go/wiki/GOPATH)
3. [Install Glide](https://github.com/Masterminds/glide)
4. [Setup CockrochDB] (https://www.cockroachlabs.com/docs/stable/) to store movie details in local

#### Getting Started

1. Clone the repo under `$GOPATH/src`. If folder does not exist than create it first. Then run `git clone https://github.com/Gohelraj/imdb-rest-apis.git imdb-rest-apis`
2. Run `glide install`
3. Add `config.toml` file path in `path` variable in `config/config.go` file
4. Change configuration setting of Database in `Config.toml`
4. Run `go run cmd/main.go`