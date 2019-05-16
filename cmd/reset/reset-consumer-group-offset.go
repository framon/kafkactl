package reset

import (
	"github.com/deviceinsight/kafkactl/operations/consumergroupoffsets"
	"github.com/spf13/cobra"
)

var offsetFlags consumergroupoffsets.ResetConsumerGroupOffsetFlags

var cmdResetOffset = &cobra.Command{
	Use:     "consumer-group-offset GROUP",
	Aliases: []string{"cgo", "offset"},
	Short:   "reset a consumer group offset",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		(&consumergroupoffsets.ConsumerGroupOffsetOperation{}).ResetConsumerGroupOffset(offsetFlags, args[0])
	},
}

func init() {
	cmdResetOffset.Flags().BoolVarP(&offsetFlags.OldestOffset, "oldest", "", false, "set the offset to oldest offset (for all partitions or the specified partition)")
	cmdResetOffset.Flags().BoolVarP(&offsetFlags.NewestOffset, "newest", "", false, "set the offset to newest offset (for all partitions or the specified partition)")
	cmdResetOffset.Flags().Int64VarP(&offsetFlags.Offset, "offset", "", -1, "set initial offset for a partition")
	cmdResetOffset.Flags().Int32VarP(&offsetFlags.Partition, "partition", "p", -1, "partition to apply the offset")
	cmdResetOffset.Flags().StringVarP(&offsetFlags.Topic, "topic", "t", offsetFlags.Topic, "topic to change offset for")
	cmdResetOffset.Flags().BoolVarP(&offsetFlags.Execute, "execute", "e", false, "execute the reset (as default only the results are displayed for validation)")
	cmdResetOffset.Flags().StringVarP(&offsetFlags.OutputFormat, "output", "o", offsetFlags.OutputFormat, "output format. One of: json|yaml")
}
