package banner

import (
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func getDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		return nil, nil, err
	}
	return gdb, mock, nil
}

func TestCreate(t *testing.T) {
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.LogMode(true)

	r := Repository{DB: db}

	PromotionCode := "JustInTime"
	ContentUrl:= "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_92x30dp.png"
	StartedAt := time.Add(-time.Hour)
	ExpiredAt := time.Add(time.Hour)

	// Mock設定
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "banners" ("promotion_code","content_url","started_at","expired_at",) VALUES ($1,$2,$3,$4)
         RETURNING "promotion_code"."content_url"."started_at"."expired_at"`)).
		WithArgs(PromotionCode, ContentUrl, StartedAt, ExpiredAt).
		WillReturnRows(
			sqlmock.NewRows([]string{"promotion_code"}).AddRow(id))

	// 実行
	err = r.Create(id, name)
	if err != nil {
		t.Fatal(err)
	}
}