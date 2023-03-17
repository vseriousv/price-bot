prebuild:
	rm -rf build

build: prebuild
	go build -o ./build/price_alerts_bot -v ./cmd/price_alerts_bot/main.go

dev:
	go run ./cmd/price_alerts_bot/main.go

worker-dev:
	go run ./cmd/worker/price_monitoring.go

prod: build
	./build/main

.DEFAULT_GOAL := build
