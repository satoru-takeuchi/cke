package cmd

import (
	"context"

	"github.com/cybozu-go/well"
	"github.com/spf13/cobra"
)

var sabakanDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable sabakan integration",
	Long:  `Disable sabakan integration.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		well.Go(func(ctx context.Context) error {
			return storage.EnableSabakan(ctx, false)
		})
		well.Stop()
		return well.Wait()
	},
}

func init() {
	sabakanCmd.AddCommand(sabakanDisableCmd)
}
