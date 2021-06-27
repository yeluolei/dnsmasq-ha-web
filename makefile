.PHONY: build

frontend-install:
	cd frontend && yarn install

frontend:
	cd frontend && yarn build

build:
	CGO_ENABLED=1 go build -o bin/dnsmasq-ha-web cmd/dnsmasq-ha-web/*.go