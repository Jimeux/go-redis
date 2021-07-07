run:
	go run main.go
bench:
	go test ./cache -bench=. -benchmem

redis:
	redis-cli -h 127.0.0.1 -p 26379
