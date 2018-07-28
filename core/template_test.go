package core

import "testing"

func TestPageUrl(t *testing.T) {
	s := PageUrl("/search?keyword=n", 1)

	if s !=  "/search?keyword=n&page=1" {
		t.Error("page url error")
	}
	s = PageUrl("/search?keyword=n&page=2", 1)

	if s !=  "/search?keyword=n&page=1" {
		t.Error("page url error")
	}
}
