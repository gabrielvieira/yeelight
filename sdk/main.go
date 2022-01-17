package main

import (
	yeelight "github.com/gabrielvieira/yeelight/sdk/pkg"
)

func main() {
	y := yeelight.New("1", "1", "192.168.15.58:55443", "12", []string{"asd"})
	y.Toggle()
}
