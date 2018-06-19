package dbdao

import (
	"reflect"
)

func init() {
	DbStructs["student"] = reflect.TypeOf(Student{})
	DbDaos["student"] = reflect.TypeOf(StudentDao{})
}

type Student struct {
	Id        int    `xorm:"not null pk autoincr INT(11)"`
	Name      string `xorm:"not null default '' VARCHAR(100)"`
	Age       int    `xorm:"not null default 0 INT(2)"`
	Birthdate int    `xorm:"not null default 0 INT(11)"`
}

type StudentDao struct {
	DbBaseDao
}

func NewStudentDao(v ...interface{}) *StudentDao {
	this := new(StudentDao)
	this.UpdateEngine(v...)
	return this
}

func (this *StudentDao) Get(mId Param) (ret []Student, err error) {
	ret = make([]Student, 0)
	this.buildQuery(mId, "id")
	err = this.Session.Find(&ret)
	return
}
func (this *StudentDao) GetLimit(mId Param, pn, rn int) (ret []Student, err error) {
	ret = make([]Student, 0)
	this.buildQuery(mId, "id")
	err = this.Session.Limit(rn, pn).Find(&ret)
	return
}
func (this *StudentDao) GetCount(mId Param) (ret int64, err error) {
	this.buildQuery(mId, "id")
	ret, err = this.Session.Count(new(Student))
	return
}
