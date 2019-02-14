BIN=pomos

build:
	cd cmd && go build -o $(BIN)

run: build
	./cmd/pomos -c ./cmd/config.yaml
