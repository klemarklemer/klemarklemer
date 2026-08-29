.PHONY: dev up deps api web compose-up compose-down migrate test-api typecheck-web lint-web build-web test-web test predeploy

# up runs the API and the web app together in one terminal.
#
# `trap 'kill 0'` sends the signal to the whole process group, so Ctrl-C stops
# both children rather than orphaning one of them. `wait` keeps make attached
# until they exit.
up:
	@echo "API  http://localhost:8000"
	@echo "Web  http://localhost:5173"
	@echo "Ctrl-C stops both. Dependencies: make deps"
	@trap 'kill 0' EXIT INT TERM; \
		( cd apps/api && candi -run -service core ) & \
		( cd apps/web && npm run dev ) & \
		wait

dev: up

deps:
	docker compose -f deployments/compose/docker-compose.dev.yml up -d postgres redis

api:
	cd apps/api && candi -run -service core

web:
	cd apps/web && npm run dev

compose-up:
	docker compose -f deployments/compose/docker-compose.dev.yml up --build

compose-down:
	docker compose -f deployments/compose/docker-compose.dev.yml down

migrate:
	cd apps/api && make migration service=core

test-api:
	cd apps/api && go test ./services/core/...

typecheck-web:
	cd apps/web && npm run typecheck

lint-web:
	@echo "Linting web..."

build-web:
	cd apps/web && npm run build

test-web: typecheck-web build-web

test: test-api test-web

predeploy: test
