package tool

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"implement-oidc/global"
	"strings"
)

/*
   统一采用某种固定的编码方式和加密算法
*/

// Header token的header必要字段
type Header struct {
	TokenType string `json:"typ"` //令牌的类型，如JWT
	Algorithm string `json:"alg"` //正在使用的散列算法，例如HMAC、SHA256或RSA
}

//Payload 必要的jwt payload字段,添加id token、access token所需字段
//id token: 不存储用户敏感信息，增加存储用户名及对应手机号。此token用于认证用户信息
//access token:不存用户相关信息，此token用于授权
type Payload struct {
	Audience  string `json:"aud,omitempty"` //接受者的url地址
	ExpiresAt int64  `json:"exp,omitempty"` //该jwt销毁的时间；unix时间戳
	Id        string `json:"jti,omitempty"` //该jwt的唯一ID编号
	IssuedAt  int64  `json:"iat,omitempty"` //该jwt的发布时间；unix时间戳
	Issuer    string `json:"iss,omitempty"` //发布者的url地址
	NotBefore int64  `json:"nbf,omitempty"` //该jwt的使用时间不能早于该时间；unix时间戳
	Subject   string `json:"sub,omitempty"` //该JWT所面向的用户，用于处理特定应用，不是常用的字段

	UserName string `json:"user_name,omitempty"` //用户名，用以身份认证
	Phone    string `json:"phone,omitempty"`     //用户手机号，用以身份认证
}

func CreateToken(header Header, payload Payload) (string, error) {
	//json化数据
	headJson, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	fmt.Println("headJson:", string(headJson))
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	fmt.Println("payloadJson:", string(payloadJson))
	headEncode, err := encode(headJson)
	if err != nil {
		return "", err
	}
	fmt.Println("headEncode:", headEncode)
	payloadEncode, err := encode(payloadJson)
	if err != nil {
		return "", err
	}
	fmt.Println("payloadEncode:", payloadEncode)
	token, err := secretSignature(headEncode, payloadEncode, global.JwtKey)
	if err != nil {
		return "", err
	}
	fmt.Println("token:", token)
	return token, nil
}

//编码jwt`s header\payload
func encode(jsonStr []byte) (string, error) {
	//校验json数据
	res := json.Valid(jsonStr)
	if !res {
		return "", errors.New("the header json string is wrong")
	}
	//编码header，采用BASE64压缩编码,便于http传输url,除去填充字符串
	encodeStr := base64.RawURLEncoding.EncodeToString(jsonStr)
	return encodeStr, nil
}

//json jwt header\payload
func decode(encodeStr string) (string, error) {
	//解码
	jsonStr, err := base64.RawURLEncoding.DecodeString(encodeStr)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}

//加密jwt`s signature
func secretSignature(headerEncode string, payloadEncode, secretKey string) (string, error) {
	//合并header+payload
	sigStr := headerEncode + "." + payloadEncode
	//使用密钥进行签名
	hash := hmac.New(sha256.New, []byte(secretKey))
	_, err := hash.Write([]byte(sigStr))
	if err != nil {
		return "", err
	}
	//使用base64编码
	return sigStr + "." + base64.RawURLEncoding.EncodeToString(hash.Sum(nil)), nil
}

func ParseAccessToken(tokenString string) (Header, Payload, error) {

	//校验jwt格式
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return Header{}, Payload{}, errors.New("token contains an invalid number of segments")
	}

	//验证签名
	sigStr := parts[0] + "." + parts[1]
	hash := hmac.New(sha256.New, []byte(global.JwtKey))
	_, err := hash.Write([]byte(sigStr))
	if err != nil {
		return Header{}, Payload{}, err
	}
	accessToken := global.AccessTokenMaps[1]
	rightParts := strings.Split(accessToken, ".")
	if rightParts[2] != base64.RawURLEncoding.EncodeToString(hash.Sum(nil)) {
		return Header{}, Payload{}, errors.New("the signature is wrong")
	}

	//解码并json化header
	h, err := decode(parts[0])
	if err != nil {
		return Header{}, Payload{}, err
	}
	var head = Header{}
	err = json.Unmarshal([]byte(h), &head)
	if err != nil {
		return Header{}, Payload{}, err
	}

	//解析payload
	p, err := decode(parts[1])
	if err != nil {
		return Header{}, Payload{}, err
	}
	var payload = Payload{}
	err = json.Unmarshal([]byte(p), &payload)
	if err != nil {
		return Header{}, Payload{}, err
	}

	return head, payload, nil
}

// ParseIdToken 主要用于解析获取身份认证信息，及payload中的包含用户信息字段：phone\userName
func ParseIdToken(tokenString string) (Payload, error) {

	//校验jwt格式
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return Payload{}, errors.New("token contains an invalid number of segments")
	}

	//验证签名
	sigStr := parts[0] + "." + parts[1]
	hash := hmac.New(sha256.New, []byte(global.JwtKey))
	_, err := hash.Write([]byte(sigStr))
	if err != nil {
		return Payload{}, err
	}
	idToken := global.IdTokenMaps[1]
	rightParts := strings.Split(idToken, ".")
	if rightParts[2] != base64.RawURLEncoding.EncodeToString(hash.Sum(nil)) {
		return Payload{}, errors.New("the signature is wrong")
	}

	//解析payload
	p, err := decode(parts[1])
	if err != nil {
		return Payload{}, err
	}
	var payload = Payload{}
	err = json.Unmarshal([]byte(p), &payload)
	if err != nil {
		return Payload{}, err
	}

	return payload, nil
}
