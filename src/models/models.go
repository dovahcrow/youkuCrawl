package models

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"time"
)

var Engine *xorm.Engine

func init() {
	var err error
	Engine, err = xorm.NewEngine("sqlite3", "database.db")
	if err != nil {
		log.Fatalf("init database error: %v", err)
	}

	//Engine.ShowSQL = true

	err = Engine.Sync(new(Video))
	if err != nil {
		log.Fatalf("sync database error: %v", err)
	}
}

type Video struct {
	Id       int64
	V_Id     string `xorm:"not null unique varchar(30)"`
	Name     string `xorm:"not null"`
	Time     time.Time
	VideoUrl string `xorm:"not null unique varchar(255)"`
}

var ErrExist = fmt.Errorf("duplicated")

func InsertVideo(v ...*Video) (err error) {
	_, err = Engine.Insert(v)
	if err != nil {

		if strings.Contains(err.Error(), `constraint failed`) {
			return ErrExist
		}
	}
	return
}
func GetAllVideo() (ret []*Video, err error) {
	err = Engine.Find(&ret)
	return
}

func GetVideoNum(i int) (ret []*Video, err error) {
	err = Engine.Desc("id").Limit(i).Find(&ret)
	return
}
func IfVideoExist(id string) (exist bool, err error) {
	if c, err := Engine.Where("v__id=?", id).Count(new(Video)); err == nil {
		if c > 0 {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, err
	}
}
