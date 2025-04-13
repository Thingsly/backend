# Notes

- Do not modify the contents of *.gen.go files to avoid changes being overwritten when regenerated.

- When writing interfaces, add *.http.go files for validation.

- Required fields should not be pointers, and optional fields should be pointers with `omitempty`.

- The struct tags for GET requests must include the `form` tag for fields, e.g., `form:"page"`.
