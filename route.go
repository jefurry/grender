package main

import (
	"fmt"
	"github.com/fogleman/fauxgl"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"net/http"
	"path"
	"strings"
)

const (
	scale  = 1    // optional supersampling
	width  = 1920 // output width in pixels
	height = 1080 // output height in pixels
	fovy   = 40   // vertical field of view in degrees
	near   = 1    // near clipping plane
	far    = 10   // far clipping plane
)

var (
	eye            = fauxgl.V(-3, 1, -0.75)               // camera position
	center         = fauxgl.V(0, -0.07, 0)                // view center position
	up             = fauxgl.V(0, 1, 0)                    // up vector
	light          = fauxgl.V(-0.75, 1, 0.25).Normalize() // light direction
	defaultFgColor = "#0000FF"                            // object color
	defaultBgColor = "#CCCCCC"
)

// curl localhost:1323/render3d -H "Authorization: Bearer token" -X POST --data "bgcolor=13421772&fgcolor=205" --data-urlencode
// curl localhost:1323/render3d -H "Authorization: Bearer token" -X POST --data "bg-color=13421772&fg-color=205&stl-file=examples/cab.stl&image-path=examples/images/"
func render3d(c *gin.Context) {
	if *config.Server.Mode != gin.DebugMode {
		authorization := c.GetHeader("Authorization")
		token := strings.TrimPrefix(authorization, "Bearer ")
		//claims, err := parseToken(token)
		_, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
			return
		}
		//fmt.Println(claims["iss"], claims["nbf"])
	}

	//modelId := c.PostForm("model-id")
	bgcolor := fauxgl.HexColor(c.DefaultPostForm("bg-color", defaultBgColor))
	fgcolor := fauxgl.HexColor(c.DefaultPostForm("fg-color", defaultFgColor))
	//size := c.PostForm("size")
	stlFile := c.PostForm("stl-file")
	imagePath := c.PostForm("image-path")

	if !FindPrefixInStringArray(stlFile, config.Render.StlFilePaths) || !FindPrefixInStringArray(imagePath, config.Render.ImageFilePaths) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "400 Bad Request"})
		return
	}

	// load a mesh
	mesh, err := fauxgl.LoadSTL(stlFile)
	if err != nil {
		panic(err)
	}

	// fit mesh in a bi-unit cube centered at the origin
	mesh.BiUnitCube()

	// smooth the normals
	mesh.SmoothNormalsThreshold(fauxgl.Radians(30))

	// create a rendering context
	context := fauxgl.NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(bgcolor)

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := fauxgl.LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	//matrix := mesh.Center().LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	fmt.Println(matrix)

	// use builtin phong shader
	shader := fauxgl.NewPhongShader(matrix, light, eye)
	shader.ObjectColor = fgcolor
	//shader := NewSolidColorShader(matrix, color)
	context.Shader = shader

	// render
	context.DrawMesh(mesh)

	// downsample image for antialiasing
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)

	name := GenName("")
	dir, err := GetHashDir(imagePath, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server Internal Error"})
		return
	}

	fullPath := path.Join(dir, GenMd5(name)+".png")
	// save image
	fauxgl.SavePNG(fullPath, image)

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
