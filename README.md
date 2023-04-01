# Example API in golang (no framework)

This is a banking application API, create account, delete account, list accounts and transfer money

!!! This project is not to be used in any production environment, only a tutorial example. !!!
*** Still work in progress ***

## Contains
- PostGresql database support
- SqLite database support
- JWT token authentication
- Middlevare
- CORS headers

# Commands

```
make build	(build application to bin folder)
make run (build and run)
make run-sqlite (build and run with sqlite database)
make seed (buld and seed a user into the databse)
make seed-sqlite (build and seed a user to sqilte database)
make serve (serve locally)
make serve-sqlite (serve localy with sqlite db)
make test (run test)
make deploy (not to be used)
```

# Parameters
```
go run . -<parameter>

-seed (seeds new user)
-sqlite (uses the sqlite db)
-port (serve on port ex port=80) / by default it serves on port 8000
```

## JWT secret
When generates JWT token it uses a secret which can be passed by variable
```
export JWT_SECRET="bummm"
```

## Docker
The docker-compose.yml file only intended to spin up a postgres DB, the go app is not yet dockerised

# Happy learning



// TODO add validation
// Secure application
// TODO full test coverage
