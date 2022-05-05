package middlewares

import (
	"fmt"
	"net/http"

	"github.com/Dsmit05/metida/internal/cryptography"

	"github.com/gin-gonic/gin"
)

// Todo: Привести все ошибки к типовым
type ProtectedMidleware struct {
	auth cryptography.ManagerToken
}

func NewProtectedMidleware(token cryptography.ManagerToken) *ProtectedMidleware {
	return &ProtectedMidleware{auth: token}
}

func (o *ProtectedMidleware) AuthMidleware(c *gin.Context) {
	email, role, err := o.parseAuthHeader(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusExpectationFailed, gin.H{"error": err.Error()}) // 417
		return
	}

	c.Set("email", email)
	c.Set("role", role)
}

func (o *ProtectedMidleware) parseAuthHeader(c *gin.Context) (email, role string, err error) {
	header := c.GetHeader("Authorizations")
	if header == "" {
		// ToDo: привести все к стандарту и возвращать нормальный ответ в доп описании: Пожалуйста авторизуйтесь
		return "", "", fmt.Errorf("empty header")
	}

	// // ToDo: ограничить количество символов headerParts
	// headerParts := strings.Split(header, " ")

	return o.auth.ParseToken(header)
}

// ToDo: вынести в отдельый файл или пакет? здесь явно возвращать ошибку
// CheckAccessRights Проверяет у хендлера доступы, вызывает аборт
func CheckAccessRights(c *gin.Context, roles ...string) {
	roleFromContext, ok := c.Get("role")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user without role"})
	}

	email, ok := c.Get("email")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user without email"})
	}

	for _, role := range roles {
		if roleFromContext == role {
			return
		}
	}

	msg := fmt.Sprintf("User %v has no enough rights for access to resource.", email)
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": msg})

}
