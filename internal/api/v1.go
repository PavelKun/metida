package api

import (
	"time"

	_ "github.com/Dsmit05/metida/docs"
	"github.com/Dsmit05/metida/internal/api/controllers"
	"github.com/Dsmit05/metida/internal/api/middlewares"
	"github.com/Dsmit05/metida/internal/config"
	"github.com/Dsmit05/metida/internal/cryptography"
	"github.com/Dsmit05/metida/internal/logger"
	"github.com/Dsmit05/metida/internal/repositories"
	ginzap "github.com/gin-contrib/zap"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// собираем все контроллеры тут

type ApiV1 struct {
	*controllers.UserHandler
	*controllers.UserContent
	*controllers.SiteBlog
	*controllers.ProfilingHandler
	*middlewares.ProtectedMidleware
	cfg *config.Config
}

func V1(
	db *repositories.PostgresRepository,
	ManagerToken cryptography.ManagerToken,
	cfg *config.Config,
) *ApiV1 {

	userHandler := controllers.NewUserHandler(db, ManagerToken)
	wallEditorialsHandler := controllers.NewWallEditorials(db)
	siteBlog := controllers.NewSiteBlog(db)
	profilingHandler := controllers.NewProfilingHandler()
	protectedMidleware := middlewares.NewProtectedMidleware(ManagerToken)

	return &ApiV1{
		userHandler,
		wallEditorialsHandler,
		siteBlog,
		profilingHandler,
		protectedMidleware,
		cfg,
	}
}

func (o *ApiV1) Start() error {
	var r *gin.Engine

	// включаем режимы
	if o.cfg.IfDebagOn() {
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	// подключаем логирование всех запросов
	if o.cfg.IfLogRequestsOn() {
		r.Use(ginzap.Ginzap(logger.Prod, time.RFC3339, true))
	}

	// подключаем профилирование
	if o.cfg.IfPprofOn() {
		o.ReqisterRouteProfiling(r.Group("/admin"))
	}

	control := r.Group("/auth")
	{
		control.POST("/sign-up", o.CreateUser)
		control.POST("/sign-in", o.AuthenticationUser)
		control.POST("/refresh", o.RefreshTokenUser)
	}

	lk := r.Group("/lk")
	lk.Use(o.AuthMidleware)
	{
		lk.GET("/content/:id", o.ShowContent)
		lk.POST("/content", o.CreateContent)
		lk.POST("/blog", o.CreateBlog)
	}
	r.GET("/blog/:id", o.ShowBlog)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r.Run(":8080")
}
