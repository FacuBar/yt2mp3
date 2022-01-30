package yt2mp3

import (
	"fmt"
	"io"
	"log"
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

func downloader(videos <-chan *youtube.Video, downloaded chan<- result) {
	for video := range videos {
		formats := video.Formats.WithAudioChannels()
		stream, _, err := yt.GetStream(video, &formats[0])
		if err != nil {
			downloaded <- result{
				Error: fmt.Errorf("%s couldn't be downloaded", video.Title),
			}
			continue
		}

		videoFileName := fmt.Sprintf("%v.mp4", video.Title)
		mp3FileName := fmt.Sprintf("%v.mp3", video.Title)
		file, err := os.Create(videoFileName)
		if err != nil {
			downloaded <- result{
				Error: fmt.Errorf("%s couldn't be downloaded", video.Title),
			}
			continue
		}

		log.Println("Downloading", video.Title)
		if _, err = io.Copy(file, stream); err != nil {
			file.Close()
			downloaded <- result{
				Error: fmt.Errorf("%s couldn't be downloaded", video.Title),
			}
			continue
		}

		downloaded <- result{
			Response: response{Video: videoFileName, Audio: mp3FileName},
			Error:    nil,
		}
		file.Close()
		log.Println("Downloaded", video.Title)
	}
}

func converter(downloaded <-chan result, converted chan<- result) {
	for file := range downloaded {
		if file.Error != nil {
			converted <- file
		}

		log.Println("Converting", file.Response.Video)
		convert := exec.Command("ffmpeg", "-i", file.Response.Video, file.Response.Audio)
		if err := convert.Run(); err != nil {
			fmt.Println(err)
			converted <- result{
				Error: fmt.Errorf("%s couldn't be converted", file.Response.Video),
			}
		}
		os.Remove(file.Response.Video)
		converted <- file
		log.Println("Converted", file.Response.Video)
	}
}

func downloadAndConvert(videos []*youtube.Video) []result {
	home, _ := os.UserHomeDir()
	os.Chdir(filepath.Join(home, "Music"))

	videosCh := make(chan *youtube.Video)
	downloaded := make(chan result)
	converted := make(chan result)
	// @TODO: add waitgroup and close chanels
	go func() {
		for {
			<-converted
		}
	}()

	go downloader(videosCh, downloaded)
	go converter(downloaded, converted)

	for _, video := range videos {
		videosCh <- video
	}

	return nil
}

func DownloadSingle(url string) {
	video, err := yt.GetVideo(url)
	if err != nil {
		return
	}
	go downloadAndConvert([]*youtube.Video{video})
}

func DownloadPlaylist(url string) {
	playlist, _ := yt.GetPlaylist(url)

	videos := []*youtube.Video{}

	for _, video := range playlist.Videos {
		vid, err := yt.VideoFromPlaylistEntry(video)
		if err == nil {
			videos = append(videos, vid)
		}
	}

	go downloadAndConvert(videos)
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
