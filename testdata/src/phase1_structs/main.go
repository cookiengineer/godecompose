package main

type Point struct {
	X int
	Y int
}

func (p *Point) GetX() int {
	return p.X
}

func (p *Point) GetY() int {
	return p.Y
}

func (p *Point) SetX(x int) {
	p.X = x
}

func (p *Point) SetY(y int) {
	p.Y = y
}

func main() {
	p := &Point{X: 1, Y: 2}
	p.GetX()
	p.GetY()
	p.SetX(10)
	p.SetY(20)
}
