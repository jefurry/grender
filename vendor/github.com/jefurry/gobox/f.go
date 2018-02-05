package gobox

import (
	"fmt"
	"math"
	"strconv"
)

func min(x, y, z float64) float64 {
	m := math.Min(x, y)

	return math.Min(m, z)
}

func max(x, y, z float64) float64 {
	m := math.Max(x, y)

	return math.Max(m, z)
}

// 毫米转厘米
func unit_cm(v float64, _type uint8) float64 {
	if _type == TYPE_SQUARE {
		return v * 1e-2
	} else if _type == TYPE_CUBE {
		return v * 1e-3
	}

	// TYPE_DIRECT
	return v * 1e-1
}

// 毫米转分米
func unit_dm(v float64, _type uint8) float64 {
	if _type == TYPE_SQUARE {
		return v * 1e-4
	} else if _type == TYPE_CUBE {
		return v * 1e-6
	}

	// TYPE_DIRECT
	return v * 1e-2
}

// 毫米转米
func unit_m(v float64, _type uint8) float64 {
	if _type == TYPE_SQUARE {
		return v * 1e-6
	} else if _type == TYPE_CUBE {
		return v * 1e-9
	}

	// TYPE_DIRECT
	return v * 1e-3
}

// 毫米转英寸
func unit_inch(v float64, _type uint8) float64 {
	if _type == TYPE_SQUARE {
		return v * (1 / (25.4 * 25.4))
	} else if _type == TYPE_CUBE {
		return v * (1 / (25.4 * 25.4 * 25.4))
	}

	// TYPE_DIRECT
	return v * (1 / 25.4)
}

func get_value(v float64, percision uint8) (float64, error) {
	value, err := Round(v, percision)
	if err != nil {
		return 0.0, err
	}

	return value, nil
}

func get_area(v float64, percision uint8) (float64, error) {
	return get_value(v, percision)
}

func get_volume(v float64, percision uint8) (float64, error) {
	return get_value(v, percision)
}

func get_bbox(l, w, h float64, percision uint8) (*BoundsBox, error) {
	length, err := get_value(l, percision)
	if err != nil {
		return nil, err
	}

	width, err := get_value(w, percision)
	if err != nil {
		return nil, err
	}

	height, err := get_value(h, percision)
	if err != nil {
		return nil, err
	}

	return &BoundsBox{
		Length: length, // 长
		Width:  width,  // 宽
		Height: height, // 高
	}, nil
}

// 四舍五入保留小数精度
func Round(f float64, precision uint8) (float64, error) {
	s := fmt.Sprintf("%%.%df", precision)
	v := fmt.Sprintf(s, f)
	vv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0.0, err
	}

	return vv, nil
}
