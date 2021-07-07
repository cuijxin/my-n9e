package models

import (
	"encoding/json"
	"errors"

	"github.com/toolkits/pkg/logger"
	"gorm.io/gorm"
)

type User struct {
	Id       int64           `json:"id"`
	Username string          `json:"username"`
	Nickname string          `json:"nickname"`
	Password string          `json:"-"`
	Phone    string          `json:"phone"`
	Email    string          `json:"email"`
	Portrait string          `json:"portrait"`
	Status   int             `json:"status"`
	Role     string          `json:"role"`
	Contacts json.RawMessage `json:"contacts"` // 内容为 map[string]string 结构
	CreateAt int64           `json:"create_at"`
	CreateBy string          `json:"ceate_by"`
	UpdateAt int64           `json:"update_at"`
	UpdateBy string          `json:"update_by"`
}

func PassLogin(username, pass string) (*User, error) {
	user, err := UserGetByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		logger.Infof("password auth fail, no such user: %s", username)
		return nil, loginFailError
	}

	loginPass, err := CryptoPass(pass)
	if err != nil {
		return nil, internalServerError
	}

	if loginPass != user.Password {
		logger.Infof("password auth fail, password error, user: %s", username)
		return nil, loginFailError
	}

	return user, nil
}

func UserGetByUsername(username string) (*User, error) {
	return UserGet("username=?", username)
}

func UserGetById(id int64) (*User, error) {
	return UserGet("id=?", id)
}

func UserGet(where string, args ...interface{}) (*User, error) {
	var obj User
	err := DB.Where(where, args).First(&obj).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		logger.Errorf("mysql.error: query user(%s)%+v fail: %s", where, args, err)
		return nil, internalServerError
	}
	return &obj, nil
}
