package models

// 新纪录
func (p *Health) New() (affected int64, err error) {
	affected, err = Xorm.Alias("o").Insert(p)
	return
}
