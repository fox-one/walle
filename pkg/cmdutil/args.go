package cmdutil

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func ValidatePin(pin string) error {
	const pattern = `^\d{6}$`
	if !govalidator.Matches(pin, pattern) {
		return errors.New("pin invalid")
	}

	return nil
}

func NoArgsQuoteReminder(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	errMsg := fmt.Sprintf("unknown argument %q", args[0])
	if len(args) > 1 {
		errMsg = fmt.Sprintf("unknown arguments %q", args)
	}

	hasValueFlag := false
	cmd.Flags().Visit(func(f *pflag.Flag) {
		if f.Value.Type() != "bool" {
			hasValueFlag = true
		}
	})

	if hasValueFlag {
		errMsg += "; please quote all values that have spaces"
	}

	return errors.New(errMsg)
}

func AnyArgs(cmd *cobra.Command, args []string) error { return nil }
