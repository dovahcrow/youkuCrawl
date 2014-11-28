package main

import (
	"flag"
	"github.com/astaxie/beego/logs"
	//"log"
	"fmt"
	//"github.com/astaxie/beego/httplib"
	"flags"
	"io/ioutil"
	"models"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"youku"
)

func main() {
	//等待1秒，让log中的缓冲全部输出到屏幕上
	defer time.Sleep(1 * time.Second)

	fmt.Println(`----Youku Getter By 42. version: 1.41----`)

	//设置日志管理器
	l := logs.NewLogger(1024)  //缓冲长度1024
	l.SetLevel(logs.LevelInfo) //默认等级Warn
	flag.Parse()
	if *flags.Consolelog { //设置console日志
		l.SetLogger(`console`, ``)
	}
	if *flags.Filelog { //设置file日志
		l.SetLogger(`file`, `{"filename":"log.log"}`)
	}
	if *flags.Verbose {
		l.SetLevel(logs.LevelDebug)
	}
	if *flags.Verboseverbose {
		l.SetLevel(logs.LevelTrace)
	}

	//选择命令
	switch *flags.Cmd {
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
			directory := flags.Directory
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
			if *flags.DownNumber == 0 {
				l.Error(`Download Number Is 0！`)
				return
			}

			vs, err := models.GetVideoNum(int(*flags.DownNumber))
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

				file :=
					*directory +
						`/` +
						func() string {
							r := []rune(v.Name)
							if len(r) > 20 {
								return string(r[:20])
							} else {
								return v.Name
							}
						}() +
						".flv"

				_, err := os.Stat(file)
				if !os.IsNotExist(err) {
					l.Info("Video %v Exited, Skip", v.V_Id)
					continue
				}

				re, err := http.Get(v.VideoUrl)
				if err != nil {
					l.Critical("Download Video %v Error: %v", v.V_Id, err)
					continue
				}
				if re.StatusCode != 200 {
					//视频地址会过期，重新更新
					l.Warn("Get Video %v Error, Maybe Out Of Date", v.V_Id)

					if re.StatusCode == 404 {

						l.Warn("Its 404. Re-Get It")
						re.Body.Close()
						err = youku.GetVideoUrl(v)
						if err != nil {
							l.Critical("Re-get Video %v URL Error: %v", v.V_Id, err)
							continue
						}
						err = models.UpdateVideo(v)
						if err != nil {
							l.Critical("Update Video %v URL Error: %v", v.V_Id, err)
							continue
						}
						re, err = http.Get(v.VideoUrl)
						if err != nil {
							l.Critical("Download Video %v Error: %v", v.V_Id, err)
							continue
						}

					} else {
						l.Critical("Download Video %v Error Code: %v", v.V_Id, re.StatusCode)
					}
				}

				b, _ := ioutil.ReadAll(re.Body)
				re.Body.Close()

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
