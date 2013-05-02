package heightmap

import (
	"fmt"
)

type HeightMap struct {
	Heights   [][]float64
	dimension int
	perlin    *PerlinGenerator
}

func NewHeightMap(dimension int) *HeightMap {
	hm := &HeightMap{
		Heights:   make([][]float64, dimension),
		dimension: dimension,
		perlin:    NewPerlinGenerator(0),
	}
	for i := 0; i < dimension; i++ {
		hm.Heights[i] = make([]float64, dimension)
	}
	return hm
}

func (self *HeightMap) String() string {
	s := ""
	for i := 0; i < self.dimension; i++ {
		for j := 0; j < self.dimension; j++ {
			s += fmt.Sprintf("%v", self.Heights[i][j])
		}
	}
	return s
}

func (self *HeightMap) AddPerlinNoise(f float64) {
	for i := 0; i < self.dimension; i++ {
		for j := 0; j < self.dimension; j++ {
			self.Heights[i][j] += self.perlin.Noise(f*float64(i)/float64(self.dimension), f*float64(j)/float64(self.dimension), 0.0)
		}
	}
}

func (self *HeightMap) Perturb(f, d float64) {
	temp := make([][]float64, self.dimension)
	var u, v int
	for i := 0; i < self.dimension; i++ {
		temp[i] = make([]float64, self.dimension)
		for j := 0; j < self.dimension; j++ {
			u = i + int(self.perlin.Noise(f*float64(i)/float64(self.dimension), f*float64(j)/float64(self.dimension), 0.0)*d)
			v = j + int(self.perlin.Noise(f*float64(i)/float64(self.dimension), f*float64(j)/float64(self.dimension), 1.0)*d)
			if u < 0 {
				u = 0
			} else if u >= self.dimension {
				u = self.dimension - 1
			}
			if v < 0 {
				v = 0
			} else if v >= self.dimension {
				v = self.dimension - 1
			}
			temp[i][j] = self.Heights[u][v]
		}
	}
	self.Heights = temp
}
