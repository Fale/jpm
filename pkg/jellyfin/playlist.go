package jellyfin

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
	"path"
)

type Playlist struct {
	XMLName           xml.Name `xml:"Item"`
	Added             string
	LockData          bool
	LocalTitle        string
	RunningTime       int
	Genres            []string        `xml:"Genres>Genre"`
	Studios           []string        `xml:"Studios>Studio"`
	PlaylistItems     []PlaylistItem  `xml:"PlaylistItems>PlaylistItem"`
	Shares            []PlaylistShare `xml:"Shares>Share"`
	PlaylistMediaType string
}

type PlaylistShare struct {
	UserId  string
	CanEdit bool
}

type PlaylistItem struct {
	Path string
}

func ReadPlaylist(p string, name string) (Playlist, error) {
	xmlFile, err := os.Open(path.Join(p, name, "playlist.xml"))
	if err != nil {
		return Playlist{}, err
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)
	var playlist Playlist

	xml.Unmarshal(byteValue, &playlist)
	return playlist, nil
}

func (playlist Playlist) Save(p string, name string) error {
	if err := os.MkdirAll(path.Join(p, name), os.ModePerm); err != nil {
		return err
	}

	xmlFile, err := os.Create(path.Join(p, name, "playlist.xml"))
	if err != nil {
		return err
	}
	xmlFile.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n")
	encoder := xml.NewEncoder(xmlFile)
	encoder.Indent("", "  ")
	err = encoder.Encode(playlist)
	if err != nil {
		return err
	}

	return nil
}

func PlaylistExists(p string, name string) bool {
	if _, err := os.Stat(path.Join(p, name, "playlist.xml")); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
