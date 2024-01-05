package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"unicode"

	"github.com/dhowden/tag"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/fale/jpm/pkg/jellyfin"
)

type song struct {
	Path   string
	Artist string
	Title  string
}

func normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return strings.ToLower(result)
}

func backupPlaylist(name string) error {
	src := path.Join(playlistsPath, name, "playlist.xml")
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(backupLocation)
	if err == nil {
		return fmt.Errorf("File %s already exists.", backupLocation)
	}

	destination, err := os.Create(backupLocation)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, 100)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

func sortFunc(name string) error {
	if len(backupLocation) > 0 {
		if err := backupPlaylist(name); err != nil {
			return err
		}
	}

	p, err := jellyfin.ReadPlaylist(playlistsPath, name)
	if err != nil {
		return err
	}
	songs := make([]song, len(p.PlaylistItems))

	for i, s := range p.PlaylistItems {
		f, err := os.Open(s.Path)
		if err != nil {
			return err
		}
		defer f.Close()

		m, err := tag.ReadFrom(f)
		if err != nil {
			return err
		}

		songs[i] = song{
			Path:   s.Path,
			Title:  m.Title(),
			Artist: m.Artist(),
		}
	}
	sort.Slice(songs, func(i, j int) bool {
		if songs[i].Artist != songs[j].Artist {
			return normalize(songs[i].Artist) < normalize(songs[j].Artist)
		}
		return normalize(songs[i].Title) < normalize(songs[j].Title)
	})
	paths := make([]jellyfin.PlaylistItem, len(songs))
	for i, s := range songs {
		paths[i] = jellyfin.PlaylistItem{Path: s.Path}
	}
	p.PlaylistItems = paths

	if err := p.Save(playlistsPath, name); err != nil {
		return err
	}
	return nil
}
