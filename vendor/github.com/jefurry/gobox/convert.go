package gobox

import (
	"errors"
	"reflect"
)

const (
	// 毫米
	UNIT_MM uint8 = iota

	// 厘米
	UNIT_CM

	// 分米
	UNIT_DM

	// 米
	UNIT_M

	// 英寸
	UNIT_INCH
)

const (
	TYPE_DIRECT uint8 = iota

	// 平方单位
	TYPE_SQUARE

	// 立方单位
	TYPE_CUBE
)

var (
	DefaultConverter = NewConverter()
)

type Handler func(float64, uint8) float64

type Converter struct {
	m map[uint8]Handler
}

func NewConverter() *Converter {
	return &Converter{
		m: make(map[uint8]Handler, 20),
	}
}

// 添加单位转换器
func (c *Converter) Add(unit uint8, handler Handler) {
	c.m[unit] = handler
}

func (c *Converter) Remove(unit uint8) {
	delete(c.m, unit)
}

func (c *Converter) ValueOf(v float64, unit, _type uint8) (float64, error) {
	if unit == UNIT_MM {
		return v, nil
	}

	f, ok := c.m[unit]
	if !ok {
		return 0.0, errors.New("convert handler not exists.")
	}

	fn := reflect.ValueOf(f)
	in := []reflect.Value{
		reflect.ValueOf(v),
		reflect.ValueOf(_type),
	}

	res := fn.Call(in)
	return res[0].Float(), nil
}

func (c *Converter) GetArea(v float64, unit, percision uint8) (float64, error) {
	v, err := c.ValueOf(v, unit, TYPE_SQUARE)
	if err != nil {
		return v, err
	}

	return get_area(v, percision)
}

func (c *Converter) GetVolume(v float64, unit, percision uint8) (float64, error) {
	v, err := c.ValueOf(v, unit, TYPE_CUBE)
	if err != nil {
		return v, err
	}

	return get_volume(v, percision)
}

func (c *Converter) GetDirect(v float64, unit, percision uint8) (float64, error) {
	v, err := c.ValueOf(v, unit, TYPE_DIRECT)
	if err != nil {
		return v, err
	}

	return v, nil
}

func (c *Converter) GetBBox(length, width, height float64, percision uint8) (*BoundsBox, error) {
	return get_bbox(length, width, height, percision)
}

func init() {
	m := map[uint8]Handler{
		UNIT_CM:   unit_cm,
		UNIT_DM:   unit_dm,
		UNIT_M:    unit_m,
		UNIT_INCH: unit_inch,
	}

	for k, v := range m {
		DefaultConverter.Add(k, v)
	}
}
