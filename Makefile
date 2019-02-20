BIN=pomos

build:
	cd cmd/pomos && go build -o $(BIN)

run: build
	./cmd/pomos/$(BIN) -c ./config.yaml
