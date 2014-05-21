package main

import (
	"flag"
	"github.com/astaxie/beego/logs"
	//"log"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"models"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"youku"
)

var cmd = flag.String("c", "", `选择命令，可用的为 syncDB-获取可用视频并存入到数据库 getVideo-将数据库中的所有视频下载到指定目录`)

//var specDown = flag.String("s", "", "下载特定的视频，输入视频id")
var downNumber = flag.Uint("n", 1, "由于视频数量过多，请手动选择要下载的数量。按时间顺序下载。默认下载一个视频")
var directory = flag.String("d", filepath.Dir(os.Args[0]), `要下载到的目录`)
var verbose = flag.Bool("v", false, "日志烦琐程度 -vv 烦死你")
var verboseverbose = flag.Bool("vv", false, "特烦琐日志模式")
var consolelog = flag.Bool("co", true, "是否启用控制台日志")
var filelog = flag.Bool("fo", false, "是否启用文件日志")

func main() {
	//等待1秒，让log中的缓冲全部输出到屏幕上
	defer time.Sleep(1 * time.Second)

	fmt.Println(`----Youku Getter By 42. version: 1.20----`)

	//设置日志管理器
	l := logs.NewLogger(1024)  //缓冲长度1024
	l.SetLevel(logs.LevelInfo) //默认等级Warn
	flag.Parse()
	if *consolelog { //设置console日志
		l.SetLogger(`console`, ``)
	}
	if *filelog { //设置file日志
		l.SetLogger(`file`, `{"filename":"log.log"}`)
	}
	if *verbose {
		l.SetLevel(logs.LevelDebug)
	}
	if *verboseverbose {
		l.SetLevel(logs.LevelTrace)
	}

	//选择命令
	switch *cmd {
	case `syncDB`:
		{

			l.Trace(`Getting All Video Lists`)
			videos, err := youku.GetVideoIdListRange() //获取视频列表
			if err != nil {
				l.Critical(`Get Video List Critical! System Exit`)
				os.Exit(1)
			}
			l.Info(`Get Video List Successful`)

			l.Trace(`Get Video Download Urls`)
			for _, v := range videos {

				l.Debug("Get Video Id: %s And Name: %s", v.V_Id, v.Name)

				//从数据库中检查是否该视频已经下载过了
				if ex, err := models.IfVideoExist(v.V_Id); err != nil {
					l.Warn("Get Video: %s Fail: %v", v.V_Id, err)
					continue
				} else if ex {
					l.Info("Video: %s Existed. Skip", v.V_Id)
					continue
				}

				//获取视频下载地址
				err = youku.GetVideoUrl(v)
				if err != nil {
					l.Warn("Get Video Id: %s Fail: %v", v.V_Id, err)
					if err == youku.ErrVideoEncrypted {
						l.Error("Oh The Video With Id: %s Is Encrypted!", v.V_Id)
					}
					continue
				}
			}
			l.Info(`Get Video Download Url Succssful`)

			//插入数据库
			l.Trace(`Start Inserting To Database`)
			for _, v := range videos {
				err = models.InsertVideo(v)
				if err != nil {
					if err != models.ErrExist {
						l.Critical("Insert Video With Id: %s FAIL: %v", v.V_Id, err)
						continue
					} else {
						l.Warn("Video %s Exited. Skip", v.V_Id)
					}
				}
			}
			l.Info(`Insert Video To Database Succssful`)
		}
	case `getVideo`:
		{
			//分析目录是否合法
			*directory = regexp.MustCompile(`^~/`).ReplaceAllString(*directory, os.Getenv(`HOME`)+"/") //把目录中的～用home目录替换
			//转换成绝对路径
			dir, err := filepath.Abs(*directory)
			if err != nil {
				l.Critical(`Directory Not Correct`)
				return
			}
			//去掉右斜杠
			dir = strings.TrimRight(dir, "/")
			l.Trace("Download Directory Is %s", dir)

			*directory = dir
			if *downNumber == 0 {
				l.Error(`Download Number Is 0！`)
				return
			}

			vs, err := models.GetVideoNum(int(*downNumber))
			if err != nil {
				l.Critical(`Get Videos Error`)
				return
			}
			f, err := os.Open(*directory)
			if err != nil {
				l.Critical("Cannot Open Directory: %v", err)
				return
			}
			stat, _ := f.Stat()
			if !stat.IsDir() {
				l.Critical("Target Not A Directory")
				return
			}
			f.Close()

			l.Trace(`Start Download Videos`)
			for _, v := range vs {

				l.Debug("Start Download Video: %s", v.V_Id)
				file := *directory + `/` + v.V_Id + ".flv"
				_, err := os.Stat(file)
				if !os.IsNotExist(err) {
					l.Info("Video %v Exited, Skip", v.V_Id)
					continue
				}

				b, err := httplib.Get(v.VideoUrl).Bytes()
				if err != nil {
					l.Critical("Download Video %v Error: %v", v.V_Id, err)
					continue
				}
				l.Info("Download Video %v Successful", v.V_Id)

				l.Trace("Start Writing %v To File", v.V_Id)
				f, err = os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
				if err != nil {
					l.Critical("Open File %v Error: %v", file, err)
					continue
				}
				f.Write(b)
				f.Close()
				l.Info("Write File: %v Successful", file)
			}
			l.Info(`Download Finished`)
		}
	default:
		{
			flag.Usage()
		}
	}

}
