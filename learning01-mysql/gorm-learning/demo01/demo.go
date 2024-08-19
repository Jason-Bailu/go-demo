package main

// 导入 gorm mysql包
import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type User struct {
	gorm.Model
	Username string
	Password string
	Gender   uint
}

func main() {
	dsn := "root:dxy000216@tcp(127.0.0.1:3306)/godata?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		panic(err)
	}

	// 自动迁移
	db.AutoMigrate(&User{})

	// 创建记录
	u1 := User{Username: "bailu", Password: "123456", Gender: 1}
	db.Create(&u1)

	// 搜索记录
	var u = new(User)
	db.First(u)
	fmt.Printf("%#v\n", u)

	var uu User
	db.Find(&uu, "username=?", "bailu")
	fmt.Printf("%#v\n", uu)

	// 更新记录
	db.Model(&u).Update("gender", 0)

	// 删除记录
	db.Delete(&u)
}
