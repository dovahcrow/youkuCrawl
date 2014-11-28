package models

import (
	"flags"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"time"
)

var Engine *xorm.Engine

func init() {
	var err error
	switch *flags.DB {
	case `sqlite`:
		{

			Engine, err = xorm.NewEngine("sqlite3", *flags.DBPath)
		}
	case `mysql`:
		{
			Engine, err = xorm.NewEngine("mysql", *flags.DBPath)
		}
	case `postgre`:
		{
			Engine, err = xorm.NewEngine("postgre", *flags.DBPath)
		}
	default:
		{
			log.Fatalf("Unsupported database. Candidates are 'sqlite','mysql','postgre'")
		}
	}

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
	V_Id     string    `xorm:"not null unique varchar(30)"`
	Name     string    `xorm:"not null"`
	Time     time.Time `xorm:"created"`
	VideoUrl string    `xorm:"not null unique varchar(255)"`
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

func UpdateVideo(v *Video) (err error) {
	_, err = Engine.Id(v.Id).Update(v)
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

func IfVideoExistName(name string) (exist bool, err error) {
	if c, err := Engine.Where("name like ?", name).Count(new(Video)); err == nil {
		if c > 0 {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, err
	}
}
