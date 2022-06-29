package logic

import (
	"fmt"
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
)

func SignUp(p *models.ParamsSignUp) (err error) {
	// 判断用户存在不存在
	err1 := mysql.CheckUserExist(p.Username)
	if err1 != nil {
		// 数据库查询出错
		//fmt.Println(err1)
		return err1
	}

	// 生成UID
	userID := snowflake.GenID()
	fmt.Println(userID)
	// 构建user实例
	user := &models.User{
		UserId:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	// 保存进数据库

	return mysql.InsertUser(user)
}

func Login(p *models.ParamsLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是指针，就能拿到user。UserID

	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成jwt
	token, err := jwt.GenToken(user.UserId, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
