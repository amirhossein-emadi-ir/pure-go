package main

import (
	"fmt"
)

func main() {
	const name, age = "Amirhossein", 25
	fmt.Println(name, "is", age, "years old.")

	// It is conventional not to worry about any error returned by Println.
}
