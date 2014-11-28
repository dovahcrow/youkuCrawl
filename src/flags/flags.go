package flags

import (
	"flag"
	"os"
	"path/filepath"
)

var Cmd = flag.String("c", "", `选择命令，可用的为 syncDB-获取可用视频并存入到数据库 getVideo-将数据库中的所有视频下载到指定目录`)

//var specDown = flag.String("s", "", "下载特定的视频，输入视频id")
var DownNumber = flag.Uint("n", 1, "由于视频数量过多，请手动选择要下载的数量。按时间顺序下载。默认下载一个视频")
var Directory = flag.String("d", filepath.Dir(os.Args[0]), `要下载到的目录`)
var Verbose = flag.Bool("v", false, "日志烦琐程度 -vv 烦死你")
var Verboseverbose = flag.Bool("vv", false, "特烦琐日志模式")
var Consolelog = flag.Bool("co", true, "是否启用控制台日志")
var Filelog = flag.Bool("fo", false, "是否启用文件日志")
var DB = flag.String("db", "", "使用的数据库")
var DBPath = flag.String("dbp", "", "数据库路径,sqlite为文件路径，mysql和postgre为'root:123@tcp(127.0.0.1:3306)/test?charset=utf8',即'用户名:密码@tcp(ip地址:端口)/数据库名?charset=utf8'")

func init() {
	flag.Parse()
}
