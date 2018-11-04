# imdb-rest-apis

> Apis to search movie details from local DB and  IMDB

#### Prerequisite

1. [Install Golang](https://golang.org/doc/install)
2. Setup GOPATH [Link1](https://golang.org/doc/code.html#GOPATH) and [Link2](https://github.com/golang/go/wiki/GOPATH)
3. [Install Glide](https://github.com/Masterminds/glide)
4. [Setup CockrochDB](https://www.cockroachlabs.com/docs/stable/) to store movie details in local

#### Getting Started

1. Clone the repo under `$GOPATH/src`. If folder does not exist than create it first. Then run `git clone https://github.com/Gohelraj/imdb-rest-apis.git imdb-rest-apis`
2. To install dependencies Run `glide install`
3. Add `config.toml` file path in `path` variable of `config/config.go` file
4. Change configuration setting of Database in `config.toml`
5. Add `omdb` API key in `config.toml` file. If you don't have API key than first [Generate key](http://www.omdbapi.com/apikey.aspx) and than add.
6. Run `go run cmd/main.go`
7. Import `imdb-movies-apis.postman_collection.json` file in [Postman](https://www.getpostman.com/download?platform=win64) and read the description of all apis to understand how all apis will work