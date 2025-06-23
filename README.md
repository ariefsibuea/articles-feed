# Articles Feed

Articles Feed is a simple RESTful API responsible for creating and fetching articles.

## Features

- Create new articles
- Fetch a list of articles

## Prerequisites

- Go installed
- Docker, if you want to run dependencies in containers

## Getting Started

Clone the repository and install dependencies:

```bash
git clone https://github.com/ariefsibuea/articles-feed.git
cd articles-feed
make vendor
```

## Running the API

To start and stop the API, use the following commands:

```bash
# Run the API server
make api-run

# Stop and cleanup resources
make api-stop
```

Remember to prepare the `.env` file before running the API. You can use the provided sample as a starting point. By default, the API will be available at `http://localhost:8080`.

## API Documentation

### Create Article

- **Endpoint:** `POST /articles`
- **Description:** Create a new article.
- **Request Body:**

    ```json
    {
        "title": "Article Title",
        "body": "Article content goes here.",
        "authorName": "John Doe"
    }
    ```

- **Response:**
  - **201 Created**

        ```json
        {
            "success": true,
            "data": {
                "id": "acdb113a-60ae-4643-92c7-2d15f675b3f5",
                "title": "Article Title",
                "body": "Article content goes here.",
                "createdAt": "2025-06-23T11:14:55Z",
                "authorName": "John Doe"
            },
            "meta": {}
        }
        ```

  - **400 Bad Request:** Invalid input.
  - **500 Internal Server Error:** Internal server error.

### Fetch Articles

- **Endpoint:** `GET /articles`
- **Description:** Retrieve a list of articles.
- **Response:**
  - **200 OK**

        ```json
        {
            "success": true,
            "data": {
                "articles": [
                    {
                        "id": "acdb113a-60ae-4643-92c7-2d15f675b3f5",
                        "title": "Article Title",
                        "body": "Article content goes here.",
                        "created_at": "2025-06-23T11:14:55Z",
                        "authorName": "John Doe"
                    }
                ]
            },
            "meta": {
                "page": 1,
                "pageSize": 10,
                "totalItems": 1
            }
        }
        ```

  - **500 Internal Server Error:** Internal server error.

## Testing

This project includes integration tests. To run them, use:

```bash
# Prepare dependencies
make vendor

# Set up test environment
make test-setup

# Run tests
make test-run

# Cleanup test resources
make test-cleanup

# Run all integration tests (covers all commands above)
make integration-test
```

Before running the tests, ensure that the `.env.test` file is present. You can use the provided sample file.
