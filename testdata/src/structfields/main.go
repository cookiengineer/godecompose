package main

import "fmt"

type Point struct {
	X   int
	Y   int
	Name string
	Ok   bool
}

func (p *Point) setX(n int) {
	p.X = n
}

func (p *Point) setY(n int) {
	p.Y = n
}

func (p *Point) getX() int {
	return p.X
}

func (p *Point) getY() int {
	return p.Y
}

func (p *Point) setOk(v bool) {
	p.Ok = v
}

func (p *Point) isOk() bool {
	return p.Ok
}

func (p *Point) setName(s string) {
	p.Name = s
}

func main() {
	p := &Point{X: 10, Y: 20, Name: "origin", Ok: true}

	p.setX(42)
	p.setY(100)
	fmt.Println(p.getX(), p.getY())

	p.setOk(false)
	if p.isOk() {
		fmt.Println("ok is true")
	}

	p.setName("updated")
	fmt.Println("structfields e2e done")
}
