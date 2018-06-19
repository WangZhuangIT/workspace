package model

import (
	"go_libs/dao/dbdao"
	"go_libs/dao/redisdao"
	"go_libs/service/logger"
)

type Student struct {
	Id        int
	Name      string
	Age       int
	BirthDate int
}

func init() {

}

func (stu *Student) StuInfoList() (list []dbdao.Student, err error) {
	students, es := dbdao.NewStudentDao().Get(nil)
	if es != nil {
		err = es
		return
	}
	redao := redisdao.GetInstance("cache")
	conn := redao.Get()
	defer conn.Close()
	key := "demo-student"
	_, e := conn.Do("set", key, "demo-student set redis ok!")
	if e != nil {
		err = logger.NewRedisError(e)
		return
	}
	cache := redisdao.NewRedisCache("demo/cache")
	cache.Id(key).Value("demo-student set redis obj ok!").Set()
	return students, nil
}
