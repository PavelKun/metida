package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dsmit05/metida/internal/api/response"
	"github.com/Dsmit05/metida/internal/consts"
	"github.com/Dsmit05/metida/internal/cryptography"
	"github.com/Dsmit05/metida/internal/logger"
	"github.com/Dsmit05/metida/internal/repositories"
	"github.com/gin-gonic/gin"
)

// UserHandler defines the user controller methods
type UserHandler struct {
	db    *repositories.PostgresRepository
	token cryptography.ManagerToken
}

func NewUserHandler(db *repositories.PostgresRepository, token cryptography.ManagerToken) *UserHandler {
	return &UserHandler{db: db, token: token}
}

//Todo: add tag example
type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// @Summary Sign Up
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body CreateUserInput true "credentials"
// @Success 200 {object} response.Success
// @Failure 400 {object} response.Error
// @Router /auth/sign-up [post]
func (o *UserHandler) CreateUser(c *gin.Context) {
	// Validate input
	var inputData CreateUserInput
	logger.L.Info("Start Handler CreateUser")
	if err := c.ShouldBindJSON(&inputData); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeInvalidParams, "bad data, try again", nil)
		return
	}

	if err := o.db.CreateUser(inputData.Username, inputData.Password, inputData.Email, consts.RoleUser); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, err.Error(), err)
		return
	}

	// Здесь нужно возращать OK и предлагать подтвердить почту
	// Дальше пользователь должен подтвердить почту, после чего создаем ему сессию
	// Todo: данную реализацию делать в следующих версиях апи
	ip, agent := o.getIPandUserAgent(c)

	// Create refresh Token
	rToken, err := o.token.CreateRefreshToken()
	if err != nil {
		response.GinError(c, http.StatusInternalServerError, response.CodeCryptoError, "", nil)
		return
	}
	err = o.db.CreateSession(inputData.Email, rToken, agent, ip, time.Now().Add(consts.RefreshTokenTTL).Unix())
	if err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, err.Error(), err)
		return
	}

	// Create access token
	aToken, err := o.token.CreateToken(inputData.Email, consts.RoleUser, consts.AccessTokenTTL)
	if err != nil {
		response.GinError(c, http.StatusInternalServerError, response.CodeCryptoError, "", nil)
		return
	}

	c.JSON(http.StatusOK, response.Success{
		Code:        response.CodeOk,
		Description: "Create New User",
		Data:        gin.H{"aToken": aToken, "rToken": rToken},
	})
}

type AuthenticationUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Sign In
// @Tags auth
// @Description log in account
// @ID login-account
// @Accept json
// @Produce json
// @Param input body AuthenticationUserInput true "credentials"
// @Success 200 {object} response.Success
// @Failure 400 {object} response.Error
// @Router /auth/sign-in [post]
func (o *UserHandler) AuthenticationUser(c *gin.Context) {
	var inputData AuthenticationUserInput

	if err := c.ShouldBindJSON(&inputData); err != nil { // Todo: органиизвать функции валидации данных
		response.GinError(c, http.StatusBadRequest, response.CodeInvalidParams, "bad data, try again", nil)
		return
	}

	user, err := o.db.ReadUser(inputData.Email)
	if err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeDBError, err.Error(), err)
		return
	}

	// check password
	if inputData.Password != user.Password {
		// Todo: сделать блек лист ip, защита от брутфорса?
		err := fmt.Errorf("wrong password")
		response.GinError(c, http.StatusUnauthorized, response.CodeBadRequest, "", err)
		return
	}

	newRefreshToken, err := o.token.CreateRefreshToken()
	if err != nil {
		response.GinError(c, http.StatusInternalServerError, response.CodeCryptoError, "", nil)
		return
	}

	ip, agent := o.getIPandUserAgent(c)

	// при каждом логине создаем новую сессию
	err = o.db.CreateSession(
		inputData.Email, newRefreshToken, agent, ip, time.Now().Add(consts.RefreshTokenTTL).Unix())
	if err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, err.Error(), err)
		return
	}

	// Create access token
	aToken, err := o.token.CreateToken(inputData.Email, consts.RoleUser, consts.AccessTokenTTL)
	if err != nil {
		response.GinError(c, http.StatusInternalServerError, response.CodeCryptoError, "", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"aToken": aToken, "rToken": newRefreshToken}})

}

type RefreshTokenInput struct {
	RefreshToken string `json:"rtoken" binding:"required"`
}

// @Summary Refresh token
// @Tags auth
// @Description refresh access token
// @ID refresh-token
// @Accept json
// @Produce json
// @Param input body RefreshTokenInput true "credentials"
// @Success 200 {object} response.Success
// @Failure 400 {object} response.Error
// @Router /auth/refresh [post]
func (o *UserHandler) RefreshTokenUser(c *gin.Context) {
	var inputData RefreshTokenInput

	if err := c.ShouldBindJSON(&inputData); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeInvalidParams, "bad data, try again", nil)
		return
	}

	userData, err := o.db.ReadEmailRoleWithRefreshToken(inputData.RefreshToken)
	if err != nil {
		response.GinError(c, http.StatusUnauthorized, response.CodeBadRequest, err.Error(), err)
		return
	}

	// Todo: Здесь получаем время жизни токена, и проверяем не прошло ли оно

	// Create access token
	aToken, err := o.token.CreateToken(userData.Email, userData.Role, consts.AccessTokenTTL)
	if err != nil {
		response.GinError(c, http.StatusInternalServerError, response.CodeCryptoError, "", nil)
		return
	}

	rToken, err := o.token.CreateRefreshToken()
	if err != nil {
		response.GinError(c, http.StatusInternalServerError, response.CodeCryptoError, "", nil)
		return
	}

	err = o.db.UpdateSessionTokenOnly(
		inputData.RefreshToken, rToken, time.Now().Add(consts.RefreshTokenTTL).Unix())
	if err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, "", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"aToken": aToken, "rToken": rToken}})
}

func (o *UserHandler) AuthenticationRole(c *gin.Context) {

}

func (o *UserHandler) getIPandUserAgent(c *gin.Context) (IP, UserAgent string) {
	val, ok := c.Request.Header["User-Agent"]

	if !ok {
		logger.L.Error("not have User-Agent Header")
	} else {
		// Todo: переделать
		// val содержит: [Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36]
		UserAgent = val[0]
	}

	IP = c.ClientIP()

	return
}
