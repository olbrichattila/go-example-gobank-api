module example.com/storage

go 1.18

replace example.com/types => ../types

require (
	example.com/types v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.10.7
	github.com/mattn/go-sqlite3 v1.14.16
)

require golang.org/x/crypto v0.7.0 // indirect
