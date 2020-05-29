package main

import (
	"fmt"
	"os"

	"github.com/liamhendricks/auth-backend/src/cmd"
)

func main() {
	//update cmd import to your path
	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
