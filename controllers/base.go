package controllers

import "github.com/gin-gonic/gin"

func defaultH(c *gin.Context) gin.H {
	options := c.MustGet("options").(map[string]string)

	return gin.H{
		"SiteTitle": options["site_name"], //page title:w
	}
}
