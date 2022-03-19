package main

import "fmt"

type TwitterAuth struct {
	bearer string
}

func (auth TwitterAuth) DebugText() string {
	return auth.bearer
}

func wtf() {
	fmt.Println("Wtf wtf wtf")
}
