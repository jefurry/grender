package gobox

type Transform struct {
	area, volume float64
	bbox         *BoundsBox
}

func NewTransform(area, volume float64, bbox *BoundsBox) *Transform {
	return &Transform{
		area:   area,
		volume: volume,
		bbox:   bbox,
	}
}

// 缩放
// length_ratio 长度缩放比例
// width_ratio 宽度缩放比例
// height_ratio 调试缩放比例
func (t *Transform) Scale(length_ratio, width_ratio, height_ratio float64) {
	t.SetLength(length_ratio * t.bbox.Length)
	t.SetWidth(width_ratio * t.bbox.Width)
	t.SetHeight(height_ratio * t.bbox.Height)
}

func (t *Transform) SetLength(l float64) {
	ratio := l / t.bbox.Length
	t.bbox.Length = l
	t.bbox.Width = t.bbox.Width * ratio
	t.bbox.Height = t.bbox.Height * ratio
	t.area = t.area * ratio * ratio
	t.volume = t.volume * ratio * ratio * ratio
}

func (t *Transform) SetWidth(w float64) {
	ratio := w / t.bbox.Width
	t.bbox.Width = w
	t.bbox.Length = t.bbox.Length * ratio
	t.bbox.Height = t.bbox.Height * ratio
	t.area = t.area * ratio * ratio
	t.volume = t.volume * ratio * ratio * ratio
}

func (t *Transform) SetHeight(h float64) {
	ratio := h / t.bbox.Height
	t.bbox.Height = h
	t.bbox.Length = t.bbox.Length * ratio
	t.bbox.Width = t.bbox.Width * ratio
	t.area = t.area * ratio * ratio
	t.volume = t.volume * ratio * ratio * ratio
}

func (t *Transform) GetArea(percision uint8) (float64, error) {
	return get_area(t.area, percision)
}

func (t *Transform) GetVolume(percision uint8) (float64, error) {
	return get_volume(t.volume, percision)
}

func (t *Transform) GetBBox(percision uint8) (*BoundsBox, error) {
	return get_bbox(t.bbox.Length, t.bbox.Width, t.bbox.Height, percision)
}
