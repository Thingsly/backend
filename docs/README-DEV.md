# Project Documentation

## Gen Tool (Must Read)

**Model & Query Code Generation**: [Refer to the README](/cmd/gen/REANME.md)

## Running & Debugging

- **Run**: `go run .`
- **API Debugging**: [http://localhost:9999/swagger/index.html](http://localhost:9999/swagger/index.html)

## Important Notes

- Logging is handled using `logrus`.
- Use `viper` to load config files. Always set default values in case a config is missing.
- Both `model` and `query` should be generated using the Gen tool.
- When writing code for the first time, refer to interfaces related to user management and notification groups.
- Pay special attention to authorization logic:
  - There are three user types distinguished by the `authority` field:
    - `SYS_ADMIN`: System Administrator
    - `TENANT_ADMIN`: Tenant Administrator
    - `TENANT_USER`: Tenant User
  - Each tenant has one Tenant Admin. Data is tenant-isolated using the `tenant_id` field.
  - Default system admin account: `demo-super@thingsly.vn` / `123456`
  - SYS_ADMIN can create TENANT_ADMIN accounts. TENANT_ADMIN can create TENANT_USER accounts.

- All newly added SQL statements should be updated in `/sql/1.sql`.
- There are helper methods in `service/enter.go` for pointer conversions.
- Database error messages can be unclear, so JSON field formats must be validated manually, for example:

```go
// Validate if input is a JSON string
if CreateDeviceConfigReq.AdditionalInfo != nil && !IsJSON(*CreateDeviceConfigReq.AdditionalInfo) {
  return fmt.Errorf("additional_info is not a valid JSON")
}
deviceconfig.AdditionalInfo = CreateDeviceConfigReq.AdditionalInfo
```

Or validate during map conversion:

```go
condsMap, err := StructToMapAndVerifyJson(req, "additional_info", "protocol_config")
if err != nil {
  return nil, err
}
deviceconfig.AdditionalInfo = CreateDeviceConfigReq.AdditionalInfo
```

- `First()` returns an error if no record is found. Use `errors.Is(err, gorm.ErrRecordNotFound)` to check.
- See `initialize/caching-overview.md` for caching documentation.

---

## Development Notes & Conventions

### Basic RESTful Route Naming & Methods

- **Create Resource**
  - Path: `POST /resources`
  - Use `POST` without `/create` suffix; pass new resource data in request body.

- **Get Resource Details**
  - Path: `GET /resources/{id}`
  - Use `GET` with resource ID as a path parameter.

- **Update Resource**
  - Path: `PUT /resources/{id}`
  - Use `PUT` with updated data in the request body.

- **Delete Resource**
  - Path: `DELETE /resources/{id}`
  - Use `DELETE` with the resource ID in the path.

- **Get Resource List (with filters/sorting)**
  - Path: `GET /resources`
  - Use `GET` with query parameters.

### Request Data Transmission

- **Query Parameters**: Used in `GET` for filtering, pagination, sorting, e.g.,  
  `GET /resources?role=admin&page=2&limit=10`

- **Path Parameters**: For identifying specific resources in `GET`, `PUT`, `DELETE` requests.  
  Example: `DELETE /resources/{id}`

- **Request Body**: Used in `POST` and `PUT` for transmitting complex data, typically in JSON format.

### API Responses

- `POST` should return a complete representation of the newly created resource, including server-generated fields like ID and timestamps.
- `PUT` should return the updated resource in full to confirm successful changes.

### Complex Queries & Raw SQL

- For advanced querying scenarios, prefer raw SQL for maintainability and clarity.  
  Refer to `dal/ota_upgrade_tasks.go -> GetOtaUpgradeTaskListByPage` for implementation details.

---

## Git Commit Message Convention

> All commit messages should be written in English.

- `Feat`: New feature
- `Fix`: Bug fix
- `Docs`: Documentation changes
- `Style`: Code formatting (no logic changes)
- `Refactor`: Code refactoring without functional changes
- `Test`: Adding or updating tests
- `Chore`: Changes to build tools or infrastructure

---

## Go Programming Philosophy

- **Simplicity**  
  > "Less is exponentially more." — Rob Pike  
  Reducing language features increases expressiveness and flexibility.

- **Concurrency over Parallelism**  
  > "Concurrency is not parallelism." — Rob Pike  
  Go's concurrency model (goroutines & channels) is about handling multiple tasks efficiently, not just multithreading.

- **Clear Error Handling**  
  > "Clear is better than clever."  
  Go avoids exception-based logic in favor of explicit error checks.

- **Unified Code Formatting**  
  > "Gofmt's style is no one's favorite, yet gofmt is everyone's favorite."  
  Code formatting is enforced via `gofmt` to ensure consistency.

- **Avoid Hidden Complexity**  
  > "Do less. Enable more."  
  Prefer straightforward implementations over overly abstract solutions.

- **Small Interfaces Are Better**  
  > "The bigger the interface, the weaker the abstraction."  
  Design focused, minimal interfaces that do one thing well.

---

## Swagger Guide

```bash
# Run from the project root
swag init
```

- Swagger UI: [http://localhost:9999/swagger/index.html](http://localhost:9999/swagger/index.html)
- Writing Swagger annotations:  
  Refer to the `/api/v1/login [post]` endpoint for examples.
  