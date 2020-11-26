package models

import (
	"database/sql"
	db "photoshare/database"
	"time"
)

type User struct {
	Id         int32     `gorm:"primary_key;AUTO_INCREMENT" json:"Id,string"`
	Name       string    `gorm:"not null"`
	Headimg    string    `gorm:"not null"`
	Phone      string    `gorm:"not null;type:char(11)"`
	City       int32     `gorm:"not null"`
	Brithday   time.Time `gorm:"not null;type:date"`
	Ismale     bool      `gorm:"not null"`
	Password   string    `gorm:"not null"`
	Updatetime time.Time `gorm:"not null;type:datetime"`
	Writetime  time.Time `gorm:"not null;type:datetime"`
	Token      string    `gorm:"not null;type:varchar(128)"`
}

func init() {
	if hasTable := db.GormDB.HasTable(&User{}); !hasTable {
		db.GormDB.Create(&User{})
	}
}

//定义表名
func (User) TableName() string {
	return "t_users"
}

//通过id获取用户数据
// func (u *User) GetFirst() (err error) {
// 	err = db.SqlDB.QueryRow("select id, name, headimg, city, brithday, ismale, password, updatetime, writetime from t_users where id = @id", sql.Named("id", u.Id)).Scan(&u.Id,
// 		&u.Name, &u.Headimg, &u.City, &u.Brithday, &u.Ismale, &u.Password, &u.Updatetime, &u.Writetime)
// 	return
// }

func (u *User) GetFirst() (err error) {
	db.GormDB.First(u)
	return
}

//判断是否用户手机号码是否存在
func (u *User) IsPhoneExits() (exists bool, err error) {
	var result int
	err = db.SqlDB.QueryRow("select case when exists(select 1 from t_users where phone = @phone) then 1 else 0 end", sql.Named("phone", u.Phone)).Scan(&result)
	if err != nil {
		return
	}
	if result == 1 {
		exists = true
	}
	return
}

//添加用户数据
func (u *User) Insert() (rowcount int32, err error) {
	var result sql.Result
	var temp int64

	result, err = db.SqlDB.Exec("insert into t_users(name,phone,headimg,city,brithday,ismale,password) values(@name,@phone,@headimg,@city,@brithday,@ismale,@password)",
		sql.Named("name", u.Name), sql.Named("phone", u.Phone), sql.Named("headimg", u.Headimg), sql.Named("city", u.City),
		sql.Named("brithday", u.Brithday), sql.Named("ismale", u.Ismale), sql.Named("password", u.Password))
	if err != nil {
		return
	}
	temp, err = result.RowsAffected()
	rowcount = int32(temp)
	//temp, err = result.LastInsertId()
	//u.Id = int32(temp)
	return
}
