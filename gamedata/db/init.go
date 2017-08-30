package db

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	registerDB("default")
}

func registerDB(name string) error {
	orm.RegisterDataBase("default", "mysql", "root:@tcp(192.168.8.81:3307)/bear")
	orm.SetMaxOpenConns(name, 5000)
	orm.SetMaxIdleConns(name, 4000)
	return nil
}

func getOrm(name string) orm.Ormer {
	o := orm.NewOrm()
	o.Using(name)
	return o
}
