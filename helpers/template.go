package helpers

import (
	"time"
	"html/template"
	"strings"
	"github.com/pkg/errors"
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
	}

	return exists
}

//模板传递多个变量
//{{template "userlist" Args "Users" .MostPopular "Current" .CurrentUser}}
func Args(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
