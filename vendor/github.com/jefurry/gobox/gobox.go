package gobox

import (
	"math"
)

var (
	INF = math.Inf(1)
)

// 顶点
type Vertex3 struct {
	X, Y, Z float64
}

// 三角形
type Triangle struct {
	V1 *Vertex3
	V2 *Vertex3
	V3 *Vertex3
}

type BoundsBox struct {
	Length, Width, Height float64
}

type Box struct {
	min_x, min_y, min_z, max_x, max_y, max_z float64
	// 三角形个数
	triangle_count int32
	// 表面积
	total_area float64
	// 体积
	total_volume float64
	// 模型名称
	model_name []byte
	// 单位转换器
	ct *Converter
}

func NewBox(ct *Converter) *Box {
	if ct == nil {
		ct = DefaultConverter
	}

	return &Box{
		min_x:          INF,
		max_x:          -INF,
		min_y:          INF,
		max_y:          -INF,
		min_z:          INF,
		max_z:          -INF,
		triangle_count: 0,
		total_area:     0.0,
		total_volume:   0.0,
		ct:             ct,
	}
}

// 设置单位转换器
func (box *Box) SetConverter(ct *Converter) {
	box.ct = ct
}

// 获取单位转换器
func (box *Box) GetConverter() *Converter {
	return box.ct
}

// 计算单个三角形的表面积
func (box *Box) AreaOfTriangle(triangle *Triangle) float64 {
	v1 := triangle.V1
	v2 := triangle.V2
	v3 := triangle.V3

	var aside float64 = math.Sqrt(math.Pow(v2.X-v1.X, 2) + math.Pow(v2.Y-v1.Y, 2) + math.Pow(v2.Z-v1.Z, 2))
	var bside float64 = math.Sqrt(math.Pow(v3.X-v2.X, 2) + math.Pow(v3.Y-v2.Y, 2) + math.Pow(v3.Z-v2.Z, 2))
	var cside float64 = math.Sqrt(math.Pow(v3.X-v1.X, 2) + math.Pow(v3.Y-v1.Y, 2) + math.Pow(v3.Z-v1.Z, 2))

	var totalSides float64 = aside + bside + cside
	var halfSides float64 = totalSides / 2
	var triangle_area float64 = math.Sqrt(halfSides * (halfSides - aside) * (halfSides - bside) * (halfSides - cside))

	return triangle_area
}

// 计算单个三角形的容量
func (box *Box) SignedVolumeOfTriangle(triangle *Triangle) float64 {
	v1 := triangle.V1
	v2 := triangle.V2
	v3 := triangle.V3

	v321 := v3.X * v2.Y * v1.Z
	v231 := v2.X * v3.Y * v1.Z
	v312 := v3.X * v1.Y * v2.Z
	v132 := v1.X * v3.Y * v2.Z
	v213 := v2.X * v1.Y * v3.Z
	v123 := v1.X * v2.Y * v3.Z

	return (1.0 / 6.0) * (-v321 + v231 + v312 - v132 - v213 + v123)
}

func (box *Box) SeekArea(triangle *Triangle) float64 {
	area := box.AreaOfTriangle(triangle)
	box.total_area += area

	return box.total_area
}

func (box *Box) SeekVolume(triangle *Triangle) float64 {
	volume := box.SignedVolumeOfTriangle(triangle)
	box.total_volume += volume

	return box.total_volume
}

// x, y, z的查找最大值和最小值
// 用于计算长宽高
func (box *Box) SeekBoundsBox(triangle *Triangle) {
	v1 := triangle.V1
	v2 := triangle.V2
	v3 := triangle.V3

	// x
	min_x := min(v1.X, v2.X, v3.X)
	max_x := max(v1.X, v2.X, v3.X)
	if min_x < box.min_x {
		box.min_x = min_x
	}
	if max_x > box.max_x {
		box.max_x = max_x
	}

	// y
	min_y := min(v1.Y, v2.Y, v3.Y)
	max_y := max(v1.Y, v2.Y, v3.Y)
	if min_y < box.min_y {
		box.min_y = min_y
	}
	if max_y > box.max_y {
		box.max_y = max_y
	}

	// z
	min_z := min(v1.Z, v2.Z, v3.Z)
	max_z := max(v1.Z, v2.Z, v3.Z)
	if min_z < box.min_z {
		box.min_z = min_z
	}
	if max_z > box.max_z {
		box.max_z = max_z
	}
}

func (box *Box) GetTriangleCount() int32 {
	return box.triangle_count
}

func (box *Box) SetTriangleCount(triangle_count int32) {
	box.triangle_count = triangle_count
}

// 表面积
// 单位默认为毫米
func (box *Box) GetArea(unit, percision uint8) (float64, error) {
	return box.ct.GetArea(math.Abs(box.total_area), unit, percision)
}

// 体积
// 单位默认为毫米
func (box *Box) GetVolume(unit, percision uint8) (float64, error) {
	return box.ct.GetVolume(math.Abs(box.total_area), unit, percision)
}

// 长宽高
// 单位默认为毫米
func (box *Box) GetBoundsBox(unit, percision uint8) (*BoundsBox, error) {
	length, err := box.ct.GetDirect(math.Abs(box.max_x-box.min_x), unit, TYPE_DIRECT)
	if err != nil {
		return nil, err
	}
	width, err := box.ct.GetDirect(math.Abs(box.max_y-box.min_y), unit, TYPE_DIRECT)
	if err != nil {
		return nil, err
	}
	height, err := box.ct.GetDirect(math.Abs(box.max_z-box.min_z), unit, TYPE_DIRECT)
	if err != nil {
		return nil, err
	}

	return box.ct.GetBBox(length, width, height, percision)
}

func (box *Box) GetMinX() float64 {
	return box.min_x
}

func (box *Box) GetMaxX() float64 {
	return box.max_x
}

func (box *Box) GetMinY() float64 {
	return box.min_y
}

func (box *Box) GetMaxY() float64 {
	return box.max_y
}

func (box *Box) GetMinZ() float64 {
	return box.min_z
}

func (box *Box) GetMaxZ() float64 {
	return box.max_z
}

func (box *Box) SetModelName(model_name []byte) []byte {
	box.model_name = model_name

	return model_name
}

func (box *Box) GetModelName() []byte {
	return box.model_name
}
