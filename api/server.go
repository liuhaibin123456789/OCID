package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"implement-oidc/global"
	"implement-oidc/model"
	"implement-oidc/service"
	"implement-oidc/tool"
)

// GetCode
// @summary  生成服务端授权码，redirect_uri放在最后，防止其后面得参数失效
// @produce  json
// @Param    client_id     query  string  true  "本平台申请得应用ID，固定用于测试"
// @param    response_type  query  string  true  "必为'code'"
// @param    scope  query  string  true  "必为'openid'"
// @param    state  query  string  false  "可选"
// @param    redirect_uri  query  string  true  "授权成功后，跳转的地址"
// @router   /api/v1/oauth2/authorize [get]
func GetCode(c *gin.Context) {
	//获取数据
	clientId := c.Query("client_id")
	redirectUri := c.DefaultQuery("redirect_uri", global.RedirectUri)
	responseType := c.Query("response_type")
	scope := c.Query("scope")
	state := c.Query("state")
	//组装
	reqCode := &model.ReqCode{
		ClientId:     clientId,
		RedirectUri:  redirectUri,
		ResponseType: responseType,
		Scope:        scope,
		State:        state,
	}
	errMod := service.GetCode(reqCode, c)
	if errMod.Err != nil {
		tool.FormatJson(errMod.Status, errMod.Err, "", c)
	}
}

func CreateToken(c *gin.Context) {
	reqToken := model.ReqToken{}
	err := c.ShouldBind(&reqToken)
	if err != nil {
		fmt.Println("ShouldBind reqToken failed.")
		return
	}
	repToken, errMod := service.CreateToken(reqToken)
	if errMod.Err != nil {
		tool.FormatJson(errMod.Status, errMod.Err, "", c)
	}
	tool.FormatJson(errMod.Status, nil, repToken, c)
}

func GetUser(c *gin.Context) {
	accessToken := c.PostForm("access_token")
	idToken := c.PostForm("id_token")

	user, errMod := service.GetUser(accessToken, idToken)
	if errMod.Err != nil {
		tool.FormatJson(errMod.Status, errMod.Err, "", c)
	}
	tool.FormatJson(errMod.Status, nil, user, c)

}
