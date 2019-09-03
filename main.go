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
)

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}
func getHandler(c echo.Context) error {
    b := banner.GetActiveBanner(c)
    const dateLayout = "2006-01-02 15:04"
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
    return c.Render(http.StatusOK, "index", data)
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
    // Routes
    e.GET("/", getHandler)
    e.POST("/insert", banner.Insert)
    e.POST("/find", banner.Find)
    e.POST("/update", banner.Update)
    e.POST("/delete", banner.Delete)
    // Start server
    e.Logger.Fatal(e.Start(":8080"))
}
