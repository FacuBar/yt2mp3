package cmd

import (
	"github.com/FacuBar/yt2mp3/internal"
	"github.com/spf13/cobra"
)

// playlistCmd represents the playlist command
var playlistCmd = &cobra.Command{
	Use:   "playlist",
	Short: "A brief description",
	Long:  `A longer description`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.DownloadPlaylist(args[0])
	},
}

func init() {
	rootCmd.AddCommand(playlistCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playlistCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playlistCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
