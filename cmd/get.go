package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/jsonw23/saw/blade"
	"github.com/jsonw23/saw/config"
	"github.com/spf13/cobra"
)

var getConfig config.Configuration
var getOutputConfig config.OutputConfiguration

var getCommand = &cobra.Command{
	Use:   "get <log group>",
	Short: "Get log events",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("getting events requires log group argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		getConfig.Group = args[0]
		b := blade.NewBlade(&getConfig, nil, &getOutputConfig)
		if getConfig.Prefix != "" {
			if streams, err := b.GetLogStreams(); err != nil {
				if len(streams) == 0 {
					fmt.Printf("No streams found in %s with prefix %s\n", getConfig.Group, getConfig.Prefix)
					fmt.Printf("To view available streams: `saw streams %s`\n", getConfig.Group)
					os.Exit(3)
				}
				getConfig.Streams = streams
			} else {
				panic(err)
			}
		}
		b.GetEvents()
	},
}

func init() {
	getCommand.Flags().StringVar(&getConfig.Prefix, "prefix", "", "log group prefix filter")
	getCommand.Flags().StringVar(
		&getConfig.Start,
		"start",
		"",
		`start getting the logs from this point
Takes an absolute timestamp in RFC3339 format, or a relative time (eg. -2h).
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`,
	)
	getCommand.Flags().StringVar(
		&getConfig.End,
		"stop",
		"now",
		`stop getting the logs at this point
Takes an absolute timestamp in RFC3339 format, or a relative time (eg. -2h).
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`,
	)
	getCommand.Flags().StringVar(&getConfig.Filter, "filter", "", "event filter pattern")
	getCommand.Flags().BoolVar(&getOutputConfig.Pretty, "pretty", false, "print timestamp and stream name prefix")
	getCommand.Flags().BoolVar(&getOutputConfig.Expand, "expand", false, "indent JSON log messages")
	getCommand.Flags().BoolVar(&getOutputConfig.Invert, "invert", false, "invert colors for light terminal themes")
	getCommand.Flags().BoolVar(&getOutputConfig.RawString, "rawString", false, "print JSON strings without escaping")
}
