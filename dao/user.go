package dao

import (
	"implement-oidc/model"
	"implement-oidc/tool"
)

func SelectUser(phone string) (model.UserInfo, error) {
	var user = model.UserInfo{}
	err := tool.DB.QueryRow("select * from `user` where phone=?;", phone).Scan(
		&user.ID,
		&user.Phone,
		&user.Avatar,
		&user.UserName,
		&user.Email,
		&user.WeChat,
		&user.QQ,
	)
	if err != nil {
		return model.UserInfo{}, err
	}
	return user, nil
}
