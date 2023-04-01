module example.com/gobank

go 1.18

require (
	example.com/api v0.0.0-00010101000000-000000000000
	example.com/storage v0.0.0-00010101000000-000000000000
	example.com/types v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.0.0-rc.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace example.com/types => ./types

replace example.com/storage => ./storage

replace example.com/api => ./api
