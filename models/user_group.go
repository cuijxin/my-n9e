package models

import "github.com/toolkits/pkg/str"

type UserGroup struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Note     string `json:"note"`
	CreateAt int64  `json:"create_at"`
	CreateBy string `json:"create_by"`
	UpdateAt int64  `json:"update_at"`
	UpdateBy string `json:"update_by"`
}

func (ug *UserGroup) TableName() string {
	return "user_groups"
}

func (ug *UserGroup) Validate() error {
	if str.Dangerous(ug.Name) {
		return _e("Group name has invalid characters")
	}

	if str.Dangerous(ug.Note) {
		return _e("Group note has invalid characters")
	}

	return nil
}
