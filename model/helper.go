package model

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"xorm.io/xorm"
)

type beforeInsertProcess interface {
	beforeInsert() error
}

func GetById(id string, obj interface{}) error {
	has, err := db.ID(id).Get(obj)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("未找到对象")
	}

	return nil
}
func IsExist(obj interface{}) bool {
	has, _ := db.Get(obj)
	return has
}
func DeleteById(id string, obj interface{}) error {
	_, err := db.ID(id).Delete(obj)
	return err
}
func UpdateById(id string, obj interface{}, cols ...string) error {
	_, err := db.Cols(cols...).ID(id).Update(obj)
	return err
}

func UpdateByIdWithTransaction(s *xorm.Engine, id string, obj interface{}, cols ...string) error {
	_, err := s.Cols(cols...).ID(id).Update(obj)
	return err
}

func Insert(obj interface{}) error {
	var err error
	m, ok := obj.(beforeInsertProcess)
	if ok {
		err = m.beforeInsert()
		if err != nil {
			return err
		}
	}

	_, err = db.Insert(obj)
	return err
}

func checkMethodExist(m interface{}, name string) (reflect.Method, bool) {
	st := reflect.TypeOf(m)
	return st.MethodByName(name)
}

func InsertWithSession(s *xorm.Session, obj interface{}) error {
	_, err := s.Insert(obj)
	return err

}
func UpdateByIdWithSession(s *xorm.Session, id string, obj interface{}, cols ...string) error {
	_, err := s.Cols(cols...).ID(id).Update(obj)
	return err
}
func splitToInt64(s string, sep string) []int64 {
	list := strings.Split(s, sep)
	arr := make([]int64, 0)
	for _, t := range list {
		ti, _ := strconv.ParseInt(t, 10, 64)
		arr = append(arr, ti)
	}
	return arr
}

// 导出数据库
func ExportDB(path string) error {
	err := db.DumpAllToFile(path)
	return err
}

// 导入数据库
func ImportDB(path string) error {
	_, err := db.ImportFile(path)
	return err
}
