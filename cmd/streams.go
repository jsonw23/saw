package cmd

import (
	"errors"
	"fmt"

	"github.com/jsonw23/saw/blade"
	"github.com/jsonw23/saw/config"
	"github.com/spf13/cobra"
)

var streamsConfig config.Configuration

var streamsCommand = &cobra.Command{
	Use:   "streams <log group>",
	Short: "List streams in log group",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("listing streams requires log group argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		streamsConfig.Group = args[0]
		b := blade.NewBlade(&streamsConfig, nil, nil)

		if logStreams, err := b.GetLogStreams(); err == nil {
			for _, stream := range logStreams {
				fmt.Println(*stream.LogStreamName)
			}
		} else {
			panic(err)
		}
	},
}

func init() {
	streamsCommand.Flags().StringVar(&streamsConfig.Prefix, "prefix", "", "stream prefix filter")
	streamsCommand.Flags().StringVar(&streamsConfig.OrderBy, "orderBy", "LogStreamName", "order streams by LogStreamName or LastEventTime")
	streamsCommand.Flags().BoolVar(&streamsConfig.Descending, "descending", false, "order streams descending")
}
