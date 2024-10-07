package main

import (
	bcdedit_cmd "github.com/jc-lab/go-bcdedit/pkg/bcdedit-cmd"
	"os"
)

func main() {
	bcdedit_cmd.Main(os.Args)
}
