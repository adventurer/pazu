package cache

import (
	"log"
	"publish/models"
)

var (
	MemUsers        = make(map[int]string, 0)
	MemProject      = make(map[int]string, 0)
	MemUserHsaTable = make(map[string]models.User, 0)
	MemHealthTable  = make([]models.Health, 0)
)

func init() {
	CacheUsers()
	CacheProject()
	CacheUserHasTable()
	CacheHealthTable()
}

func CacheUsers() {
	list := make([]models.User, 0)
	err := models.Xorm.Alias("o").Find(&list)
	if err != nil {
		log.Println(err.Error())
	}
	for _, v := range list {
		MemUsers[v.Id] = v.Username
	}
}

func CacheProject() {
	list := make([]models.Project, 0)
	err := models.Xorm.Alias("o").OrderBy("Level asc").Find(&list)
	if err != nil {
		log.Println(err.Error())
	}
	for _, v := range list {
		MemProject[v.Id] = v.Name
	}
}

func CacheUserHasTable() {
	list := make([]models.User, 0)
	err := models.Xorm.Alias("o").Find(&list)
	if err != nil {
		log.Println(err.Error())
	}
	for _, v := range list {
		MemUserHsaTable[v.AuthKey] = v
	}
}

func CacheHealthTable() {
	list := make([]models.Health, 0)
	err := models.Xorm.Alias("o").Find(&list)
	if err != nil {
		log.Println(err.Error())
	}
	MemHealthTable = list
}
