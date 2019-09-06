package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"testing"
	"strings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/bluemooninc/docker-go-banner/configs"
)

/*
** Config for Database table
 */
type Banner struct {
	gorm.Model
	PromotionCode string
	ContentUrl string
	StartedAt time.Time
	ExpiredAt time.Time
}

func remoteHello(domain string) string {
	// /greetingにクエリパラメータgreet=Helloを渡してGet問い合わせする
	res, err := http.Get(domain + "/")

	// エラー処理
	if err != nil {
		fmt.Println("Error")
		return "error"
	}
	defer res.Body.Close()

	// レスポンスを戻り値にする
	resStr, _ := ioutil.ReadAll(res.Body)
	return string(resStr)
}
/**
A banner is expired when the display period is over.
After a banner expires, it should not be displayed again.
 */
func TestExpiredBanner(t *testing.T) {
	// Setup
	db := configs.GormConnect()
	now := time.Now()
	st := now.Add(-time.Hour)
	ex := now.Add(-30 * time.Minute)

	// Truncate and Create
	db.Exec("TRUNCATE TABLE banners")
	db.Create(&Banner{
		PromotionCode: "ExpiredBanner",
		ContentUrl: "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_92x30dp.png",
		StartedAt: st,
ExpiredAt: ex,
})
// curl and get http body
res := remoteHello("http://localhost:8080")
// Find expect strings for promotion code
strPos := strings.Index(res, "ExpiredBanner")
if strPos > 0 {
t.Errorf("TestExpiredBanner ERROR.")
}
}

/**
A banner’s display period is the duration the banner is active on the screen.
*/
func TestSingleRecord(t *testing.T) {
	// Setup
	db := configs.GormConnect()
	now := time.Now()
	st := now.Add(-time.Hour)
	ex := now.Add(time.Hour)

	// Truncate and Create
	db.Exec("TRUNCATE TABLE banners")
	db.Create(&Banner{
		PromotionCode: "JustOnTimePromotion",
		ContentUrl: "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_92x30dp.png",
		StartedAt: st,
		ExpiredAt: ex,
	})
	// curl and get http body
	res := remoteHello("http://localhost:8080")
	// Find expect strings for promotion code
	strPos := strings.Index(res, "JustOnTimePromotion")
	if strPos == -1 {
		t.Errorf("No exist JustOnTimePromotion.")
	}
}

/**
there may be occasions where two banners are considered active. In this case, the banner with the earlier expiration should be displayed.
*/
func TestDoubleRecord(t *testing.T) {
	// Setup
	db := configs.GormConnect()
	now := time.Now()
	st := now.Add(-time.Hour)
	ex := now.Add(30 * time.Minute)
	ex2 := now.Add(time.Hour)

	// Truncate and Create
	db.Exec("TRUNCATE TABLE banners")
	db.Create(&Banner{
		PromotionCode: "ShorterExpiredPromotion",
		ContentUrl: "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_92x30dp.png",
		StartedAt: st,
		ExpiredAt: ex,
	})
	// 30 min Short than 1st record.
	db.Create(&Banner{
		PromotionCode: "LongerExpiredPromotion",
		ContentUrl: "https://upload.wikimedia.org/wikipedia/commons/thumb/f/fa/Apple_logo_black.svg/170px-Apple_logo_black.svg.png",
		StartedAt: st,
		ExpiredAt: ex2,
	})
	// curl and get http body
	res := remoteHello("http://localhost:8080")
	// Find expect strings for promotion code
	if strings.Index(res, "ShorterExpiredPromotion") == -1 {
		t.Errorf("We can't find ShorterExpiredPromotion.")
	}
	// Find unexpect strings for LongerExpiredPromotion
	if strings.Index(res, "LongerExpiredPromotion") > 0 {
		t.Errorf("We expect LongerExpiredPromotion is not find here.")
	}
}
/**
display the banner if the user has an internal IP address (10.0.0.1, 10.0.0.2), even if the current time is before the display period of the banner.
*/
func TestInternalIpAddress(t *testing.T) {
	// Setup
	db := configs.GormConnect()
	now := time.Now()
	st := now.Add(-time.Hour)
	ex := now.Add(3 * time.Hour)
	st2 := now.Add(time.Hour)
	ex2 := now.Add(2 * time.Hour)

	// Truncate and Create
	db.Exec("TRUNCATE TABLE banners")
	db.Create(&Banner{
		PromotionCode: "LongerExpiredPromotion",
		ContentUrl: "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_92x30dp.png",
		StartedAt: st,
		ExpiredAt: ex,
	})
	// 30 min Short than 1st record.
	db.Create(&Banner{
		PromotionCode: "ShorterExpiredFuturePromotion",
		ContentUrl: "https://upload.wikimedia.org/wikipedia/commons/thumb/f/fa/Apple_logo_black.svg/170px-Apple_logo_black.svg.png",
		StartedAt: st2,
		ExpiredAt: ex2,
	})
	// curl and get http body
	res := remoteHello("http://localhost:8080")
	// Find expect strings for promotion code
	if strings.Index(res, "ShorterExpiredFuturePromotion") == -1 {
		t.Errorf("We can't find ShorterExpiredFuturePromotion.")
	}
	// Find unexpect strings for LongerExpiredPromotion
	if strings.Index(res, "LongerExpiredPromotion") > 0 {
		t.Errorf("We expect LongerExpiredPromotion is not find here.")
	}
}