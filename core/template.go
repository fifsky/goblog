package core

import (
	"bytes"
	"html/template"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fifsky/goblog/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/ilibs/logger"
)

// 格式化时间
func DateFormat(date time.Time, layout string) string {
	return date.Format(layout)
}

// 格式化时间
func DateFormatString(date string, layout string) string {
	d, err := time.Parse(time.RFC3339, date)

	if err != nil {
		return err.Error()
	}
	return d.Format(layout)
}

func WeekDayFormat(date time.Time) string {
	m := map[time.Weekday]string{
		0: "星期日",
		1: "星期一",
		2: "星期二",
		3: "星期三",
		4: "星期四",
		5: "星期五",
		6: "星期六",
	}

	return m[date.Weekday()]
}

// 截取字符串
func Substr(source string, start, end int) string {
	rs := []rune(source)
	length := len(rs)
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}

	return string(rs[start:end])
}

func Unescaped(x string) interface{} {
	return template.HTML(x)
}

func Truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}

//返回资源路径
func StaticUrl(url ...string) string {
	if len(url) > 0 {
		return "/static/" + strings.Trim(url[0], "/")
	}

	return "/static/"
}

func PageUrl(uri string, page int) string {
	u, err := url.Parse(uri)
	if err != nil {
		return uri
	}
	uv := u.Query()
	uv.Set("page", strconv.Itoa(page))
	u.RawQuery = uv.Encode()
	return u.String()
}

func IsPage(url ...string) bool {

	if len(url) < 2 {
		return false
	}

	currurl := strings.Trim(url[0], "/")
	currurls := strings.Split(currurl, "/")
	exists := false

	for _, page := range url[1:] {
		page = strings.Trim(page, "/")
		pages := strings.Split(page, "/")
		plen := len(pages)

		if plen == 0 || len(currurls) < plen {
			continue
		}

		suburls := currurls[:plen]
		exists = strings.Join(suburls, "/") == strings.Join(pages, "/")
		if exists {
			return true
		}
	}

	return exists
}

//模板传递多个变量
//{{template "userlist" Args "Users" .MostPopular "Current" .CurrentUser}}
func Args(values ...interface{}) (map[string]interface{}) {
	if len(values)%2 != 0 {
		logger.Error("invalid dict call")
		return nil
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			logger.Error("dict keys must be strings")
			return nil
		}
		dict[key] = values[i+1]
	}
	return dict
}

var funcMap = template.FuncMap{
	"WeekDayFormat":    WeekDayFormat,
	"DateFormatString": DateFormatString,
	"DateFormat":       DateFormat,
	"Substr":           Substr,
	"Truncate":         Truncate,
	"Unescaped":        Unescaped,
	"StaticUrl":        StaticUrl,
	"IsPage":           IsPage,
	"Args":             Args,
	"PageUrl":          PageUrl,
}

type HTMLProduction struct{}

func (r HTMLProduction) Instance(name string, data interface{}) render.Render {
	htmlFiles := make([]string, 0)
	htmlFiles = append(htmlFiles, filepath.Join(config.App.Common.Path, "views/layout/base.html"))
	htmlFiles = append(htmlFiles, filepath.Join(config.App.Common.Path, "views/", name+".html"))
	var tpl = template.Must(template.New("base.html").Funcs(funcMap).Funcs(template.FuncMap{"include": tplInclude}).ParseFiles(htmlFiles...))

	// 如果没有定义css和js模板，则定义之
	if jsTpl := tpl.Lookup("header"); jsTpl == nil {
		tpl.Parse(`{{define "header"}}{{end}}`)
	}
	if cssTpl := tpl.Lookup("footer"); cssTpl == nil {
		tpl.Parse(`{{define "footer"}}{{end}}`)
	}

	return render.HTML{
		Template: tpl,
		Data:     data,
	}
}

func SetTemplate(engine *gin.Engine) {
	engine.HTMLRender = HTMLProduction{}
}

func tplInclude(file string, dot interface{}) template.HTML {
	var buffer = &bytes.Buffer{}
	tpl, err := template.New(filepath.Base(file)+".html").Funcs(funcMap).Funcs(template.FuncMap{"include":tplInclude}).ParseFiles(filepath.Join(config.App.Common.Path, "views/", file+".html"))
	if err != nil {
		logger.Errorf("parse template file(%s) error:%v\n", file, err)
		return ""
	}
	err = tpl.Execute(buffer, dot)
	if err != nil {
		logger.Errorf("template file(%s) syntax error:%v", file, err)
		return ""
	}
	return template.HTML(buffer.String())
}
