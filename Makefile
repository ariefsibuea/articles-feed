.PHONY: vendor
vendor: go.mod go.sum
	@echo "Tidying and vendoring Go modules..."
	go mod tidy
	go mod vendor


.PHONY: test-setup
test-setup:
	docker-compose up -d postgres_test
	@echo "Waiting for test database to be ready..."
	@sleep 3

.PHONY: test-run
test-run:
	TEST_DB_HOST=localhost \
	TEST_DB_PORT=5433 \
	TEST_DB_USER=user_articles_feed_test \
	TEST_DB_PASSWORD=pass_articles_feed_test \
	TEST_DB_NAME=articles_feed_test \
	go test -v -count=1 ./test

.PHONY: test-cleanup
test-cleanup:
	docker-compose stop postgres_test
	docker-compose rm -f postgres_test
	docker volume prune -f

.PHONY: test-integration
test-integration:
	@$(MAKE) vendor
	@$(MAKE) test-setup
	-@$(MAKE) test-run || true
	@$(MAKE) test-cleanup

.PHONY: migrate
migrate:
	migrate -path ./migrations -database "postgres://user_articles_feed:pass_articles_feed@localhost:5432/articles_feed?sslmode=disable" up

.PHONY: api-start
api-start: vendor
	docker-compose up --build -d api postgres
	@echo "Waiting for database to be ready..."
	@sleep 3
	$(MAKE) migrate

.PHONY: api-stop
api-stop:
	docker-compose down -v
	docker image prune -f
	docker volume prune -f
