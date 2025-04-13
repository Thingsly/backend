# Development Help

## Project Structure

```text
project-root/
├── main.go                         # Main entry point for the application
├── middleware/                     # Middleware
│   └── response/                  # Response handling
│       ├── middleware.go          # Response middleware
│       └── response.go            # Response structure definitions
├── pkg/                           # Common packages
│   ├── errcode/                  # Error code handling
│   │   ├── code.go              # Error code constants
│   │   ├── error.go             # Error definitions
│   │   ├── manager.go           # Error code manager
│   │   └── validator.go         # Error code validator
│   └── metrics/                  # Performance monitoring
│       └── metrics.go            # Metric definitions
├── config/                       # Configuration files
│   ├── config.yaml              # Main configuration file
│   └── messages.yaml            # Error message configuration
├── go.mod                        # Go module file
└── README.md                     # Project documentation
```

## Usage Example

💡 For core error codes (such as system errors, general business errors), use constants from pkg/errcode/code.go. For specific business error codes, use strings directly. The main focus should be on maintaining and ensuring the completeness of messages.yaml.

### Example: Error Handling in Services

```go
// internal/service/user.go
package service

import (
   "context"
   "project/pkg/errcode"
)

type UserService struct {
   repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
   return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*User, error) {
   // Using Newf to handle parameter error
   if id <= 0 {
       return nil, errcode.Newf(errcode.CodeParamError, id)
   }

   user, err := s.repo.GetUser(ctx, id)
   if err != nil {
       if isNotFoundError(err) {
           // Using WithVars to handle user not found error
           return nil, errcode.WithVars(errcode.CodeNotFound, map[string]interface{}{
               "id": id,
               "error": "User not found",
           })
       }
       // Using WithData to handle database error
       return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
           "sql_error": err.Error(),
           "user_id": id,
       })
   }

   return user, nil
}
```

### Example: API Handlers

```go
// internal/api/handler/user.go
package handler

import (
   "github.com/gin-gonic/gin"
   "project/internal/service"
   "project/pkg/errcode"
)

type UserHandler struct {
   userSvc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
   return &UserHandler{userSvc: svc}
}

func (h *UserHandler) GetUser(c *gin.Context) {
   id := c.GetInt64Param("id")
   
   // Parameter validation
   if id <= 0 {
       // Using NewWithMessage to handle invalid ID error
       c.Error(errcode.NewWithMessage(200015, "Please provide a valid user ID"))
       return
   }

   user, err := h.userSvc.GetUser(c, id)
   if err != nil {
       c.Error(err)
       return
   }

   c.Set("data", user)
}
```
