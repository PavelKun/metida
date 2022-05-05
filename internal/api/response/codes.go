package response

// Внутренние коды ответов сервера
const (
	CodeOk = iota
	CodeUnknownException
	CodeForbidden
	CodeInvalidParams
	CodeNotFound
	CodeBadRequest
	CodeInvalidJsonConversion
	CodeUserPasswordIsEmpty

	CodeCryptoError
	CodeDBError
)
