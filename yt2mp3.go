package yt2mp3

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kkdai/youtube/v2"
)

var (
	yt = youtube.Client{}
)

type result struct {
	Response response
	Error    error
}

type response struct {
	Video string
	Audio string
}

func downloader(urls <-chan string, downloaded chan<- result) {
	for url := range urls {
		video, err := yt.GetVideo(url)
		if err != nil {
			downloaded <- result{
				Error: fmt.Errorf("%s couldn't be fetched", url),
			}
			continue
		}

		formats := video.Formats.WithAudioChannels()
		stream, _, err := yt.GetStream(video, &formats[0])
		if err != nil {
			downloaded <- result{
				Error: fmt.Errorf("%s [%s] couldn't be downloaded", video.Title, url),
			}
			continue
		}

		videoFileName := fmt.Sprintf("%v.mp4", video.Title)
		mp3FileName := fmt.Sprintf("%v.mp3", video.Title)
		file, err := os.Create(videoFileName)
		if err != nil {
			downloaded <- result{
				Error: fmt.Errorf("%s [%s] couldn't be downloaded", video.Title, url),
			}
			continue
		}

		fmt.Println("Donwloading", video.Title)
		if _, err = io.Copy(file, stream); err != nil {
			file.Close()
			downloaded <- result{
				Error: fmt.Errorf("%s [%s] couldn't be downloaded", video.Title, url),
			}
			continue
		}

		fmt.Println("Converting to mp3")
		downloaded <- result{
			Response: response{Video: videoFileName, Audio: mp3FileName},
			Error:    nil,
		}
		file.Close()
	}
}

func converter(downloaded <-chan result, converted chan<- result) {
	for file := range downloaded {
		if file.Error != nil {
			converted <- file
		}

		convert := exec.Command("ffmpeg", "-i", file.Response.Video, file.Response.Audio)
		if err := convert.Run(); err != nil {
			fmt.Println(err)
			converted <- result{
				Error: fmt.Errorf("%s couldn't be converted", file.Response.Video),
			}
		}
		os.Remove(file.Response.Video)
		converted <- file
	}
}

func downloadAndConvert(urls []string) []result {
	home, _ := os.UserHomeDir()
	os.Chdir(filepath.Join(home, "Music"))

	urlsCh := make(chan string)
	downloaded := make(chan result)
	converted := make(chan result)

	go downloader(urlsCh, downloaded)
	go converter(downloaded, converted)

	for _, url := range urls {
		urlsCh <- url
	}

	close(urlsCh)

	return nil
}

func DownloadSingle(url string) error {
	downloadAndConvert([]string{url})
	return nil
}

// func AddTags(fileName string, tags *id3tags) error {
// 	tag, err := id3v2.Open(fileName, id3v2.Options{Parse: false})
// 	if err != nil || tag == nil {
// 		return errors.New(err.Error())
// 	}
// 	defer tag.Close()

// 	switch {
// 	case tags.CoverArt != "":
// 		client := http.Client{}

// 		resp, err := client.Get(tags.CoverArt)
// 		if err != nil {
// 			return err
// 		}
// 		defer resp.Body.Close()

// 		artwork, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			return err
// 		}

// 		tag.AddAttachedPicture(id3v2.PictureFrame{
// 			Encoding:    tag.DefaultEncoding(),
// 			MimeType:    "image/jpeg",
// 			PictureType: id3v2.PTFrontCover,
// 			Description: "Front cover",
// 			Picture:     artwork,
// 		})
// 		fallthrough

// 	case tags.Album != "":
// 		tag.SetAlbum(tags.Album)
// 		fallthrough

// 	case tags.Artist != "":
// 		tag.SetArtist(tags.Artist)
// 		fallthrough

// 	case tags.Title != "":
// 		tag.SetTitle(tags.Title)
// 	}

// 	tag.Save()

// 	return nil
// }

// type id3tags struct {
// 	CoverArt string `json:"cover_art"`
// 	Artist   string `json:"artist"`
// 	Title    string `json:"title"`
// 	Album    string `json:"album"`
// }
