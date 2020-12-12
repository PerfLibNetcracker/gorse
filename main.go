package main

import (
	"fmt"
	"github.com/PerfLibNetcracker/gorse/cmd"
)
import "os"

func main() {
	var tes = os.Args[0]
	fmt.Println(tes)
	cmd.Main()
}
