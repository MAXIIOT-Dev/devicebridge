package routers

import (
	"github.com/maxiiot/devicebridge/controllers"
	_ "github.com/maxiiot/devicebridge/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginswagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title vbase bridge API
// @version 0.1.0
// @description vbase bridge swagger.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath /api
func Route(mode string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.Default()

	// swagger
	r.GET("/swagger/*any", ginswagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))

	// cors start-----
	corsCfg := cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"DELETE", "PUT", "GET", "POST", "PATCH"},
	}
	// cors end-------
	r.Use(cors.New(corsCfg))

	r.Static("/static", "./ui/static/")
	r.StaticFile("/", "./ui/index.html")

	gpRoot := r.Group("/api", controllers.JWTAuth())
	{
		gpRoot.GET("/version", controllers.GetVersion) // app version

		gpRoot.GET("/device", controllers.ListDevice)               // 设备列表
		gpRoot.POST("/device", controllers.CreateDevice)            // 新增设备
		gpRoot.GET("/device/:dev_eui", controllers.GetDevice)       // 设备明细
		gpRoot.PUT("/device", controllers.UpdateDevice)             // 修改设备信息
		gpRoot.DELETE("/device/:dev_eui", controllers.DeleteDevice) // 删除设备

	}

	gpUser := r.Group("/api/user")
	{
		gpUser.POST("/add", controllers.JWTAuth(), controllers.CreateUser)              // 新增用户
		gpUser.POST("/login", controllers.Login)                                        // 用户登陆
		gpUser.PUT("/changepwd", controllers.JWTAuth(), controllers.ChangeUserPassword) // 修改密码
	}
	return r
}
