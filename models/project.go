package models

import (
	"log"
)

func (p *Project) List() (list []Project) {
	list = make([]Project, 0)
	err := Xorm.Alias("o").OrderBy("Level asc").Find(&list)
	if err != nil {
		log.Println(err.Error())
	}
	return
}

func (p *Project) Find(id interface{}) (record *Project, err error) {
	record = new(Project)
	_, err = Xorm.Alias("o").Where("id=?", id).Get(record)
	if err != nil {
		log.Println(err.Error())
	}
	return
}

// 新纪录
func (p *Project) New() (affected int64, err error) {
	affected, err = Xorm.Alias("o").Insert(p)
	return
}

// 修改
func (p *Project) Edit() (affected int64, err error) {
	affected, err = Xorm.Id(p.Id).Update(p)
	return
}

// 删除
func (p *Project) Del() (affected int64, err error) {
	affected, err = Xorm.Id(p.Id).Delete(p)
	return
}
