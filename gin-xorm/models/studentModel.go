package models

import (
	"errors"
	"fmt"
	"gin-xorm/dao"
	"time"
)

type Student struct {
	Id        int       `xorm:"pk autoincr"`
	Name      string    `json:"Name"`
	Sex       string    `json:"Sex"`
	Age       int       `json:"Age"`
	Address   string    `json:"Address"`
	CreatedAt time.Time `xorm:"created"`
}

// if err := dao.Db.Sync2(new(Student)); err != nil {
// 	log.Fatalf("Fail to sync database: %v\n", err)
// }
func (stu *Student) AddStu() error {
	_, err := dao.Db.Insert(stu)
	return err
}

func (stu *Student) GetStu() error {
	exist, err := dao.Db.Where("age = ?", 18).Get(stu)
	if !exist {
		return errors.New("无此数据")
	}
	if err != nil {
		return err
	}
	return nil
}

func (stu *Student) GetStuSlice(start int, end int) ([]Student, error) {
	stuSlice := []Student{}
	err := dao.Db.Where("id >= ? and id <= ?", start, end).Find(&stuSlice)
	fmt.Println(stuSlice)
	//对查询结果进行挨个遍历
	// err = dao.Db.Iterate(new(Student), func(i int, bean interface{}) error {
	// 	student := bean.(*Student)
	// 	fmt.Println(student.Id)
	// 	return nil
	// })

	return stuSlice, err
}

func (stu *Student) UpStu() error {
	_, err := dao.Db.Where("id = ?", 1).Update(stu)
	return err
}

func (stu *Student) DelStu() error {
	_, err := dao.Db.Delete(stu)
	return err
}
