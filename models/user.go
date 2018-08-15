package models

import "log"

// 新纪录
func (p *User) New() (affected int64, err error) {
	affected, err = Xorm.Alias("o").Insert(p)
	return
}

func (p *User) FindByUsername() (user User, err error) {
	user.Username = p.Username
	_, err = Xorm.Get(&user)
	if err != nil {
		log.Println(err.Error())
	}
	return
}

func (p *User) SetAccessTocken() (affected int64, err error) {
	affected, err = Xorm.Id(p.Id).Cols("auth_key").Cols("updated_at").Update(p)
	return
}
