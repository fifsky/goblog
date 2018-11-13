package helpers

import (
	"fmt"
	"os"
	"testing"
)

func TestMd5File(t *testing.T) {
	f, err := os.Open("./util.go")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	str := Md5File(f)
	fmt.Println(str)
}