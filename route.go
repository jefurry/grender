package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strings"
)

// curl localhost:1323/render3d -H "Authorization: Bearer token" -X POST --data "bgcolor=13421772&fgcolor=205" --data-urlencode
// curl localhost:1323/render3d -H "Authorization: Bearer token" -X POST --data "bg-color=13421772&fg-color=205&stl-file=examples/cab.stl&image-path=examples/images/"
func render3d(c *gin.Context) {
	defer Logger.Flush()

	if *config.Server.Mode != gin.DebugMode {
		authorization := c.GetHeader("Authorization")
		token := strings.TrimPrefix(authorization, "Bearer ")
		Logger.Infof("Received Token for %s", token)
		//claims, err := parseToken(token)
		_, err := ParseToken(token)
		if err != nil {
			Logger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
			return
		}
		//fmt.Println(claims["iss"], claims["nbf"])
	}

	//modelId := c.PostForm("model-id")
	//size := c.PostForm("size")
	modelId := 1 //c.PostForm("model-id")
	size := 0    //c.PostForm("size")
	bgcolor := c.DefaultPostForm("bg-color", defaultBgColor)
	fgcolor := c.DefaultPostForm("fg-color", defaultFgColor)
	stlFile := c.PostForm("stl-file")
	imagePath := c.PostForm("image-path")

	if !FindPrefixInStringArray(stlFile, config.Render.StlFilePaths) || !FindPrefixInStringArray(imagePath, config.Render.ImageFilePaths) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "400 Bad Request"})
		return
	}

	name := GenName("")
	dir, err := GetHashDir(imagePath, name)
	if err != nil {
		Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server Internal Error"})
		return
	}

	imageFile := path.Join(dir, GenMd5(name)+".png")

	h, err := fauxgl_render(modelId, size, fgcolor, bgcolor, stlFile, imageFile)
	if err != nil {
		Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server Internal Error"})
		return
	}

	c.JSON(http.StatusOK, h)
}
