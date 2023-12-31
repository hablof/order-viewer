### генерация моков
mock-for-httpcontroller:
	minimock -i ./internal/httpcontroller.* -o ./internal/httpcontroller

### запуски
run-pub:
	go run cmd/producer/main.go

run-app:
	go run cmd/order-viewer/main.go

run-stan:
	./nats-streaming-server.exe -sc streaming.config