data "external_schema" "gorm" {
  program = [
    "go", "run", "./cmd/atlas/main.go",
  ]
}

env "local" {
  src = data.external_schema.gorm.url
  dev = "postgres://${getenv("DB_USER")}:${getenv("DB_PASSWORD")}@atlas-dev-db:5432/dev?sslmode=disable"
  url = "postgres://${getenv("DB_USER")}:${getenv("DB_PASSWORD")}@${getenv("DB_HOST")}:5432/${getenv("DB_NAME")}?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
}

env "prod" {
  url = "postgres://${getenv("DB_USER")}:${getenv("DB_PASSWORD")}@${getenv("DB_HOST")}:5432/${getenv("DB_NAME")}?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
}
