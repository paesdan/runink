package barnctl

import (
    "github.com/spf13/cobra"
)

func newCreateCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "create [name]",
        Short: "Create a new herd",
        RunE: func(cmd *cobra.Command, args []string) error {
            if len(args) < 1 {
                return fmt.Errorf("herd name is required")
            }
            // TODO: Implement herd creation logic
            return nil
        },
    }
    return cmd
}