// +build mage

package main

import (
	"fmt"
	// mage:import gitr
	_ "github.com/getcouragenow/booty/tools/mage/gitr"
)

// Hello says hello
func Hello() {
	fmt.Println("hello")
}
