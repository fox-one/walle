package migrate

import (
	"github.com/fox-one/pkg/store/db"
	"github.com/fox-one/walle/pkg/cmd/builder"
	"github.com/spf13/cobra"
)

func NewCmd(b builder.Builder) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migrate",
		Aliases: []string{"setdb"},
		Short:   "Manager DB Tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(b.DB())
		},
	}

	return cmd
}

func Run(tx *db.DB) error {
	return db.Migrate(tx)
}
