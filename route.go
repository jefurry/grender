package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// curl localhost:1323/render3d -H "Authorization: Bearer token" -X POST --data "bgcolor=13421772&fgcolor=205" --data-urlencode
// curl localhost:1323/render3d -H "Authorization: Bearer token" -X POST --data "bg-color=13421772&fg-color=205&stl-file=examples/cab.stl&image-path=examples/images/"
func render3d(c *gin.Context) {
	defer Logger.Flush()

	//fmt.Println(GenToken(500, "1", "3d@grender", "grender", "urn:grender"))
	modelId := 1
	if *config.Server.Mode != gin.DebugMode {
		authorization := c.GetHeader("Authorization")
		token := strings.TrimPrefix(authorization, "Bearer ")
		Logger.Infof("Received Token for '%s'", token)
		//claims, err := parseToken(token)
		claims, err := ParseToken(token)
		if err != nil {
			Logger.Errorf("Parse Token for: %s", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
			return
		}

		jti, ok := claims["jti"].(string)
		if !ok || jti == "" {
			Logger.Errorf("Invalid jti for '%v'", jti)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
			return
		}

		id, err := strconv.Atoi(jti)
		if err != nil {
			Logger.Errorf("Parse Token jti for: %s", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
			return
		}

		if id <= 0 {
			Logger.Errorf("Invalid modelId for '%v'", id)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
			return
		}

		modelId = id
	}

	fileSize := c.DefaultPostForm("size", "0")
	bgcolor := c.DefaultPostForm("bg-color", defaultBgColor)
	fgcolor := c.DefaultPostForm("fg-color", defaultFgColor)
	stlFile := c.PostForm("stl-file")
	imagePath := c.PostForm("image-path")

	size, err := strconv.Atoi(fileSize)
	if err != nil {
		Logger.Errorf("Parse size for: %s", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
		return
	}
	if size <= 0 {
		Logger.Errorf("Invalid size for: '%v'", fileSize)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
		return
	}

	if !FindPrefixInStringArray(stlFile, config.Render.StlFilePaths) {
		Logger.Errorf("Invalid Stl File Path for '%v'", stlFile)
		c.JSON(http.StatusBadRequest, gin.H{"message": "400 Bad Request"})
		return
	}

	if !FindPrefixInStringArray(imagePath, config.Render.ImageFilePaths) {
		Logger.Errorf("Invalid Image File Path for '%v'", imagePath)
		c.JSON(http.StatusBadRequest, gin.H{"message": "400 Bad Request"})
		return
	}

	name := GenName("")
	dir, err := GetHashDir(imagePath, name)
	if err != nil {
		Logger.Errorf("Generate dir Failed: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server Internal Error"})
		return
	}

	imageFile := path.Join(dir, GenMd5(name)+".png")

	h, err := fauxgl_render(modelId, size, fgcolor, bgcolor, stlFile, imageFile)
	if err != nil {
		Logger.Errorf("Render Faield: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server Internal Error"})
		return
	}

	Logger.Infof("Render Succed for: 'model-id=%d, file-size=%d, stl-file=%s, image-file=%s'", modelId, size, stlFile, imageFile)

	c.JSON(http.StatusOK, h)
}
