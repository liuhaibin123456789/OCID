package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"implement-oidc/tool"

	"implement-oidc/api"
	_ "implement-oidc/docs" //init函数，调用该包，手动导入
)

// @title           OIDC协议
// @version         1.0
// @host            localhost:8081
// @description     基于golang开发得OIDC协议op，用于提供 code\id token\access token\第三方用户信息
// @securityDefinitions.apikey CoreAPI
// @name Authorization
// @in header
func main() {
	err := tool.InitMysql()
	if err != nil {
		return
	}
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) //注册swagger路由
	router.GET("/api/v1/oauth2/authorize", api.GetCode)                       //模拟op：重定向某url至code、state参数
	router.POST("/api/v1/oauth2/token", api.CreateToken)                      //模拟op,生成access_token、id_token
	router.GET("/api/v1/user", api.GetUser)                                   //校验两个token，获取本服务端用户资源
	err = router.Run(":8081")
	if err != nil {
		return
	}
}
