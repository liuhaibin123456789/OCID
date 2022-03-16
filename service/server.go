package service

import (
	"errors"
	"implement-oidc/dao"
	"net/http"
	"strconv"
	"time"

	"implement-oidc/global"
	"implement-oidc/model"
	"implement-oidc/tool"

	"github.com/gin-gonic/gin"
)

func GetCode(reqCode *model.ReqCode, c *gin.Context) model.ErrMod {
	//校验数据
	if reqCode.Scope != global.Scope {
		return model.ErrMod{
			Status: http.StatusBadRequest,
			Err:    errors.New("the scope is not 'openid'"),
		}
	}
	if reqCode.ResponseType != global.ResponseType {
		return model.ErrMod{
			Status: http.StatusBadRequest,
			Err:    errors.New("the response_type is not 'code'"),
		}
	}
	if reqCode.ClientId != global.ClientId {
		return model.ErrMod{
			Status: http.StatusBadRequest,
			Err:    errors.New("the client_id is not wrong"),
		}
	}
	//正则匹配todo
	if reqCode.RedirectUri == "" {
		return model.ErrMod{
			Status: http.StatusBadRequest,
			Err:    errors.New("the redirect_uri is empty"),
		}
	}
	//生成随机code
	randNum := tool.RandNum()

	//todo写入redis，十分钟有效。

	//重定向并附上code
	c.Redirect(http.StatusMovedPermanently, reqCode.RedirectUri+"?code="+strconv.FormatInt(int64(randNum), 10)+"&state="+reqCode.State)

	return model.ErrMod{
		Status: 200,
		Err:    nil,
	}
}

func CreateToken(reqToken model.ReqToken) (model.RepToken, model.ErrMod) {
	//校验入参
	//todo校验redis里十分钟有效的code
	//if reqToken.Code {
	//
	//}
	if reqToken.ClientId != global.ClientId {
		return model.RepToken{}, model.ErrMod{
			Err:    errors.New("the client_id is wrong"),
			Status: http.StatusBadRequest,
		}
	}
	if reqToken.ClientSecret != global.ClientSecret {
		return model.RepToken{}, model.ErrMod{
			Err:    errors.New("the client_secret is wrong"),
			Status: http.StatusBadRequest,
		}
	}
	if reqToken.GrantType != global.GrantType {
		return model.RepToken{}, model.ErrMod{
			Err:    errors.New("the grant_type is wrong"),
			Status: http.StatusBadRequest,
		}
	}
	//todo正则匹配域名或ip+端口
	if reqToken.RedirectUri == "" {
		return model.RepToken{}, model.ErrMod{
			Err:    errors.New("the redirect_uri is empty"),
			Status: http.StatusBadRequest,
		}
	}

	//token编号
	var num1, num2 int
	num1++
	num2++
	//生成access_token
	header := tool.Header{Algorithm: "HS256", TokenType: "JWT"}
	payload1 := tool.Payload{
		Id:        strconv.Itoa(num1),
		ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		NotBefore: time.Now().Unix(),
		Issuer:    "cold bin",
	}
	accessToken, err := tool.CreateToken(header, payload1)
	if err != nil {
		return model.RepToken{}, model.ErrMod{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}
	//保存access_token
	global.AccessTokenMaps[num1] = accessToken

	//生成id_token
	payload2 := tool.Payload{
		Id:        strconv.Itoa(num2),
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		NotBefore: time.Now().Unix(),
		Issuer:    "cold bin",
		//根据应用id查询本服务端数据库部分用户非隐私信息，此处省略图方便，就模拟一下。
		Phone:    "15736469310", //用户手机号
		UserName: "阿冰",          //用户名
	}
	idToken, err := tool.CreateToken(header, payload2)
	if err != nil {
		return model.RepToken{}, model.ErrMod{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}
	//保存id_token
	global.IdTokenMaps[num2] = idToken

	return model.RepToken{
			AccessToken: accessToken,
			IdToken:     idToken,
			ExpiresIn:   time.Now().Add(time.Hour * 12).String(),
			TokenType:   "Bearer",
		}, model.ErrMod{
			Err:    nil,
			Status: http.StatusOK,
		}
}

func GetUser(accessToken, idToken string) (model.UserInfo, model.ErrMod) {
	//身份认证
	p, err := tool.ParseIdToken(idToken)
	if err != nil {
		return model.UserInfo{}, model.ErrMod{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	//通过access_token授权拿取资源
	_, _, err = tool.ParseAccessToken(accessToken)
	if err != nil {
		return model.UserInfo{}, model.ErrMod{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	user, err := dao.SelectUser(p.Phone)
	if err != nil {
		return model.UserInfo{}, model.ErrMod{
			Err:    err,
			Status: http.StatusNotFound,
		}
	}

	return user, model.ErrMod{
		Err:    nil,
		Status: http.StatusOK,
	}
}
