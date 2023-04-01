build:
	@export JWT_SECRET="bummm"
	@go build -o bin/gobank
run: build
	@./bin/gobank
run-sqlite: build
	@./bin/gobank -sqlite
seed: build
	@./bin/gobank -seed
seed-sqlite: build
	@./bin/gobank -sqlite -seed
serve:
	@export JWT_SECRET="bummm"
	@go run .
serve-sqlite:
	@export JWT_SECRET="bummm"
	@go run . -sqlite
test:
	@go test -v ./...
deploy:
	scp bin/gobank root@104.248.243.214:/root/gobank
