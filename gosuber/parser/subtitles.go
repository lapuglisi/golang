package parser

import (
	"regexp"
	"strings"
	"time"
)

// OsdbDate custom date for osdb
type OsdbDate struct {
	time.Time
}

// UnmarshalJSON unmarshaller
func (dt *OsdbDate) UnmarshalJSON(input []byte) error {
	strDate := string(input)
	strDate = strings.Trim(strDate, `"`)

	osdbDate, err := time.Parse("2006-01-02 15:04:05", strDate)
	if err != nil {
		return err
	}

	dt.Time = osdbDate

	return nil
}

// QueryParametersTag ftw
type QueryParametersTag struct {
	MovieHash     string `json:"moviehash"`
	MovieByteSize int64  `json:"moviebytesize"`
}

// SubtitleInfo ...
type SubtitleInfo struct {
	MatchedBy           string             `json:"MatchedBy"`
	IDSubMovieFile      string             `json:"IDSubMovieFile"`
	MovieHash           string             `json:"MovieHash"`
	MovieByteSize       uint64             `json:"MovieByteSize"`
	MovieTimeMS         int                `json:"MovieTimeMS"`
	IDSubtitleFile      string             `json:"IDSubtitleFile"`
	SubFileName         string             `json:"SubFileName"`
	SubActualCD         int                `json:"SubActualCD"`
	SubSize             int                `json:"SubSize"`
	SubHash             string             `json:"SubHash"`
	SubLastTS           string             `json:"SubLastTS"`
	SubTSGroup          int                `json:"SubTSGroup"`
	InfoReleaseGroup    string             `json:"InfoReleaseGroup"`
	InfoFormat          string             `json:"InfoFormat"`
	InfoOther           string             `json:"InfoOther"`
	IDSubtitle          int                `json:"IDSubtitle"`
	UserID              int                `json:"UserID"`
	SubLanguageID       string             `json:"SubLanguageID"`
	SubFormat           string             `json:"SubFormat"`
	SubSumCD            int                `json:"SubSumCD"`
	SubAuthorComment    string             `json:"SubAuthorComment"`
	SubAddDate          OsdbDate           `json:"SubAddDate"`
	SubBad              int                `json:"SubBad"`
	SubRating           float32            `json:"SubRating"`
	SubSumVotes         int                `json:"SubSumVotes"`
	SubDownloadsCnt     int                `json:"SubDownloadsCnt"`
	MovieReleaseName    string             `json:"MovieReleaseName"`
	MovieFPS            float32            `json:"MovieFPS"`
	IDMovie             string             `json:"IDMovie"`
	IDMovieImdb         string             `json:"IDMovieImdb"`
	MovieName           string             `json:"MovieName"`
	MovieNameEng        string             `json:"MovieNameEng"`
	MovieYear           int                `json:"MovieYear"`
	MovieImdbRating     float32            `json:"MovieImdbRating"`
	SubFeatured         int                `json:"SubFeatured"`
	UserNickName        string             `json:"UserNickName"`
	SubTranslator       string             `json:"SubTranslator"`
	ISO639              string             `json:"ISO639"`
	LanguageName        string             `json:"LanguageName"`
	SubComments         int                `json:"SubComments"`
	SubHearingImpaired  int                `json:"SubHearingImpaired"`
	UserRank            string             `json:"UserRank"`
	SeriesSeason        int                `json:"SeriesSeason"`
	SeriesEpisode       int                `json:"SeriesEpisode"`
	MovieKind           string             `json:"MovieKind"`
	SubHD               bool               `json:"SubHD"`
	SeriesIMDBParent    string             `json:"SeriesIMDBParent"`
	SubEncoding         string             `json:"SubEncoding"`
	SubAutoTranslation  bool               `json:"SubAutoTranslation"`
	SubForeignPartsOnly bool               `json:"SubForeignPartsOnly"`
	SubFromTrusted      bool               `json:"SubFromTrusted"`
	SubTSGroupHash      string             `json:"SubTSGroupHash"`
	SubDownloadLink     string             `json:"SubDownloadLink"`
	ZipDownloadLink     string             `json:"ZipDownloadLink"`
	SubtitlesLink       string             `json:"SubtitlesLink"`
	QueryNumber         int                `json:"QueryNumber"`
	QueryParamenters    QueryParametersTag `json:"QueryParameters"`
	Score               float32            `json:"Score"`
}

func (sub *SubtitleInfo) InfoString() string {
	return "[%d] \033[33mMovie: \033[0m%s, \033[33mLanguage: \033[0m%s, \033[33mScore: \033[0m%f, \033[33mSubtitle: \033[0m%s"
}

func (sub *SubtitleInfo) Matches(pattern string) (bool, error) {
	regex, err := regexp.Compile(pattern)

	if err != nil {
		return false, err
	}

	matches := regex.Match([]byte(sub.SubFileName)) ||
		regex.Match([]byte(sub.InfoReleaseGroup)) ||
		regex.Match([]byte(sub.MovieReleaseName)) ||
		regex.Match([]byte(sub.SubFormat))

	return matches, nil
}

// SubtitleInfoList .....
type SubtitleInfoList []SubtitleInfo
