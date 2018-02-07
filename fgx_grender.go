package main

import (
	"github.com/fogleman/fauxgl"
	"github.com/gin-gonic/gin"
	"github.com/jefurry/gobox"
	"github.com/nfnt/resize"
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

// curl localhost:1323/render3d -H "Authorization: Bearer token" -X POST --data "bgcolor=#cccccc&fgcolor=#0000ff" --data-urlencode
// curl localhost:1323/render3d -H "Authorization: Bearer token" -X POST --data "model-id=1&size=1&bg-color=#cccccc&fg-color=#0000ff&stl-file=examples/cab.stl&image-path=examples/images/"
func fauxgl_render(modelId, size int, fgcolor, bgcolor, stlFile, imageFile string) (gin.H, error) {
	fc := fauxgl.HexColor(fgcolor)
	bc := fauxgl.HexColor(bgcolor)

	// load a mesh
	mesh, err := fauxgl.LoadSTL(stlFile)
	if err != nil {
		return nil, err
	}

	// calc size, volume etc.
	box := gobox.NewBox(gobox.DefaultConverter)
	var triangle_count int32 = 0
	for _, t := range mesh.Triangles {
		p1 := t.V1.Position
		p2 := t.V2.Position
		p3 := t.V3.Position

		v1 := &gobox.Vertex3{X: p1.X, Y: p2.Y, Z: p3.Z}
		v2 := &gobox.Vertex3{X: p2.X, Y: p2.Y, Z: p2.Z}
		v3 := &gobox.Vertex3{X: p3.X, Y: p3.Y, Z: p3.Z}

		triangle := &gobox.Triangle{V1: v1, V2: v2, V3: v3}
		triangle_count += 1
		box.SetTriangleCount(triangle_count)
		/*
			// suface area
			box.SeekArea(triangle)
			// volume
			box.SeekVolume(triangle)
		*/
		// box
		box.SeekBoundsBox(triangle)
	}

	var unit, percision uint8 = gobox.UNIT_CM, 2
	ct := gobox.DefaultConverter
	/*
		fmt.Println(box.GetArea(unit, percision))
		fmt.Println(box.GetVolume(unit, percision))
	*/
	bbox, err := box.GetBoundsBox(unit, percision)
	if err != nil {
		return nil, err
	}

	volume, err := ct.GetVolume(mesh.Volume(), unit, percision)
	if err != nil {
		return nil, err
	}

	area, err := ct.GetArea(mesh.SurfaceArea(), unit, percision)
	if err != nil {
		return nil, err
	}
	// end

	// fit mesh in a bi-unit cube centered at the origin
	mesh.BiUnitCube()

	// smooth the normals
	mesh.SmoothNormalsThreshold(fauxgl.Radians(30))

	// create a rendering context
	context := fauxgl.NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(bc)

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := fauxgl.LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	//matrix := mesh.Center().LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// use builtin phong shader
	shader := fauxgl.NewPhongShader(matrix, light, eye)
	shader.ObjectColor = fc
	//shader := NewSolidColorShader(matrix, fc)
	context.Shader = shader

	// render
	context.DrawMesh(mesh)

	// downsample image for antialiasing
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)

	// save image
	if err := fauxgl.SavePNG(imageFile, image); err != nil {
		return nil, err
	}

	return gin.H{"data": gin.H{
		"model-id":   modelId,
		"size":       size,
		"image-file": imageFile,
		"volume":     volume,
		"area":       area,
		"box": gin.H{
			"length": bbox.Length,
			"width":  bbox.Width,
			"height": bbox.Height,
		},
	}}, nil
}
