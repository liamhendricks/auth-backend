package cmd

import (
  "fmt"
  "github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
    Use:   "api",
    Short: "Root command for api.",
}

func init() {
  fmt.Println("root init")
}
