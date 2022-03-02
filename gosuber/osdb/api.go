package osdb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	parser "github.com/lapuglisi/gosuber/parser"
)

// OsdbApiUri URI
// OsdbUserAgent user-agent
/*
-- Name Patterns as defined here --
!! Files are saved on source movie's directory
%M ..... Movie file name without extension
%L ..... LangId (eng, pob)
%X ..... Subtitle Extension (SubFormat)
%S ..... Subtitle file name
... more to come
*/
const (
	OsdbAPIURI         = "https://rest.opensubtitles.org/search"
	OsdbUserAgent      = "gosuber"
	DefaultNamePattern = "%M%L%X"
)

const (
	// ByHash is nice
	ByHash int = iota
	// ByString is nice too
	ByString
)

// Downloader asdasd
type Downloader struct {
	Languages []string
	MovieHash uint64
	MovieSize int64
	MovieFile string
	SaveAs    string
	Mode      int
}

func (dl *Downloader) generateFileName(sub parser.SubtitleInfo) (string, error) {
	var result string

	var regex *regexp.Regexp

	regex = regexp.MustCompile("%M|%L|%X|%S")
	if !regex.Match([]byte(dl.SaveAs)) {
		return "", fmt.Errorf("error: invalid name pattern '%s'", dl.SaveAs)
	}

	filePath := filepath.Dir(dl.MovieFile)
	fileName := filepath.Base(dl.MovieFile)
	fileExt := filepath.Ext(dl.MovieFile)

	// Strip extension from fileName
	if extPos := strings.LastIndex(fileName, fileExt); extPos > 0 {
		fileName = fileName[:extPos]
	}

	result = fmt.Sprintf("%s/%s", filePath, dl.SaveAs)
	result = strings.ReplaceAll(result, "%M", fileName)
	result = strings.ReplaceAll(result, "%L", sub.SubLanguageID)
	result = strings.ReplaceAll(result, "%X", sub.SubFormat)
	result = strings.ReplaceAll(result, "%S", sub.SubFileName)

	return result, nil
}

// GetSubtitles neathers
func (dl *Downloader) GetSubtitles() (subtitles parser.SubtitleInfoList, err error) {

	var requestURL string
	var foundSubs parser.SubtitleInfoList
	isAnEpisode := regexp.MustCompile("(.+)\\s*[Ss](\\d{2})[Ee](\\d{2})")

	if dl.MovieHash == 0 || dl.MovieSize == 0 {
		return nil, fmt.Errorf("either hash or size is not valid")
	}

	httpClient := http.Client{
		Timeout: time.Second * 30, // 30 seconds
	}

	for _, lang := range dl.Languages {

		switch dl.Mode {
		case ByString:
			fileName := filepath.Base(dl.MovieFile)

			// Check whether we are dealing with an episode from
			// any series
			matches := isAnEpisode.FindAllStringSubmatch(fileName, -1)
			if matches != nil && len(matches) > 3 {
				// Ok, this is an episode
				// [0] is entire string
				// [1-...] submatches
				requestURL = fmt.Sprintf("%s/query-%s/episode-%s/season-%s/sublanguageid-%s/tags-web-dl",
					OsdbAPIURI,
					matches[1], matches[2], matches[3],
					lang)
			} else {
				requestURL = fmt.Sprintf("%s/query-%s/sublanguageid-%s",
					OsdbAPIURI, url.PathEscape(fileName), lang)
			}

			// Give our httpClient more time to think
			httpClient.Timeout = time.Second * 60
			break

		case ByHash:
			fallthrough
		default:
			// Just return
			// Create request URL
			requestURL = fmt.Sprintf("%s/moviebytesize-%d/moviehash-%x/sublanguageid-%s",
				OsdbAPIURI, dl.MovieSize, dl.MovieHash, lang)
		}

		// Create new http Request
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", OsdbUserAgent)

		// Launch http request
		res, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if res.StatusCode != 200 {
			return nil, fmt.Errorf("error: %s", res.Status)
		}

		// Read the response
		jsonData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// Unmarshal our result
		foundSubs = make(parser.SubtitleInfoList, 0)

		err = json.Unmarshal(jsonData, &foundSubs)

		if len(foundSubs) > 0 {
			subtitles = append(subtitles, foundSubs...)
		}
	}

	return subtitles, nil
}

// Download FTW .................
func (dl *Downloader) Download(sub parser.SubtitleInfo) error {

	var fileDir string
	var fileName string
	var target string

	fileDir = filepath.Dir(dl.MovieFile)
	if len(fileDir) == 0 {
		fileDir = "."
	}

	fileName = filepath.Base(dl.MovieFile)
	if len(fileName) == 0 {
		fileName = sub.SubFileName
	}

	if len(dl.SaveAs) == 0 {
		dl.SaveAs = DefaultNamePattern
	}

	// Create target file name based on pattern
	target, err := dl.generateFileName(sub)
	if err != nil {
		return err
	}

	fmt.Printf("\033[33mDownloading '%s' as '%s'... \033[0m",
		sub.SubFileName, target)

	_, err = os.Stat(target)
	if err == nil {
		// File exists
		return fmt.Errorf("error: file '%s' exists", target)
	}

	targetSub, err := os.Create(target)
	if err != nil {
		return err
	}
	defer targetSub.Close()

	// Start http request
	response, err := http.Get(sub.SubDownloadLink)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(targetSub, response.Body)

	if err == nil {
		fmt.Printf("\033[32mOK!\033[0m\n")
	} else {
		fmt.Printf("\033[31mERROR: %s\033[0m\n", err.Error())
	}

	return err
}
