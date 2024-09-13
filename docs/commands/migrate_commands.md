### Working with golang-migrate SQL MIgrations

- **Creating a new up and down migrations with migrate**:
```bash
migrate create -seq -ext=.sql -dir=./migrations create_categories_table
```
/home/kcharymyrat/dev/go/e-commerce/migrations/000001_create_categories_table.up.sql
/home/kcharymyrat/dev/go/e-commerce/migrations/000001_create_categories_table.down.sql

- -seq flag indicates that we want to use sequential numbering like 0001, 0002, ...
- -ext flag indicates that we want to give the migration files the extension .sql
- -dir flag indicates that we want to store the migration files in the ./migrations

- **Executing the migrations**:
```bash
migrate -path=./migrations -database=$POSTGRES_DSN up
```

-- **Check the Migration Status: You can use the migrate tool to inspect the status of your migrations**:
```bash
migrate -path=./migrations -database="postgres://dev:dev@localhost/ecommerce" version
```

-- **Force the migration version to 1**"
```bash
migrate -path=./migrations -database="postgres://dev:dev@localhost/ecommerce" force 1
```