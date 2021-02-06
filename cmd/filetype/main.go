package main

import (
	"fmt"
	"github.com/gromgit/pour/internal/bottle"
	"os"
)

func main() {
	for _, c := range os.Args[1:] {
		if t, err := bottle.GetTypeFromPath(c); err != nil {
			fmt.Println("ERROR in " + c + ": " + err.Error())
		} else {
			fmt.Println(c + ": " + t)
		}
	}
}
