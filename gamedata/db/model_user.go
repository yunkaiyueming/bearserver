package db

import (
	"fmt"
	"time"
)

const USERINFO_TABLE_NAME = "userinfo"

type UserInfo struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Pwd            string `json:"pwd"`
	CreateTime     string `json:"createtime"`
	CreateDateTime string `json:"create_datetime"`
}

type ModelUser struct {
}

func (this *ModelUser) TableName() string {
	return USERINFO_TABLE_NAME
}

func (this *ModelUser) GetUserById(id int) (userInfo UserInfo, err error) {
	sql := fmt.Sprintf("SELECT id,name,pwd,createtime,create_datetime FROM %s WHERE id='%d'", USERINFO_TABLE_NAME, id)
	err = getOrm("default").Raw(sql).QueryRow(&userInfo)
	return
}

func (this *ModelUser) GetUserByName(name string) (userInfo UserInfo, err error) {
	sql := fmt.Sprintf("SELECT id,name,pwd,createtime,create_datetime FROM %s WHERE name='%s'", USERINFO_TABLE_NAME, name)
	err = getOrm("default").Raw(sql).QueryRow(&userInfo)
	return
}

func (this *ModelUser) RegisterUser(name, pwd string) (int64, error) {
	sql := fmt.Sprintf("INSERT INTO %s(name,pwd,createtime,create_datetime) VALUES('%s', '%s', '%d', '%s')", USERINFO_TABLE_NAME, name, pwd, time.Now().Unix(), time.Now().Format("2006-01-02 15:04:05"))
	res, err := getOrm("default").Raw(sql).Exec()
	if err == nil {
		return res.RowsAffected()
	} else {
		return 0, nil
	}
}

func (this *ModelUser) CheckUserExist(name string) bool {
	var uid int
	sql := fmt.Sprintf("SELECT id FROM %s WHERE name='%s'", USERINFO_TABLE_NAME, name)
	getOrm("default").Raw(sql).QueryRow(&uid)
	return uid > 0
}
