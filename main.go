/**
Golang upper v1.11
with echo, gorm modules with mysql database sample
by Yoshi Sakai
 */
package main
import (
    "html/template"
    "github.com/castaneai/gomodtest/banner"
    "io"
    "net/http"
    "github.com/labstack/echo"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    )

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

var DATABASE *gorm.DB
const dateLayout = "2006-01-02 15:04"

/*
** Config for MySQL setting
 */
func GormConnect() *gorm.DB {
    DBMS     := "mysql"
    USER     := "docker"
    PASS     := "docker"
    PROTOCOL := "tcp(mysql_host:3306)"
    DBNAME   := "test_database"
    // add parseTime option
    CONNECT := USER+":"+PASS+"@"+PROTOCOL+"/"+DBNAME+"?parseTime=true"
    db,err := gorm.Open(DBMS, CONNECT)

    if err != nil {
        panic(err.Error())
    }
    return db
}

/*
** Common return of Json
 */
func returnJson(c echo.Context, b *banner.Banner) error {
    if b != nil {
        return c.JSON(
            http.StatusOK,
            struct {
                Id int    `json:"Id"`
                PromotionCode string `json:"PromotionCode"`
                ContentUrl string `json:"ContentUrl"`
                StartedAt string `json:"StartedAt"`
                ExpiredAt string `json:"ExpiredAt"`
            }{
                Id: http.StatusOK,
                PromotionCode: b.PromotionCode,
                ContentUrl: b.ContentUrl,
                StartedAt: b.StartedAt.Format(dateLayout),
                ExpiredAt: b.ExpiredAt.Format(dateLayout),
            }, )
    } else {
        return c.JSON(
            http.StatusOK,
            struct {
                code int    `json:"Id"`
                body string `json:"body"`
            }{
                code: http.StatusOK,
                body: "ok",
            }, )
    }
    return nil
}

func getHandler(c echo.Context) error {
    // get banner when exist
    b := banner.GetActiveBanner(DATABASE, c)
    var data = struct {
        PromotionCode string
        ContentUrl    string
        StartedAt     string
        ExpiredAt     string
        RemoteAddr    string
    }{
        PromotionCode: b.PromotionCode,
        ContentUrl:    b.ContentUrl,
        StartedAt:     b.StartedAt.Format(dateLayout),
        ExpiredAt:     b.ExpiredAt.Format(dateLayout),
        RemoteAddr:    c.RealIP(),
    }
    // Template rendering
    return c.Render(http.StatusOK, "index", data)
}

func insert(c echo.Context) error{
    u := new(banner.UserParam)
    if err := c.Bind(u); err != nil {
        return err
    }
    banner.Insert(DATABASE, u)
    return returnJson(c, nil)
}
func find(c echo.Context) error{
    u := new(banner.UserParam)
    if err := c.Bind(u); err != nil {
        return err
    }
    b := banner.Find(DATABASE, u)
    return returnJson(c, b)
}
func update(c echo.Context) error{
    u := new(banner.UserParam)
    if err := c.Bind(u); err != nil {
        return err
    }
    b := banner.Update(DATABASE, u)
    return returnJson(c, b)
}

func delete(c echo.Context) error{
    u := new(banner.UserParam)
    if err := c.Bind(u); err != nil {
        return err
    }
    b := banner.Delete(DATABASE, u)
    return returnJson(c, b)
}

/*
** main loop
 */
func main() {
    t := &Template{
        templates: template.Must(template.ParseGlob("views/*.html")),
    }
    e := echo.New()
    e.Renderer = t
    DATABASE = GormConnect()
    DATABASE.AutoMigrate(&banner.Banner{})

    // Routes
    e.GET("/", getHandler)
    e.POST("/insert", insert)
    e.POST("/find", find)
    e.POST("/update", update)
    e.POST("/delete", delete)
    // Start server
    e.Logger.Fatal(e.Start(":8080"))
    defer DATABASE.Close()
}
