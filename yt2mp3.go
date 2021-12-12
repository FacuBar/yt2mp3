package yt2mp3

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	// "github.com/bogem/id3v2"

	"github.com/kkdai/youtube/v2"
)

func DownloadSingle(url string) {
	client := youtube.Client{}

	video, err := client.GetVideo(url)
	if err != nil {
		panic(err)
	}

	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}

	videoFileName := fmt.Sprintf("%v.mp4", video.Title)
	mp3FileName := fmt.Sprintf("%v.mp3", videoFileName)
	file, err := os.Create(videoFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err = io.Copy(file, stream); err != nil {
		panic(err)
	}

	Convert(videoFileName, mp3FileName)
}

func Convert(videFileName, mp3FileName string) {
	convert := exec.Command("ffmpeg", "-i", videFileName, mp3FileName)
	if err := convert.Run(); err != nil {
		fmt.Println(err)
	}
	os.Remove(videFileName)
}

// func addCoverArt(url string, file string) error {
// 	client := http.Client{}

// 	resp, err := client.Get(url)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	artwork, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	tag, err := id3v2.Open(file, id3v2.Options{Parse: false})
// 	if err != nil || tag == nil {
// 		return errors.New(err.Error())
// 	}
// 	defer tag.Close()

// 	tag.AddAttachedPicture(id3v2.PictureFrame{
// 		Encoding:    tag.DefaultEncoding(),
// 		MimeType:    "image/jpeg",
// 		PictureType: id3v2.PTFrontCover,
// 		Description: "Front cover",
// 		Picture:     artwork,
// 	})

// 	tag.Save()

// 	return nil
// }
