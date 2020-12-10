package broker

import (
	"github.com/fox-one/walle/pkg/cmd/broker/create"
	"github.com/fox-one/walle/pkg/cmd/builder"
	"github.com/spf13/cobra"
)

func NewCmd(b builder.Builder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broker",
		Short: "Manager agent brokers",
	}

	cmd.AddCommand(create.NewCmd(b))
	return cmd
}
