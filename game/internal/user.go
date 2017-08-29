package internal

import (
	"bearserver/msg"
	"bearserver/conf"
	//"gopkg.in/mgo.v2/bson"
	"fmt"
	//"github.com/name5566/leaf/log"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type UserData struct {//数据库的数据
	Id	int	"_id"	//用户id 自增型的
	Uid	int	//用户线上看到的id
	Name	string	//用户的昵称
	Pwd string //用户密码
	CreateTime	int64	//注册时间
}

func (data *UserData) initValue() error {
	data.CreateTime = time.Now().Unix()
	return nil
}

func register(userInfo *msg.RegisterUserInfo)  (err error) {//注册
	//var user User
	//userInfo := args[0].(*msg.RegisterUserInfo)
	//log.Debug("11111")
	skeleton.Go(func() {
		db, err := sql.Open("mysql", conf.Server.DBUrl)
		checkErr(err)
		defer db.Close()
		//方式4 insert
		//Begin函数内部会去获取连接
		tx,_ := db.Begin()
		//每次循环用的都是tx内部的连接，没有新建连接，效率高
		tx.Exec("INSERT INTO userinfo(name,pwd,createtime) values(?,?,?)",userInfo.Name,userInfo.Pwd,time.Now().Unix())
		//最后释放tx内部的连接
		tx.Commit()
	}, func() {

	})
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}


func login(user  *msg.UserLoginInfo)(err error) {
	fmt.Println("---userinfo---",user)
	skeleton.Go(func() {
		db, err := sql.Open("mysql", conf.Server.DBUrl)
		stmtOut, err := db.Query("SELECT * FROM userinfo WHERE name = ?",user.Name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtOut.Close()
	}, func() {

	})
	return
}

//检查用户是否已注册过
func checkExitedUser(userName string) (err error){
	db, err := sql.Open("mysql", conf.Server.DBUrl)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var name string
	err = db.QueryRow("SELECT name FROM userinfo WHERE name=?", userName).Scan(&name)
	fmt.Println(err)
	if err == sql.ErrNoRows{
		return err
	}

	return nil
}