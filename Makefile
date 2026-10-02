.PHONY: backend frontend test build

backend:   ## Inicia a API Go em :8080
	cd backend && go run ./cmd/server

frontend:  ## Inicia o Angular em :4200 (com proxy para a API)
	cd frontend && npm start

test:      ## Roda os testes do backend e do frontend
	cd backend && go test ./...
	cd frontend && npx ng test --watch=false

build:     ## Gera a imagem Docker completa
	docker build -t ccb .
