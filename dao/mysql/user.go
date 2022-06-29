package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"web_app/models"
)

const secret = "jacky"

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int64
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	//fmt.Println(count)
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser user表中插入数据
func InsertUser(user *models.User) (err error) {
	//对密码加密
	password := encryptPassword(user.Password)
	// 执行sql语句入库
	sqlStr := `insert into user (user_id, username, password) values (?,?,?)`
	_, err = db.Exec(sqlStr, user.UserId, user.Username, password)
	return err
}

func encryptPassword(oPassword string) string {
	hash := md5.New()
	hash.Write([]byte(secret))
	return hex.EncodeToString(hash.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	//fmt.Println(user)
	oPassword := user.Password
	sqlStr := `select user_id, username, password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}
	//fmt.Println(user)
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	//fmt.Println("密码正确")
	return
}

// GetUserById 根据id获取用户信息
func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := "select user_id, username from user where user_id=?"
	err = db.Get(user, sqlStr, uid)
	return

}
