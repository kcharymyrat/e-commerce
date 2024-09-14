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

-- **Migrate up or down to a specific version by using the goto command**:
```bash
migrate -path=./migrations -database=$EXAMPLE_DSN goto 1
```

-- **The down command to roll-back by a specific number of migrations.**:
```bash
migrate -path=./migrations -database =$EXAMPLE_DSN down 1
```

-- **Force the migration version to 1**"
```bash
migrate -path=./migrations -database="postgres://dev:dev@localhost/ecommerce" force 1
```

-- **Remote migration files**:
```bash
migrate -source="s3://<bucket>/<path>" -database=$EXAMPLE_DSN up
migrate -source="github://owner/repo/path#ref" -database=$EXAMPLE_DSN up
migrate -source="github://user:personal-access-token@owner/repo/path#ref" -database=$EXAMPLE_DSN up
```