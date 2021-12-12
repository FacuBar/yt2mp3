# yt2mp3

Little golang app that allows you to download a youtube video as mp3, and optionally embed
ID3 tags -Cover Art, Artist ...-

# Instructions

At the moment being only linux is supported

- Instal ffmpeg `sudo apt install ffmpeg`

- Build with `go build -o yt2mp3 cmd/main.go`

- Start server with `./yt2mp3`

Send post request to `localhost:{port}/donwload` with json payload that follows the signature:

```json
{
  "yt_url": "https://www.youtube.com/watch?v=wK4vbN7g6m4",
  "id3": {
    "cover_art": "https://f4.bcbits.com/img/a2238719873_10.jpg",
    "artist": "My Brightest Diamond",
    "title": "To Pluto's Moon",
    "album": "A Thousand Shark's Teeth"
  }
}
```

the mp3 file will be saved in `~/Music`
