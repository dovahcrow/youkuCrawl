package urls

import (
	"fmt"
)

var meta = "http://www.youku.com/v_showlist/c91g0d1s2p%d.html"

func GetMetaPageURL(i int) (string, error) {
	if i < 0 || i > 20 {
		return ``, fmt.Errorf(`index illegal`)
	} else {
		return fmt.Sprintf(meta, i), nil
	}
}

var m3u8 = "http://v.youku.com/player/getM3U8/vid/%s/type/mp4/v.m3u8"

func GetVideoM3U8URL(id string) string {

	return fmt.Sprintf(m3u8, id)

}
