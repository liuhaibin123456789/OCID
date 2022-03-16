package global

/*
   图方便，op不提供申请ClientId、ClientSecret的api了
*/

var ClientId = "123456" //应用标识,服务端手动分配该id

var ResponseType = "code" //授权类型,默认使用授权码模式
var Scope = "openid"      //授权范围
//var State = ""            //应用的状态值,非必选

var RedirectUri = "localhost:8084/authorization" //回调地址

var JwtKey = "123456789"    //jwt签名加密算法的密钥
var ClientSecret = "ssssss" //应用密钥

var GrantType = "authorization_code" //授权类型

var AccessTokenMaps = make(map[int]string) //服务端access_token存储
var IdTokenMaps = make(map[int]string)     //服务端id_token存储
