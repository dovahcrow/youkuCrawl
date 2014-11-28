package youku

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/httplib"
	"models"
	"net/http"
	"regexp"
	"time"
	"youku/urls"
)

func GetVideoIdListRange() (ret []*models.Video, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	for i := 1; i <= 20; i++ {
		u, _ := urls.GetMetaPageURL(i)
		re, err := http.Get(u)
		if err != nil {
			return nil, fmt.Errorf("Get Page %d Error: %v", i, err)
		}
		doc, err := goquery.NewDocumentFromResponse(re)
		if err != nil {
			return nil, fmt.Errorf("Get Page %d Error: %v", i, err)
		}

		doc.Find(`.yk-row>div`).Each(func(i int, s *goquery.Selection) {
			us, ok := s.Find(`.v-link>a`).Attr(`href`)
			if !ok {
				panic(fmt.Errorf("Parse Page %d Error", i))
			}
			title, ok := s.Find(`.v-link>a`).Attr(`title`)
			if !ok {
				panic(fmt.Errorf("Parse Page %d Error", i))
			}
			reg := regexp.MustCompile(`http://v\.youku\.com/v_show/id_(\w+)\.html`)
			ids := reg.FindStringSubmatch(us)
			if len(ids) < 2 {
				panic(
					fmt.Errorf(
						"Parse Page %d Error: Can't Find Video Id In '%v'", i, us,
					))
			}
			ret = append(ret, &models.Video{Name: title, V_Id: ids[1], Time: time.Now()})
		})
		re.Body.Close()
	}
	return
}

var ErrVideoEncrypted = fmt.Errorf("Video is Encrpyted")

func GetVideoUrl(video *models.Video) (err error) {
	u := urls.GetVideoM3U8URL(video.V_Id)
	m3u8, err := httplib.Get(u).String()
	if err != nil {
		return err
	}
	
	fs := regexp.MustCompile(
		`(http://.+\.(?:mp4|flv))\.ts\?ts_start=\d+&ts_end=\d+&ts_seg_no=\d+`,
	).FindStringSubmatch(m3u8)
	if len(fs) < 2 {
		if m3u8 == `` {
			return ErrVideoEncrypted
		} else {
			return fmt.Errorf("Can't Find Video Url: %s", m3u8)
		}
	}
	video.VideoUrl = fs[1]
	return
}
