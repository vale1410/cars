package base


const (
	OptionType CountType = iota
	ClassType
	ExactlyOne
	OptimizationType
)

type CountType int

type CountableId struct {
	Typ   CountType
	Index int
}

type Countable struct {
	CId      CountableId
	Window   int
	Capacity int
	Demand   int
	Lower    []int
	Upper    []int
}

type PosVar struct {
	CId CountableId
	Pos int
}

type CountVar struct {
	CId   CountableId
	Pos   int
	Count int
}

type AtMostVar struct {
	CId   CountableId
	First int
	Pos   int
	Count int
}

func (c *Countable) ComputeSimpleBounds(size int) {

	c.Lower = make([]int, size)
	c.Upper = make([]int, size)

	h := c.Demand

	for i := size - 1; i >= 0; i-- {
		c.Lower[i] = h
		if h > 0 {
			h--
		}
	}

	h = 2

	for i := 0; i < size; i++ {
		c.Upper[i] = h
		if h <= c.Demand {
			h++
		}
	}
}

func (c *Countable) ComputeImprovedBounds(size int) {
	c.Lower = make([]int, size)
	c.Upper = make([]int, size)

	h := c.Demand

	for i := size - 1; i >= 0; i-- {
		q := c.Window
		u := c.Capacity

		for i >= 0 {

			c.Lower[i] = h

			if u > 0 {
				u--
				if h > 0 {
					h--
				}
			}
			q--
			if q <= 0 {
				break
			}
			i--
		}
	}

	h = 1
	q := c.Window - 1
	u := c.Capacity - 1

	for i := 0; i < size; i++ {

		for i < size {

			c.Upper[i] = h + 1

			if u > 0 && h < c.Demand {
				u--
				h++
			}
			q--
			if q <= 0 {
				break
			}
			i++
		}

		q = c.Window
		u = c.Capacity

	}
}
