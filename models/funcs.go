package models

import (
	"fmt"

	"github.com/cuijxin/my-n9e/pkg/i18n"
	"github.com/toolkits/pkg/str"
)

var (
	internalServerError error
	loginFailError      error
)

func InitError() {
	internalServerError = _e("Internal server error, try again later please")
	loginFailError = _e("Login fail, check your username and password")
}

func _e(format string, a ...interface{}) error {
	return fmt.Errorf(_s(format, a...))
}

func _s(format string, a ...interface{}) string {
	return i18n.Sprintf(format, a...)
}

// CryptoPass crypto password use salt
func CryptoPass(raw string) (string, error) {
	salt, err := ConfigsGet("salt")
	if err != nil {
		return "", err
	}

	return str.MD5(salt + "<-*Uk30^96eY*->" + raw), nil
}
