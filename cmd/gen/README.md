# Gen

## Reference

[Database To Structs](https://gorm.io/gen/database_to_structs.html)

## Instructions

### 1. After creating the database table(s), use GORM Gen to generate the corresponding model and query files

### 2. Open ./main.go, and on line 29, modify the following code

```go
g.GenerateModel("users")
```

Replace "users" with the name of the table you want to generate.

### 3. Run the generator in the current directory

```go
go run .
```

### 4. This operation will overwrite the file: ./query/gen.go

After generation is complete:

- Copy the contents of the newly generated gen.go.

- Restore the original gen.go file if necessary.

- Paste the copied content back into the appropriate location.

> ⚠️ Note: Always backup any manually written code before running the generation command, as it will be overwritten.

### 5. Example gen struct and query from table 

```
docker exec -it backend-postgres-1 psql -U postgres -d Thingsly
```

```
CREATE TABLE test_table (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), name VARCHAR(255) NOT NULL, description TEXT, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);
```