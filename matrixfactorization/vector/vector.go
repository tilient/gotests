package vector

import (
	"fmt"
	"math/rand"
)

//------------------------------------------------------------

type Vector []float64

func NewVector(length int) Vector {
	return make(Vector, length)
}

func RandomVector(length int, factor float64) Vector {
	vec := make(Vector, length)
	for ix := range vec {
		vec[ix] = factor * rand.Float64()
	}
	return vec
}

func (vec Vector) Min(vec2 Vector) Vector {
	v := NewVector(len(vec))
	for ix := range v {
		v[ix] = vec[ix] - vec2[ix]
	}
	return v
}

func (vec Vector) Print() {
	for _, v := range vec {
		fmt.Printf("%5.2f ", v)
	}
}

//------------------------------------------------------------
