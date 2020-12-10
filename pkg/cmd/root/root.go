package root

import (
	"fmt"

	"github.com/fox-one/walle/pkg/cmd/broker"
	"github.com/fox-one/walle/pkg/cmd/builder"
	"github.com/spf13/cobra"
)

func NewCmd(f builder.Builder, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "walle <command> <subcommand> [flags]",
		Short:         "4swap mtg agent cli",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = fmt.Fprint(cmd.ErrOrStderr(), banner)
		},
	}

	cmd.AddCommand(broker.NewCmd(f))
	return cmd
}

const banner = `

   _____                                           __                                         __   
  /  |  |  ________  _  _______  ______     ______/  |_  ____   _____     ____   ____   _____/  |_ 
 /   |  |_/  ___/\ \/ \/ /\__  \ \____ \   /     \   __\/ ___\  \__  \   / ___\_/ __ \ /    \   __\
/    ^   /\___ \  \     /  / __ \|  |_> > |  Y Y  \  | / /_/  >  / __ \_/ /_/  >  ___/|   |  \  |  
\____   |/____  >  \/\_/  (____  /   __/  |__|_|  /__| \___  /  (____  /\___  / \___  >___|  /__|  
     |__|     \/               \/|__|           \/    /_____/        \//_____/      \/     \/      

4swap mtg agent by Fox.ONE Team
`
