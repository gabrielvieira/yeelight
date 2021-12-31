package main

import "yeelight/pkg/yeelight"

func main() {
	y := yeelight.New("1", "1", "192.168.15.58:55443", "12", []string{"asd"})
	y.SetPower(true)

}