# its-gram

Small instagram clone "its-gram" backend for "Layanan dan Aplikasi Internet" college class project.

## How to run with docker

1. make .env, fill values from .env.example

2. get firebase service account key and put it in root folder with name `firebase-service-account.json`

3. run docker

```
docker-compose up -d
```

## How to run locally

### Prerequisites

-   Golang
-   PostgreSQL
-   make (for Makefile)
-   Firebase service account key (used for Firebase Auth)

### Steps

1. Install dependencies

```
go mod tidy
```

2. Set `ENV` in .env to be `LOCAL_DEV`

3. Run the server

```
make run
```
