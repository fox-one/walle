package create

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/fox-one/walle/core"
	"github.com/fox-one/walle/pkg/cmd/builder"
	"github.com/fox-one/walle/pkg/cmdutil"
	"github.com/fox-one/walle/pkg/number"
	"github.com/spf13/cobra"
)

type Options struct {
	Name   string
	Pin    string
	UserID string
}

func NewCmd(f builder.Builder) *cobra.Command {
	opts := Options{
		UserID: f.Executor(),
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create New Broker",
		Example: heredoc.Doc(`
			$ walle broker
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := cmdutil.ValidatePin(opts.Pin); err != nil {
				if cmd.Flags().Changed("pin") {
					return err
				}

				opts.Pin = number.RandomPin()
			}

			b, err := Execute(ctx, f.Brokers(), f.Brokerz(), opts)
			if err != nil {
				return err
			}

			return f.Render().BrokerCreated(cmd.OutOrStdout(), b)
		},
	}

	cmd.Flags().StringVar(&opts.Pin, "pin", "", "broker wallet pin")
	cmd.Flags().StringVar(&opts.Name, "name", "4swap mtg broker", "broker wallet name")

	return cmd
}

func Execute(
	ctx context.Context,
	brokers core.BrokerStore,
	brokerz core.BrokerService,
	opts Options,
) (*core.Broker, error) {
	b, err := brokerz.Create(ctx, opts.Name, opts.Pin)
	if err != nil {
		return nil, err
	}

	b.UserID = opts.UserID
	if err := brokers.Create(ctx, b); err != nil {
		return nil, err
	}

	return b, nil
}
