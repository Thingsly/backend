# Naming Conventions

In Go language, following a clear set of naming conventions is crucial for maintaining the clarity and consistency of the code.

## Conventions

### 1. File Naming

- **Lowercase letters**: File names should be in lowercase letters.
- **Underscore**: If the file name consists of multiple words, use underscores (`_`) to separate them, e.g., `my_example.go`.
- **Concise and descriptive**: File names should be short and accurately describe the content.

### 2. Package Naming

- **Lowercase words**: Package names should be in lowercase letters, without underscores or camelCase.
- **Concise**: Package names should be short and meaningful.
- **Avoid redundancy**: Package names should not include the parent directory name, e.g., `package user // not recommended`.
- **Singular form**: Generally, package names should be singular, e.g., `package user` rather than `users`.

### 3. Function Naming

- **CamelCase**: Function names should use PascalCase (capitalized camel case), i.e., the first letter is uppercase, e.g., `ProcessData`.
- **Descriptive**: Function names should clearly indicate the function's purpose.

### 4. Variable Naming

- **CamelCase**: Local variables should typically use camelCase, e.g., `localVariable`.
- **Concise**: Variable names should be short yet descriptive.
- **Public variables**: If a variable needs to be accessed from outside the package, it should use PascalCase and start with an uppercase letter, e.g., `GlobalVariable`.

#### Specific Details

- **Visibility**
  - Public variables: In Go, variables starting with an uppercase letter are public (accessible from other packages).
  - Private variables: Variables starting with a lowercase letter are private (accessible only within the defining package).
- **Naming style**
  - CamelCase naming: Go uses camelCase. Local variables use lowercase camelCase (e.g., `localVariable`), while public variables use PascalCase (e.g., `PublicVariable`).
  - Short and concise: Variable names should be short while still conveying their meaning. Avoid overly long variable names unless necessary.
- **Contextual Relevance**
  - Avoid redundancy: Variable names should not contain redundant information. For example, in a `car` package, use `car.Type` instead of `car.CarType`.
  - Consider context: When naming variables, consider the context in which they are used. For example, use shorter variable names inside methods; in global contexts, use more descriptive names.
- **Special cases**
  - Temporary variables: For temporary variables, short names like `i`, `j` (loop variables), `err` (error variable), etc., are acceptable.
- **Naming conventions**
  - Boolean variables: For boolean variables, the name should be a predicate, such as `isEmpty`, `hasDone`.
  - Error variables: Error variables generally start with `err` or `Err`, such as `errNotFound`, `ErrInsufficientFunds`.
- **Use appropriate naming**
  - Avoid abbreviations: Unless the abbreviation is very common (e.g., HTTP, URL), avoid using abbreviations, as they might not be easily understood.
  - Do not use underscores or hyphens: Variable names should not contain underscores (`_`) or hyphens (`-`).

### 5. Constant Naming

- **CamelCase**: Constants are typically named using PascalCase, e.g., `ConstValue`.
- **Uppercase**: For enumeration-type constants, use all uppercase letters with underscores separating words, e.g., `CONST_ENUM_VALUE`.

### 6. Type Naming

- **CamelCase**: Custom types (including structs, interfaces, etc.) should use PascalCase, e.g., `MyStruct`.
- **Descriptive**: Type names should clearly indicate what they represent or do.

### 7. Interface Naming

- **Single-method interfaces**: If an interface contains one method, its name is typically the method name with an `-er` suffix, e.g., `Reader`.
- **Multi-method interfaces**: For interfaces with multiple methods, choose a descriptive name.

### 8. Receiver Naming

- **Short**: Receiver variable names should be short, typically the first letter of the type name in lowercase, e.g., `t` for `type MyType struct{}`.
- **Consistency**: For all methods of the same type, the receiver variable name should remain consistent.

Adhering to these naming conventions helps maintain readability and consistency in Go code and is a part of good Go programming practices.
