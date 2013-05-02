package heightmap

import (
	"math"
	"math/rand"
)

const gradientSizeTable = 256
const mask = gradientSizeTable - 1

type PerlinGenerator struct {
	gradients []float64
	perm      []int
}

func NewPerlinGenerator(seed int64) *PerlinGenerator {
	gen := &PerlinGenerator{
		gradients: make([]float64, gradientSizeTable*3),
		perm: []int{225, 155, 210, 108, 175, 199, 221, 144, 203, 116, 70, 213, 69, 158, 33, 252,
			5, 82, 173, 133, 222, 139, 174, 27, 9, 71, 90, 246, 75, 130, 91, 191,
			169, 138, 2, 151, 194, 235, 81, 7, 25, 113, 228, 159, 205, 253, 134, 142,
			248, 65, 224, 217, 22, 121, 229, 63, 89, 103, 96, 104, 156, 17, 201, 129,
			36, 8, 165, 110, 237, 117, 231, 56, 132, 211, 152, 20, 181, 111, 239, 218,
			170, 163, 51, 172, 157, 47, 80, 212, 176, 250, 87, 49, 99, 242, 136, 189,
			162, 115, 44, 43, 124, 94, 150, 16, 141, 247, 32, 10, 198, 223, 255, 72,
			53, 131, 84, 57, 220, 197, 58, 50, 208, 11, 241, 28, 3, 192, 62, 202,
			18, 215, 153, 24, 76, 41, 15, 179, 39, 46, 55, 6, 128, 167, 23, 188,
			106, 34, 187, 140, 164, 73, 112, 182, 244, 195, 227, 13, 35, 77, 196, 185,
			26, 200, 226, 119, 31, 123, 168, 125, 249, 68, 183, 230, 177, 135, 160, 180,
			12, 1, 243, 148, 102, 166, 38, 238, 251, 37, 240, 126, 64, 74, 161, 40,
			184, 149, 171, 178, 101, 66, 29, 59, 146, 61, 254, 107, 42, 86, 154, 4,
			236, 232, 120, 21, 233, 209, 45, 98, 193, 114, 78, 19, 206, 14, 118, 127,
			48, 79, 147, 85, 30, 207, 219, 54, 88, 234, 190, 122, 95, 67, 143, 109,
			137, 214, 145, 93, 92, 100, 245, 0, 216, 186, 60, 83, 105, 97, 204, 52},
	}

	source := rand.NewSource(0)
	rnd := rand.New(source)
	source.Seed(seed)

	for i := 0; i < gradientSizeTable; i++ {
		z := 1.0 - 2.0*rnd.Float64()
		r := math.Sqrt(1.0 - z*z)
		theta := 2.0 * math.Pi * rnd.Float64()
		gen.gradients[i*3] = r * math.Cos(theta)
		gen.gradients[i*3+1] = r * math.Sin(theta)
		gen.gradients[i*3+2] = z
	}

	return gen
}

func (self *PerlinGenerator) Noise(x, y, z float64) float64 {
	ix := int(math.Floor(x))
	fx0 := x - float64(ix)
	fx1 := fx0 - 1.0
	wx := self.smooth(fx0)

	iy := int(math.Floor(y))
	fy0 := y - float64(iy)
	fy1 := fy0 - 1.0
	wy := self.smooth(fy0)

	iz := int(math.Floor(z))
	fz0 := z - float64(iz)
	fz1 := fz0 - 1.0
	wz := self.smooth(fz0)

	vx0 := self.lattice(ix, iy, iz, fx0, fy0, fz0)
	vx1 := self.lattice(ix+1, iy, iz, fx1, fy0, fz0)
	vy0 := self.lerp(wx, vx0, vx1)

	vx0 = self.lattice(ix, iy+1, iz, fx0, fy1, fz0)
	vx1 = self.lattice(ix+1, iy+1, iz, fx1, fy1, fz0)
	vy1 := self.lerp(wx, vx0, vx1)

	vz0 := self.lerp(wy, vy0, vy1)

	vx0 = self.lattice(ix, iy, iz+1, fx0, fy0, fz1)
	vx1 = self.lattice(ix+1, iy, iz+1, fx1, fy0, fz1)
	vy0 = self.lerp(wx, vx0, vx1)

	vx0 = self.lattice(ix, iy+1, iz+1, fx0, fy1, fz1)
	vx1 = self.lattice(ix+1, iy+1, iz+1, fx1, fy1, fz1)
	vy1 = self.lerp(wx, vx0, vx1)

	vz1 := self.lerp(wy, vy0, vy1)
	return self.lerp(wz, vz0, vz1)
}

func (self *PerlinGenerator) permutate(x int) int {
	return self.perm[x&mask]
}

func (self *PerlinGenerator) index(ix, iy, iz int) int {
	return self.permutate(ix + self.permutate(iy+self.permutate(iz)))
}

func (self *PerlinGenerator) lattice(ix, iy, iz int, fx, fy, fz float64) float64 {
	index := self.index(ix, iy, iz)
	g := index * 3
	return self.gradients[g]*fx + self.gradients[g+1]*fy + self.gradients[g+2]*fz
}

func (self *PerlinGenerator) lerp(t, value0, value1 float64) float64 {
	return value0 + t*(value1-value0)
}

func (self *PerlinGenerator) smooth(x float64) float64 {
	return x * x * (3.0 - 2.0*x)
}
