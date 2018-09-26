
.PHONY: all
all:
	go get -v ./...
	go build -v .

.PHONY: generate
generate:
	swagger generate server -A flamedb -P models.Principal -f ./swagger.yaml

.PHONY: install_swagger
install_swagger:
	brew tap go-swagger/go-swagger
	brew install go-swagger
