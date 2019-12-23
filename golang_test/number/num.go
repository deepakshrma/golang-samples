package number

type Num float64

func (c Num) Add(num Num) Num {
	return c + num
}
