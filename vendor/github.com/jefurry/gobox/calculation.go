package gobox

const (
	// PLA: 1.26g/cm3
	// 默认材料密度
	DEFAULT_DENSITY = 1.04
	// 默认因子
	DEFAULT_FACTOR = 1.0
)

type Calculation struct {
	// 材料密度
	density float64
	// 单价
	price float64
	// 因子
	factor float64
}

// density 密度
// price 单价
// factor 因子
func NewCalculation(density, price, factor float64) *Calculation {
	return &Calculation{
		density: density,
		price:   price,
		factor:  factor,
	}
}

// 计算价格
// volume 体积
// percision 精度
func (c *Calculation) GetPrice(volume float64, percision uint8) (float64, error) {
	return get_value(volume*c.density*c.price*c.factor, percision)
}
