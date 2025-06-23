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
	go test -v -count=1 ./test -run TestArticleIntegrationTestSuite

.PHONY: test-cleanup
test-cleanup:
	docker-compose stop postgres_test
	docker-compose rm -f postgres_test
	docker volume prune -f

.PHONY: integration-test
integration-test: vendor test-setup test-run test-cleanup

.PHONY: migrate
migrate:
	migrate -path ./migrations -database "postgres://user_articles_feed:pass_articles_feed@localhost:5432/articles_feed?sslmode=disable" up

.PHONY: api-run
api-run: vendor
	docker-compose up --build -d api postgres
	@echo "Waiting for database to be ready..."
	@sleep 3
	$(MAKE) migrate

.PHONY: api-stop
api-stop:
	docker-compose down -v
	docker image prune -f
	docker volume prune -f
