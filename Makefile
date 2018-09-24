
.PHONY: all
all:
	go build -v ./cmd/flamedb-server

.PHONY: generate
generate:
	swagger generate server -A flamedb -f ./swagger.yaml

.PHONY: install_swagger
install_swagger:
	brew tap go-swagger/go-swagger
	brew install go-swagger
