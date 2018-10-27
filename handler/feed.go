package handler

import (
	"fmt"
	"time"

	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/helpers"
	"github.com/fifsky/goblog/models"
	"github.com/gorilla/feeds"
)

var FeedGet core.HandlerFunc = func(c *core.Context) core.Response {
	now := time.Now()
	options := c.GetStringMapString("options")

	feed := &feeds.Feed{
		Title:       options["site_name"],
		Link:        &feeds.Link{Href: "https://fifsky.com"},
		Description: options["site_desc"],
		Author:      &feeds.Author{Name: "fifsky", Email: "fifsky@gmail.com"},
		Created:     now,
	}

	cid := helpers.StrTo(c.DefaultQuery("cid", "0")).MustInt()

	post := &models.Posts{}
	if cid > 0 {
		post.CateId = cid
	}

	posts, err := models.PostGetList(post, 1, 10, "", "")

	if err != nil {
		return c.ErrorMessage(err)
	}

	for _, v := range posts {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       v.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://fifsky.com/article/%d", v.Id)},
			Description: v.Content,
			Author:      &feeds.Author{Name: v.NickName, Email: "fifsky@gmail.com"},
			Created:     now,
		})
	}

	err = feed.WriteAtom(c.Writer)
	if err != nil {
		return c.ErrorMessage(err)
	}
	return nil
}
