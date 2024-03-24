# its-gram

Small instagram clone "its-gram" backend for "Layanan dan Aplikasi Internet" college class project.

## How to run with docker

### Prerequisites

-   Firebase service account key (used for Firebase Auth)

### Steps

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

### How to configure swagger

1. Configure the swagger main api info (see comments on top of `func main()` in `main.go`), then add the info at your handlers too.

2. Then first try to do a `swag init -g cmd/main.go` to generate the swagger docs package.

3. Add `docs` package to `main.go`.

```
import (
    _ "github.com/rrab-0/its-gram/docs"
)
```

4. Add handler for swagger at `router.go` so we can access the docs at `/swagger/index.html`.

```
import (
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter() {
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

3. Then do `swag init`

```
swag init -g cmd/main.go --parseDependency --parseInternal
```

need `--parseDependency` and `--parseInternal` so swagger will parse dependencies and internal packages.
