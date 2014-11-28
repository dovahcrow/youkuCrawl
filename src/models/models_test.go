package models

import (
	"testing"
	//"time"
)

func TestDup(t *testing.T) {
	//v := new(Video)
	//v.Name = "aaad"
	//v.Time = time.Now()
	//v.V_Id = `LODPCDd`
	//v.VideoUrl = "http://SDa"
	//InsertVideo(v)
	//a, _ := GetAllVideo()
	//t.Log(IfVideoExist(v.V_Id))

	a, _ := GetVideoNum(1)
	t.Log(a[0])
	t.Log(UpdateVideo(a[0]))
}
