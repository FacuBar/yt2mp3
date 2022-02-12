# yt2mp3

CLI app writen in Go that allows you to download videos and playlists as mp3

## Usage

- ffmpeg required `sudo apt install ffmpeg`
- get repo `git clone https://github.com/FacuBar/yt2mp3.git && cd yt2mp3`
- get dependencies `go get .`
- build `go build -o yt2mp3 main.go`

If you want to install it you can either move it to bin or install it through go tool
- `sudo mv yt2mp3 /bin`
or
- `go install`

###### commands
- `yt2mp3 single [url]` downloads a single video and converts it to mp3
- `yt2mp3 playlist [url]` downloads all videos from a playlist and converts each one to mp3

### TODO
- [x] Download playlists
- [x] Download and convert files in a concurrent manner when possible
- [x] CLI
- [ ] Add interactive and selective id3 tags insertion
- [ ] Download metadata and complete id3 tags automatically (only if I find external API of my liking)
