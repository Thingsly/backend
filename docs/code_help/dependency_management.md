# Dependency Management

## Common Commands

1. `go mod tidy -compat=1.24`
   Specifies the Go version compatibility to consider when running `go mod tidy`.

2. `go list -m -versions github.com/xxx`
   View the version list of a package.

3. Use `go get xxx` in the terminal to download packages. Do not click to download packages within the code (this may lead to downloading the latest package version, which might be incompatible with the required version).
