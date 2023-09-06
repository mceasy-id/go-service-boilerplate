setup:
	docker compose -f ./infrastructure/docker-compose.dev.yml up -d
	docker exec db-migrator alembic upgrade head

down: 
	docker compose -f ./infrastructure/docker-compose.dev.yml down -t 1

restart:
	docker compose down -t 1
	docker compose up -d
	docker exec alembic-migrator alembic upgrade head

mock-repo: 
	mockgen -source internal/$(domain)/repository.go -destination internal/$(domain)/mock/repository_mock.go -package=mocks

mock-usecase:  
	mockgen -source internal/$(domain)/usecase.go -destination internal/$(domain)/mock/usecase_mock.go -package=mocks

mock-httpclient:
	mockgen -source internal/$(domain)/http_client.go -destination internal/$(domain)/mock/httpclient_mock.go -package=mocks

migrate-up: 
	docker exec db-migrator alembic upgrade head

migrate-revision:
	docker exec db-migrator alembic revision --autogenerate -m "$(message)"

migrate-down: 
	docker exec db-migrator alembic downgrade -1

integration-test:
	docker compose -f ./infrastructure/docker-compose.dev.yml down -t 1
	docker compose -f ./infrastructure/docker-compose.dev.yml up -d
	docker exec db-migrator alembic upgrade head
	go clean -testcache
	go test ./...
	docker compose -f ./infrastructure/docker-compose.dev.yml down -t 1

run-w-observability:
	docker compose up -d
	docker exec db-migrator alembic upgrade head
	go run cmd/api/main.go