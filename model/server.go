package model

type ReqCode struct {
	ClientId     string `json:"client_id"  example:"122334" form:"client_id"`                      //应用标识
	RedirectUri  string `json:"redirect_uri" example:"http://www.example.com" form:"redirect_uri"` //回调地址
	ResponseType string `json:"response_type" example:"code" form:"response_type"`                 //授权类型，固定值：code。
	Scope        string `json:"scope" example:"openid" form:"scope"`                               //授权范围，固定值：openid。
	State        string `json:"state" example:"test string" form:"state"`                          //应用的状态值。可用于防止CSRF攻击，成功授权后回调应用时会原样带回，应用用它校验认证请求与回调请求的对应关系。可以包含字母和数字。
}

type ReqToken struct {
	ClientId     string `json:"client_id"   form:"client_id" example:"122334"`                     //应用标识
	ClientSecret string `json:"client_secret" form:"client_secret" example:"11111"`                //应用密钥
	Code         string `json:"code" form:"code" example:"123414"`                                 //授权码，认证登录后回调获取的授权码
	GrantType    string `json:"grant_type" form:"grant_type" example:"authorization_code"`         //授权类型，固定值：authorization_code
	RedirectUri  string `json:"redirect_uri" form:"redirect_uri" example:"http://www.example.com"` //回调地址
}

type RepToken struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
	TokenType   string `json:"token_type" example:"jwt"`
	ExpiresIn   string `json:"expires_in" example:"457388"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Phone    string `json:"phone" `     //当前登录用户的手机号
	Avatar   string `json:"avatar" `    //用户头像绝对路径
	UserName string `json:"user_name" ` //用户的昵称,依据隐私表数据
	Email    string `json:"email"`      //用户的邮箱
	WeChat   string `json:"we_chat"`    //用户微信
	QQ       string `json:"qq"`         //用户qq
}
