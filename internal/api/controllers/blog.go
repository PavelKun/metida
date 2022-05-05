package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Dsmit05/metida/internal/api/middlewares"
	"github.com/Dsmit05/metida/internal/api/response"
	"github.com/Dsmit05/metida/internal/consts"
	"github.com/Dsmit05/metida/internal/repositories"
	"github.com/gin-gonic/gin"
)

// UserContent
type SiteBlog struct {
	db *repositories.PostgresRepository
}

func NewSiteBlog(db *repositories.PostgresRepository) *SiteBlog {
	return &SiteBlog{db}
}

type CreateBlogInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Img         []byte `json:"img" binding:"required"`
}

// @Summary CreateBlog
// @Tags blog
// @Description Create Blog
// @ID protected-create-blog
// @Accept json
// @Produce json
// @Param input body CreateBlogInput true "credentials"
// @Success 200 {string} string "data"
// @Failure 400 {object} response.Error
// @Failure 403 {object} string "error"
// @Security ApiKeyAuth
// @Router /lk/blog [POST]
func (o *SiteBlog) CreateBlog(c *gin.Context) {
	// может создать только юзер с определенными правами
	// Todo: не прерывает работу функции, изменить
	middlewares.CheckAccessRights(c, consts.RoleAdmin)

	var inputData CreateBlogInput

	if err := c.ShouldBindJSON(&inputData); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeInvalidParams, "bad data, try again", nil)
		return
	}

	err := o.db.CreatBlog(inputData.Name, inputData.Description, inputData.Img)
	if err != nil {
		response.GinError(c, http.StatusUnauthorized, response.CodeBadRequest, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Поле для сайта создано"})
}

// @Summary ShowBlog
// @Tags blog
// @Description Show Blog by ID
// @ID show-blog
// @Accept json
// @Produce json
// @Param id path int true "blog_id"
// @Success 200 {object} models.Blog
// @Failure 400 {object} response.Error
// @Router /blog/{id} [GET]
func (o *SiteBlog) ShowBlog(c *gin.Context) {
	blogId := c.Param("id")
	if blogId == "" {
		err := fmt.Errorf("blog not selected")
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, "you need to select the content", err)
	}

	id, err := strconv.Atoi(blogId)
	if err != nil {
		err := fmt.Errorf("blog not selected")
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, "input content number", err)
	}

	blog, err := o.db.ReadBlog(int32(id))
	if err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeDBError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, blog)
}
