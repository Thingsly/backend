# Unit Testing

## Unit Test File Naming Convention

- Unit test files should follow the naming convention `xxx_test.go`.
- Unit test files are typically stored in the `./test` directory.
- It is also acceptable to place test files in the same directory as the module being tested for easier proximity, for example: `./api/test/board_api_test.go`.

## Test Functions

- Test function names should follow the naming convention `TestXXX`.
- The first parameter of the test function must be `t *testing.T`.
- The first line of the test function must initialize the `require` object using the following code:

  ```go
  require := require.New(t)   // Initialize the require library


Example:

```go
def test_add(t *testing.T):
    require.Equal(123, 123, "they should be equal")
```

## Running Tests

- To run the tests, use the command:：`run_env=localdev go test -v ./...`
- The go test command will automatically find all files ending with *_test.go in the ./test directory and execute the test functions.
- If a test function contains t.Error() or t.Fail(), the test will fail.
- If the test function does not contain t.Error() or t.Fail(), the test will pass.

- To run all unit tests in the project locally, use the following command: `run_env=localdev go test -v ./...`
- Before running the tests, make sure to update the database configuration in the configs/conf-localdev.yaml file and ensure the local environment has the database running.
  