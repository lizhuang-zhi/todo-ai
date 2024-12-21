.PHONY: run
run: ## Startup AI Dev Mate Server
	go run main.go -c=./configs/local/config.yaml

.PHONY: prod
prod: ## Production Startup AI Dev Mate Server
	go run main.go -c=./configs/prod/config.yaml
