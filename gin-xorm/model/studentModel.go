package model

import (
	"gin-xorm/dao"
)

type Student struct {
	Name    string `json:"Name"`
	Sex     string `json:"Sex"`
	Age     int    `json:"Age"`
	Address string `json:"Address"`
}

// if err := dao.Db.Sync2(new(Student)); err != nil {
// 	log.Fatalf("Fail to sync database: %v\n", err)
// }
func (stu *Student) AddStuModel() error {
	_, err := dao.Db.Insert(stu)
	return err
}
