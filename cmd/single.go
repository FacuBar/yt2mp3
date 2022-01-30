package cmd

import (
	"github.com/FacuBar/yt2mp3/internal"
	"github.com/spf13/cobra"
)

// singleCmd represents the single command
var singleCmd = &cobra.Command{
	Use:   "single",
	Short: "A brief description",
	Long:  `A longer description`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.DownloadSingle(args[0])
	},
}

func init() {
	rootCmd.AddCommand(singleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// singleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// singleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
