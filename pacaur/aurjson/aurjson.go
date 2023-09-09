package aurjson

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	AURSearchURL      string = "https://aur.archlinux.org/rpc/?v=5&type=search&by=name-desc&arg=%s"
	AURPackageInfoURL string = "https://aur.archlinux.org/rpc/?v=5&type=info&arg[]=%s"
)

// / AURSearchResult is suave
type AURSearchResult struct {
	ID             int32    `json:"ID"`
	Name           string   `json:"Name"`
	PackageBaseID  int32    `json:"PackageBaseID"`
	PackageBase    string   `json:"PackageBase"`
	Version        string   `json:"Version"`
	Description    string   `json:"Description"`
	URL            string   `json:"URL"`
	NumVotes       int32    `json:"NumVotes"`
	Popularity     float32  `json:"Popularity"`
	OutOfDate      int32    `json:"OutOfDate"`
	Maintainer     string   `json:"Maintainer"`
	FirstSubmitted uint32   `json:"FirstSubmitted"`
	LastModified   uint32   `json:"LastModified"`
	URLPath        string   `json:"URLPath"`
	Depends        []string `json:"Depends"`
	MakeDepends    []string `json:"MakeDepends"`
	License        []string `json:"License"`
	Keywords       []string `json:"Keywords"`
}

// / AURJson is da hora
type AURJson struct {
	Version     int32             `json:"version"`
	Type        string            `json:"type"`
	ResultCount int32             `json:"resultcount"`
	Results     []AURSearchResult `json:"results"`
	Error       string            `json:"error"`
}

func (aurData *AURJson) FromJsonBytes(jsonBytes []byte) (err error) {
	err = nil

	err = json.Unmarshal(jsonBytes, aurData)

	return err
}

func SearchAUR(search string) (result *AURJson, err error) {

	var searchUrl string = fmt.Sprintf(AURSearchURL, search)

	resp, err := http.Get(searchUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result = new(AURJson)

	err = result.FromJsonBytes(respBytes)
	if err != nil {
		return nil, err
	}

	if result.Type == "error" {
		return nil, fmt.Errorf("JSON error: %s", result.Error)
	}

	return result, nil
}

func GetPackageInfo(aurPack string) (result *AURJson, err error) {

	var infoUrl string = fmt.Sprintf(AURPackageInfoURL, aurPack)

	resp, err := http.Get(infoUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result = new(AURJson)

	err = result.FromJsonBytes(respBytes)
	if err != nil {
		return nil, err
	}

	return result, nil
}
