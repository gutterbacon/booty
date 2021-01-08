// +build mage

package main

import (
	"fmt"

	"github.com/getcouragenow/booty/tools/mage/gitr"
)

// Hello says hello
func Hello() {
	fmt.Println("hello")
}

func GitrForkCloneTemplate() error {
	return gitr.ForkCloneTemplate()
}
