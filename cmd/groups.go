package cmd

import (
	"fmt"

	"github.com/jsonw23/saw/blade"
	"github.com/jsonw23/saw/config"
	"github.com/spf13/cobra"
)

// TODO: colorize based on logGroup prefix (/aws/lambda, /aws/kinesisfirehose, etc...)
var groupsConfig config.Configuration

var groupsCommand = &cobra.Command{
	Use:   "groups",
	Short: "List log groups",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		b := blade.NewBlade(&groupsConfig, nil, nil)
		if logGroups, err := b.GetLogGroups(); err == nil {
			for _, group := range logGroups {
				fmt.Println(*group.LogGroupName)
			}
		} else {
			panic(err)
		}
	},
}

func init() {
	groupsCommand.Flags().StringVar(&groupsConfig.Prefix, "prefix", "", "log group prefix filter")
}
