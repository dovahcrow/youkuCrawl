package youku

import (
	"path/filepath"
	"testing"
)

func TestGetVideoList(t *testing.T) {
	t.SkipNow()
	a, _ := GetVideoIdListRange()
	t.Log(a)

}
func TestGetVideo(t *testing.T) {
	t.SkipNow()
	vs, _ := GetVideoIdListRange()
	if err := GetVideoUrl(vs[4]); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", vs[4])
}
func TestFilePath(t *testing.T) {
	t.Log(filepath.EvalSymlinks("~/"))
}
