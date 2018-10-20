package main

import (
	"bytes"
	"fmt"
	"html/template"
)

func main()  {
	buf := &bytes.Buffer{}
	t, err := template.ParseFiles("views/layout/base2.html","views/layout/content.html")
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println(t.Name())
		t.Execute(buf, "Hello world")
		fmt.Println(buf.Len(),err)
	}
}
