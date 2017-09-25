package internal

import (
	"bearserver/gamedata/db"
	"bearserver/msg"
)

type UserModule struct{}

//注册
func (u *UserModule) Register(userInfo *msg.RegisterUserInfo) (int64, error) {
	userModel := db.ModelUser{}
	uid, err := userModel.RegisterUser(userInfo.Name, userInfo.Pwd)
	return uid, err
}

//登录
func (u *UserModule) Login(user *msg.UserLoginInfo) (int, error) {
	userModel := db.ModelUser{}
	uInfo, err := userModel.GetUserByName(user.Name)
	return uInfo.Id, err
}

//检查用户是否已注册过
func (u *UserModule) CheckExitedUser(name string) bool {
	userModel := db.ModelUser{}
	return userModel.CheckUserExist(name)
}
