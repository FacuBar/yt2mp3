package yt2mp3

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bogem/id3v2"
	"github.com/kkdai/youtube/v2"
)

func DownloadSingle(url string, tags *id3tags) error {
	home, _ := os.UserHomeDir()
	os.Chdir(filepath.Join(home, "Music"))

	client := youtube.Client{}

	video, err := client.GetVideo(url)
	if err != nil {
		return err
	}

	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return err
	}

	videoFileName := fmt.Sprintf("%v.mp4", video.Title)
	mp3FileName := fmt.Sprintf("%v.mp3", video.Title)
	file, err := os.Create(videoFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("Donwloading", video.Title)
	if _, err = io.Copy(file, stream); err != nil {
		return err
	}

	fmt.Println("Converting to mp3")
	Convert(videoFileName, mp3FileName)

	if tags != nil {
		fmt.Println("Adding tags")
		addTags(mp3FileName, tags)
	}

	return nil
}

func Convert(videFileName, mp3FileName string) {
	convert := exec.Command("ffmpeg", "-i", videFileName, mp3FileName)
	if err := convert.Run(); err != nil {
		fmt.Println(err)
	}
	os.Remove(videFileName)
}

func addTags(fileName string, tags *id3tags) error {
	tag, err := id3v2.Open(fileName, id3v2.Options{Parse: false})
	if err != nil || tag == nil {
		return errors.New(err.Error())
	}
	defer tag.Close()

	switch {
	case tags.CoverArt != "":
		client := http.Client{}

		resp, err := client.Get(tags.CoverArt)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		artwork, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		tag.AddAttachedPicture(id3v2.PictureFrame{
			Encoding:    tag.DefaultEncoding(),
			MimeType:    "image/jpeg",
			PictureType: id3v2.PTFrontCover,
			Description: "Front cover",
			Picture:     artwork,
		})
		fallthrough

	case tags.Album != "":
		tag.SetAlbum(tags.Album)
		fallthrough

	case tags.Artist != "":
		tag.SetArtist(tags.Artist)
		fallthrough

	case tags.Title != "":
		tag.SetTitle(tags.Title)
	}

	tag.Save()

	return nil
}

type id3tags struct {
	CoverArt string `json:"cover_art"`
	Artist   string `json:"artist"`
	Title    string `json:"title"`
	Album    string `json:"album"`
}
