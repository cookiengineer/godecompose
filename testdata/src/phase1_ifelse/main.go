package main

func checkVal(v int) string {
	if v == 0 {
		return "zero"
	} else if v > 0 {
		return "positive"
	} else {
		return "negative"
	}
}

func main() {
	checkVal(-1)
}
