package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/str"
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

func (u *User) TableName() string {
	return "users"
}

func (u *User) Validate() error {
	u.Username = strings.TrimSpace(u.Username)

	if u.Username == "" {
		return _e("Username is blank")
	}

	if str.Dangerous(u.Username) {
		return _e("Username has invalid characters")
	}

	if str.Dangerous(u.Nickname) {
		return _e("Nickname has invalid characters")
	}

	if u.Phone != "" && !str.IsPhone(u.Phone) {
		return _e("Phone invalid")
	}

	if u.Email != "" && !str.IsMail(u.Email) {
		return _e("Email invalid")
	}

	return nil
}

func (u *User) Update(cols ...string) error {
	if err := u.Validate(); err != nil {
		return err
	}

	err := DB.Model(u).Select(cols).Updates(*u).Error
	if err != nil {
		logger.Errorf("mysql.error: update user fail: %v", err)
		return internalServerError
	}

	return nil
}

func (u *User) Add() error {
	result := DB.Where("username=?", u.Username).Find(u)
	if result.Error != nil {
		logger.Errorf("mysql.error: count user(%s) fail: %v", u.Username, result.Error)
		return internalServerError
	}

	if result.RowsAffected > 0 {
		return _e("Username %s already exists", u.Username)
	}

	return DBInsertOne(u)
}

func InitRoot() {
	var u User
	err := DB.Where("username=?", "root").First(&u).Error
	if err == nil {
		return
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("fatal: cannot query user root", err)
		os.Exit(1)
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		pass, err := CryptoPass("root.2020")
		if err != nil {
			fmt.Println("fatal: cannot crypto password,", err)
			os.Exit(1)
		}

		now := time.Now().Unix()

		u = User{
			Username: "root",
			Password: pass,
			Nickname: "超管",
			Portrait: "",
			Role:     "Admin",
			CreateAt: now,
			UpdateAt: now,
			CreateBy: "system",
			UpdateBy: "system",
		}

		err = DB.Create(&u).Error
		if err != nil {
			fmt.Println("fatal: cannot insert user root", err)
			os.Exit(1)
		}
	}

	fmt.Println("user root init done")
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

func UserTotal(query string) (num int64, err error) {
	if query != "" {
		q := "%" + query + "%"
		result := DB.Model(&User{}).Where("username like ?", q).Or("nickname like ?", q).Or("phone like ?", q).Or("email like ?", q).Count(&num)
		num = result.RowsAffected
		err = result.Error
	} else {
		var objs []User
		result := DB.Find(&objs)
		num = result.RowsAffected
		err = result.Error
	}

	if err != nil {
		logger.Errorf("mysql.error: count user(query: %s) fail: %v", query, err)
		return num, internalServerError
	}

	return num, nil
}

func UserGets(query string, limit, page int) ([]User, error) {
	var users []User
	if query != "" {
		q := "%" + query + "%"
		err := DB.Where("username like ?", q).
			Or("nickname like ?", q).
			Or("phone like ?", q).
			Or("email like ?", q).
			Scopes(Paginate(page, limit)).
			Find(&users).Error
		if err != nil {
			logger.Errorf("mysql.error: select user(query: %s) fail: %v", query, err)
			return users, internalServerError
		}
	} else {
		err := DB.Scopes(Paginate(page, limit)).Find(&users).Error
		if err != nil {
			logger.Errorf("mysql.error: select user fail: %v", err)
			return users, internalServerError
		}
	}

	if len(users) == 0 {
		return []User{}, nil
	}

	return users, nil
}

func (u *User) ChangePassword(oldpass, newpass string) error {
	_oldpass, err := CryptoPass(oldpass)
	if err != nil {
		return err
	}
	_newpass, err := CryptoPass(newpass)
	if err != nil {
		return err
	}

	if u.Password != _oldpass {
		return _e("Incorrect old password")
	}

	u.Password = _newpass
	return u.Update("password")
}
