// pkg/errcode/code.go
package errcode

const (
	CodeSuccess      = 200
	CodeSystemError  = 100000
	CodeParamError   = 100002
	CodeDecryptError = 100003
	CodeNotFound     = 100404
	CodeDBError      = 101001
	CodeCacheError   = 102001

	CodeTokenGenerateError = 103001
	CodeTokenSaveError     = 103002
	CodeTokenDeleteError   = 103003

	CodeFilePathGenError = 104001
	CodeFileSaveError    = 104002
)

const (
	CodeUnauthorized    = 200001
	CodeInvalidAuth     = 200002
	CodeUserLocked      = 200003
	CodeUserDisabled    = 200005
	CodeTooManyAttempts = 200006

	CodeNoPermission = 201001
	CodeOpDenied     = 201002
)

const (
	CodeFileEmpty        = 202001
	CodeFileTypeMismatch = 202002
	CodeFileTooLarge     = 202003
)
