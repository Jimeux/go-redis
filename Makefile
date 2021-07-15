run:
	go run main.go

bench:
	go test ./cache -bench=. -benchmem

init-db:
	go run dynamo.go

redis:
	redis-cli -h 127.0.0.1 -p 26379
