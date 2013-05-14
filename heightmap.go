package heightmap

import (
	"fmt"
	"math"
)

type HeightMap struct {
	width,height int
	Heights   [][]float64
	perlin    *PerlinGenerator
}

func NewHeightMap(width,height int) *HeightMap {
	hm := &HeightMap{
		width: width,
		height:height,
		Heights:   make([][]float64, width),
		perlin:    NewPerlinGenerator(0),
	}
	for x := 0; x < width; x++ {
		hm.Heights[x] = make([]float64, height)
	}
	
	return hm
}

func (self *HeightMap) String() string {
	s := ""
	for x := 0; x < self.width; x++ {
		for y := 0; y < self.height; y++ {
			s += fmt.Sprintf("%v", self.Heights[x][y])
		}
	}
	return s
}

func (self *HeightMap) Reset() {
	for x := 0; x < self.width; x++ {
		for y := 0; y < self.height; y++ {
			self.Heights[x][y] = 0.0
		}
	}
}

func (self *HeightMap) Seed(v int64) {
	self.perlin = NewPerlinGenerator(v)
}

func (self *HeightMap) AddPerlinNoise(f, scale float64) {
	for x := 0; x < self.width; x++ {
		for y := 0; y < self.height; y++ {
			self.Heights[x][y] += scale * self.perlin.Noise(f*float64(x)/float64(self.width), f*float64(y)/float64(self.height), 0.0)
		}
	}
}

func (self *HeightMap) Perturb(f, d float64) {
	temp := make([][]float64, self.width)
	var u, v int
	for x := 0; x < self.width; x++ {
		temp[x] = make([]float64, self.height)
		for y := 0; y < self.height; y++ {
			u = x + int(self.perlin.Noise(f*float64(x)/float64(self.width), f*float64(y)/float64(self.height), 0.0)*d)
			v = y + int(self.perlin.Noise(f*float64(x)/float64(self.width), f*float64(y)/float64(self.height), 1.0)*d)
			if u < 0 {
				u = 0
			} else if u >= self.width {
				u = self.width - 1
			}
			if v < 0 {
				v = 0
			} else if v >= self.height {
				v = self.height - 1
			}
			temp[x][y] = self.Heights[u][v]
		}
	}
	self.Heights = temp
}

func (self *HeightMap) Erode(smoothness float64) {
	for x := 1; x < self.width-1; x++ {
		for y := 1; y < self.height-1; y++ {
   			d_max := 0.0
   			match := make([]int, 2)
   			for u := -1; u <= 1; u++ {
    			for v := -1; v <= 1; v++ {
     				if math.Abs(float64(u)) + math.Abs(float64(v)) > 0 {
						d_i := self.Heights[x][y] - self.Heights[x + u][y + v]
      					if d_i > d_max {
       						d_max = d_i
       						match[0] = u
							match[1] = v
      					}
     				}
    			}
   			}

			// @todo why width
   			if 0 < d_max && d_max <= (smoothness / float64(self.width)) {
    			d_h := 0.5 * d_max
    			self.Heights[x][y] -= d_h
    			self.Heights[x + match[0]][y + match[1]] += d_h
   			}
  		}
 	}
}

func (self *HeightMap) Smoothen() {
	for x := 1; x < self.width-1; x++ {
		for y := 1; y < self.height-1; y++ {
   			total := 0.0
   			for u := -1; u <= 1; u++ {
    			for v := -1; v <= 1; v++ {
     				total += self.Heights[x+u][y+v]
    			}
   			}
   			self.Heights[x][y] = total / 9.0
  		}
 	}
}
