package internal

import (
	"time"

	_ "bearserver/conf"
	"bearserver/gamedata/db"
	"bearserver/msg"
)

type UserData struct { //数据库的数据
	Id         int    "_id" //用户id 自增型的
	Uid        int    //用户线上看到的id
	Name       string //用户的昵称
	Pwd        string //用户密码
	CreateTime int64  //注册时间
}

func (data *UserData) initValue() error {
	data.CreateTime = time.Now().Unix()
	return nil
}

//注册
func register(userInfo *msg.RegisterUserInfo) (int64, error) {
	userModel := db.ModelUser{}
	uid, err := userModel.RegisterUser(userInfo.Name, userInfo.Pwd)
	return uid, err
}

//登录
func login(user *msg.UserLoginInfo) (int, error) {
	userModel := db.ModelUser{}
	uInfo, err := userModel.GetUserByName(user.Name)
	return uInfo.Id, err
}

//检查用户是否已注册过
func checkExitedUser(name string) bool {
	userModel := db.ModelUser{}
	return userModel.CheckUserExist(name)
}
