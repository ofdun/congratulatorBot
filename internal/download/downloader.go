package download

import (
	"github.com/kkdai/youtube/v2"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"telegramBot/internal/model"
	"telegramBot/internal/parser"
	"telegramBot/internal/storage"
	"time"
)

func getSuffix(str string) string {
	parts := strings.Split(str, ".")
	return "." + parts[len(parts)-1]
}

func IsVideo(path string) bool {
	videoFormats := map[string]bool{
		".mp4": true, ".mov": true, ".avi": true, ".mkv": true,
	}
	suffix := getSuffix(path)
	return videoFormats[suffix]
}

func PostcardsDownloadYoutube(path, filename, url string) error {
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
	defer func() {
		if err = stream.Close(); err != nil {
			return
		}
	}()

	out, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer func() {
		if err = out.Close(); err != nil {
			return
		}
	}()

	_, err = io.Copy(out, stream)
	if err != nil {
		return err
	}

	return nil
}

// PostcardDownload path/to/file/
func PostcardDownload(path string, postcard *model.Postcard) error {

	filename := strconv.FormatInt(time.Now().UnixNano(), 10)

	if postcard.YouTube {
		err := PostcardsDownloadYoutube(path, filename+".mp4", postcard.Href)
		if err != nil {
			return err
		}
		postcard.Downloaded = true
		postcard.Name = filename
		postcard.Path = path + filename
		return nil
	}

	suffix := getSuffix(postcard.Href)
	filename += suffix

	out, err := os.Create(path + filename)
	defer func() {
		if err = out.Close(); err != nil {
			return
		}
	}()

	if err != nil {
		return err
	}

	response, err := http.Get(postcard.Href)
	defer func() {
		if err = response.Body.Close(); err != nil {
			return
		}
	}()

	_, err = io.Copy(out, response.Body)
	postcard.Downloaded = true
	postcard.Name = filename
	postcard.Path = path + filename

	if err != nil {
		return err
	}

	return nil
}

func NightlyPostcardDownload() error {
	postcardsDatabase := storage.NewDatabase()
	postcardsStorage := storage.NewPostcardsPostgresStorage(postcardsDatabase)
	if err := postcardsStorage.ClearDatabase(); err != nil {
		return err
	}

	newParser := parser.NewParser(os.Getenv("POSTCARD_SITE"))

	if err := newParser.GetHTML(); err != nil {
		return err
	}

	holidays, err := newParser.GetHolidays()
	if err != nil {
		return err
	}

	postcards, err := newParser.GetPostcardsPages(holidays[0])

	if err != nil {
		return err
	}

	// Updating postcards ( adding links )
	for i := range postcards {
		if err = newParser.GetPostcardHref(&postcards[i]); err != nil {
			return err
		}
	}
	postcardsPath := os.Getenv("POSTCARD_PATH")
	if err = os.RemoveAll(postcardsPath); err != nil {
		return err
	}
	if err = os.Mkdir(postcardsPath, 0777); err != nil {
		return err
	}

	for i := range postcards {
		if err = PostcardDownload(postcardsPath, &postcards[i]); err != nil {
			return err
		}
	}

	for i := range postcards {
		if err = postcardsStorage.AddPostcardToStorage(&postcards[i]); err != nil {
			return err
		}
	}
	return nil
}
