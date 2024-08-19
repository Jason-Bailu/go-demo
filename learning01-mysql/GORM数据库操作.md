# GORM数据库操作

- 安装GORM command：**go get -u gorm.io/gorm**
- 安装MySQL驱动 command：**go get -u gorm.io/driver/mybatis**
- 文档：**[GORM](https://gorm.io/zh_CN/docs/)开发文档**

## 连接操作

```go
package main

// 导入 gorm mysql包
import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model // gorm.Model 对应 id、created_time、updated_time、deleted_Time
	Username string // username
	Password string // password
	Gender   uint // gender
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

	// 创建记录 模型有 CreatedAt字段，该字段的值将会是初次创建记录的时间
	u1 := User{Username: "bailu", Password: "123456", Gender: 1}
	db.Create(&u1)

	// 搜索记录
	var u = new(User)
	db.First(u)
	fmt.Printf("%#v\n", u)

	var uu User
	db.Find(&uu, "username=?", "bailu")
	fmt.Printf("%#v\n", uu)

	// 更新记录 模型有UpdatedAt字段，该字段的值将会是每次更新记录的时间
	db.Model(&u).Update("gender", 0)

	// 删除记录 将会设置DeletedAt字段为当前时间
	db.Delete(&u)
}
```

## 字段操作

- 模型定义示例：**通过tag来修饰字段** 修改对应字段gorm:"column:xxx"

  ```go
  type User struct {
    gorm.Model
    Name         string  `gorm:"column:username"`
    Age          sql.NullInt64
    Birthday     *time.Time
    Email        string  `gorm:"type:varchar(100);unique_index"`
    Role         string  `gorm:"size:255"` // 设置字段大小为255
    MemberNumber *string `gorm:"unique;not null"` // 设置会员号（member number）唯一并且不为空
    Num          int     `gorm:"AUTO_INCREMENT"` // 设置 num 为自增类型
    Address      string  `gorm:"index:addr"` // 给address字段创建名为addr的索引
    IgnoreMe     int     `gorm:"-"` // 忽略本字段
  }
  ```

- **默认ID为主键**

- **自定义表名：**

  ```go
  type User struct {} // 默认表名是 `users`
  
  // 将 User 的表名设置为 `profiles`
  func (User) TableName() string {
    return "profiles"
  }
  
  func (u User) TableName() string {
    if u.Role == "admin" {
      return "admin_users"
    } else {
      return "users"
    }
  }
  ```

## CRUD操作

```go
package main

import (
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

	// 创建对象
	user := User{Username: "bailu", Password: "123456", Gender: 1}
	result := db.Create(&user)
	fmt.Println(result.Error)        // nil
	fmt.Println(result.RowsAffected) // 1

	// 批量添加
	users := []User{{Username: "aaa", Password: "123456", Gender: 1},
		{Username: "bbb", Password: "123456", Gender: 0},
		{Username: "ccc", Password: "123456", Gender: 1}}
	result := db.CreateInBatches(users, 10)
	fmt.Println(result.Error)
	fmt.Println(result.RowsAffected)

	// 一般查询
	var user User
	// 第一条数据
	db.First(&user)
	fmt.Printf("%#v\n", user)
	// 随机一条数据
	db.Take(&user)
	fmt.Printf("%#v\n", user)
	// 最后一条数据
	db.Last(&user)
	fmt.Printf("%#v\n", user)
	// 所有记录
	var users []User
	db.Find(&users)
	fmt.Printf("%#v\n", users)
	// 查询指定id
	var user User
	db.First(&user, 4)
	fmt.Printf("%#v\n", user)

	// 查询 Where Not Or Select Omit
	// =
	var user User
	db.Where("username = ?", "bbb").First(&user)
	fmt.Printf("%#v\n", user)
	var users []User
	db.Where("gender = ?", 1).Find(&users)
	fmt.Printf("%#v\n", users)
	// !=
	var users []User
	db.Where("username != ?", "bbb").Find(&users)
	fmt.Printf("%#v\n", users)
	// IN
	var users []User
	db.Where("username IN (?)", []string{"aaa", "bbb"}).Find(&users)
	fmt.Printf("%#v\n", users)
	// LIKE
	var users []User
	db.Where("username LIKE ?", "%b%").Find(&users)
	fmt.Printf("%#v\n", users)
	// AND
	var users []User
	db.Where("username = ? AND gender = ?", "bbb", "0").Find(&users)
	fmt.Printf("%#v\n", users)
	// Time
	var users []User
	db.Where("created_at > ?", time.DateTime).Find(&users)
	fmt.Printf("%#v\n", users)

	// 更新
	// 更新所有字段
	var user User
	db.First(&user)
	user.Username = "白鹭"
	user.Gender = 1
	db.Save(user)
	// 更新指定字段
	var user User
	db.First(&user)
	db.Model(&user).Update("username", "bailu")
	// 条件更新字段
	var user User
	db.Model(&user).Where("username = ?", "bbb").Update("gender", 1)
	// map 更新多个属性，只会更新其中有变化的属性
	var user User
	db.First(&user)
	db.Model(&user).Updates(map[string]interface{}{"username": "bailu", "password": "654321"})

	// 删除
	// 单个删除
	var user User
	db.First(&user)
	db.Delete(&user) // 保证id有值，才能删除
	// 批量删除
	db.Where("gender = ?", 1).Delete(&User{})
	db.Delete(&User{}, "gender = ?", 1)
	// 存在deletedAt字段都为软删除
	// 物理删除 Unscoped也可以查询软删除的字段
	var user User
	db.First(&user)
	db.Unscoped().Delete(&user)
}
```

