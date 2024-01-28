package download

import (
	"CongratulatorBot/internal/model"
	"CongratulatorBot/internal/parser"
	"CongratulatorBot/internal/storage"
	"github.com/kkdai/youtube/v2"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func PostcardsDownloadYoutube(path, filename, url string, errorChan chan error) {
	client := youtube.Client{}

	video, err := client.GetVideo(url)

	if err != nil {
		errorChan <- err
	}

	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])

	if err != nil {
		errorChan <- err
	}
	defer func() {
		if err = stream.Close(); err != nil {
			errorChan <- err
		}
	}()

	out, err := os.Create(path + filename)
	if err != nil {
		errorChan <- err
	}
	defer func() {
		if err = out.Close(); err != nil {
			errorChan <- err
		}
	}()

	_, err = io.Copy(out, stream)
	if err != nil {
		errorChan <- err
	}
}

// PostcardDownload path/to/file/
func PostcardDownload(path string, postcard *model.Postcard, errorChan chan error) {

	filename := strconv.FormatInt(time.Now().UnixNano(), 10)

	if postcard.YouTube {
		PostcardsDownloadYoutube(path, filename+".mp4", postcard.Href, errorChan)
		postcard.Downloaded = true
		postcard.Name = filename
		postcard.Path = path + filename
	}

	suffix := getSuffix(postcard.Href)
	if len(suffix) > 10 { // TODO if .com/embed...
		return
	}
	filename += suffix

	out, err := os.Create(path + filename)
	defer func() {
		if err = out.Close(); err != nil {
			errorChan <- err
		}
	}()

	if err != nil {
		errorChan <- err
	}

	response, err := http.Get(postcard.Href)
	defer func() {
		if err = response.Body.Close(); err != nil {
			errorChan <- err
		}
	}()

	_, err = io.Copy(out, response.Body)
	postcard.Downloaded = true
	postcard.Name = filename
	postcard.Path = path + filename

	if err != nil {
		errorChan <- err
	}
}

func NightlyPostcardDownload(errorChan chan error) {
	postcardsDatabase := storage.NewDatabase()
	postcardsStorage := storage.NewPostcardsPostgresStorage(postcardsDatabase)
	if err := postcardsStorage.ClearDatabase(); err != nil {
		errorChan <- err
	}

	newParser := parser.NewParser(os.Getenv("POSTCARD_SITE"))

	if err := newParser.GetHTML(); err != nil {
		errorChan <- err
	}

	holidays, err := newParser.GetHolidays()
	if err != nil {
		errorChan <- err
	}

	postcards, err := newParser.GetPostcardsPages(holidays[0])

	if err != nil {
		errorChan <- err
	}

	// Updating postcards ( adding links )
	for i := range postcards {
		if err = newParser.GetPostcardHref(&postcards[i]); err != nil {
			errorChan <- err
		}
	}
	postcardsPath := os.Getenv("POSTCARD_PATH")
	if err = os.RemoveAll(postcardsPath); err != nil {
		errorChan <- err
	}
	if err = os.Mkdir(postcardsPath, 0777); err != nil {
		errorChan <- err
	}

	for i := range postcards {
		PostcardDownload(postcardsPath, &postcards[i], errorChan)
	}

	for i := range postcards {
		if err = postcardsStorage.AddPostcardToStorage(&postcards[i]); err != nil {
			errorChan <- err
		}
	}
}
