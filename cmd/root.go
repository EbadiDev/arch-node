package cmd

import (
	"fmt"
	r "runtime"

	c "github.com/ebadidev/arch-node/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "arch-node",
}

func init() {
	cobra.OnInitialize(func() {
		fmt.Println(c.AppName, c.AppVersion, "(", r.Version(), r.Compiler, r.GOOS, "/", r.GOARCH, ")")
	})
}

func Execute() error {
	return rootCmd.Execute()
}
