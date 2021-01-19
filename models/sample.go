package models

import (
	"share/database"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Sample struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Email    string
	Password string
}

func init() {
	logrus.Info("models init...")
	database.MysqlORMConn.AutoMigrate(&Sample{})
	// if !database.MysqlORMConn.HasTable(Person{}) {
	// 	logrus.Infof("create table person")
	// 	database.MysqlORMConn.CreateTable(Person{})
	// }
}

func (s Sample) Insert() (uint, error) {
	if ret := database.MysqlORMConn.Create(&s).Error; ret != nil {
		return 0, ret
	}
	return s.ID, nil
}

func (s *Sample) Delete() error {
	if ret := database.MysqlORMConn.Delete(s).Error; ret != nil {
		return ret
	}

	return nil
}

func (s *Sample) Get(id uint) error {
	if ret := database.MysqlORMConn.First(s, id).Error; ret != nil {
		return ret
	}
	return nil
}

func (s *Sample) GetAll() (rets []Sample, err error) {
	err = database.MysqlORMConn.Find(&rets).Error
	return
}

func (s *Sample) Update() error {
	var s2 Sample
	if ret := database.MysqlORMConn.Select([]string{"id", "name"}).Where("name = ?", s.Name).Find(&s2).Error; ret != nil {
		return ret
	}
	if ret := database.MysqlORMConn.Model(&s2).Update(s).Error; ret != nil {
		return ret
	}
	logrus.Infof("updated sample: %v", s2)
	return nil
}
