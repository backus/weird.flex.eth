package main

import (
	"fmt"
)

func main() {
	twitterAuth := TwitterAuth{"asdf"}

	fmt.Println("Hello world!")
	fmt.Printf("Twitter auth debug = %s\n", twitterAuth.DebugText())
}
