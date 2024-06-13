default: help

help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

dev: ## runs the application in development mode
	@docker compose -f dev.docker-compose.yaml up


create-certs: ## generate self-signed certificates
	mkcert -install
	mkcert -cert-file ./certs/local-cert.crt \
       -key-file ./certs/local-cert.key \
       local.freelance-invoice-hub.com localhost 127.0.0.1 ::1

