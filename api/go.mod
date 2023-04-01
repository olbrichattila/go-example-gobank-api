module example.com/api

go 1.18

replace example.com/storage => ../storage

replace example.com/types => ../types

require (
	example.com/storage v0.0.0-00010101000000-000000000000
	example.com/types v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.0.0-rc.1
	github.com/gorilla/mux v1.8.0
	golang.org/x/crypto v0.7.0
)

require (
	github.com/lib/pq v1.10.7 // indirect
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
)
