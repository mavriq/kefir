package run

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RunFunc(cmd *cobra.Command, args []string) (err error) {
	fmt.Println("run called")
}
