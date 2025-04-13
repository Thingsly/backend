# Notes

- Do not use `gin`'s `ShouldBindJSON` for validation. Always use the `ValidateStruct` function for validation instead.

- `ShouldBindJSON` does not handle pointer types properly and is not specifically designed for struct validation.

- Use the `"github.com/go-playground/validator/v10"` package for validation.
