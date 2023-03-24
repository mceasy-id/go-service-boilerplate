setup:
	docker compose up -d
	docker exec alembic-migrator alembic upgrade head

restart:
	docker compose down -t 1
	docker compose up -d
	docker exec alembic-migrator alembic upgrade head

revision:
	docker exec alembic-migrator alembic revision --autogenerate -m "$(message)"

test:
	docker compose down -t 1
	docker compose up -d
	docker exec alembic-migrator alembic upgrade head
	go clean -testcache
	go test ./...
	docker compose down -t 1