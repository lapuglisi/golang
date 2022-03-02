package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lapuglisi/gosuber/osdb"
	"github.com/lapuglisi/gosuber/parser"
)

func usage(why string) {
	fmt.Println()
	fmt.Printf("\033[33moops: %s\n\033[0m", why)
	fmt.Printf("\n")
	fmt.Println("gosuber ARGS")
	flag.PrintDefaults()
}

func downloadPattern(dl osdb.Downloader, subs parser.SubtitleInfoList, pattern string) {
	// dlMode is a pattern to search within subtitles
	saveAs := dl.SaveAs
	index := 0
	for _, sub := range subs {
		if ok, err := sub.Matches(pattern); ok {
			fmt.Printf("\033[32mFound subtitle for pattern '%s'. Downloading it.\033[0m\n", pattern)

			dl.SaveAs = fmt.Sprintf("%.3d_%s", index, saveAs)
			dl.Download(sub)

			index++
		} else {
			if err != nil {
				fmt.Printf("\033[31merror: %s\033[0m", err.Error())
			}
		}
	}
}

func downloadInteractive(subs parser.SubtitleInfoList, dl osdb.Downloader) {
	var subIndex int
	var totalSubs int = len(subs) - 1
	var err error

	fmt.Printf("\033[32mSubtitles found!\033[0m\n")
	for index, sub := range subs {
		fmt.Printf(
			sub.InfoString(),
			index, sub.MovieName, sub.LanguageName,
			sub.Score, sub.SubFileName)
		fmt.Println()
	}

	for {
		fmt.Println()
		fmt.Printf("\033[33mChoose which one to download (-1 to quit) [0-%d]: \033[0m", totalSubs)
		_, err = fmt.Scanf("%d\n", &subIndex)

		if subIndex == -1 {
			fmt.Printf("OK. Goodbye!\n")
			break
		}

		if subIndex < 0 || subIndex > totalSubs || err != nil {
			fmt.Printf("\033[31mWTF? (%s)\033[0m\n", err.Error())

			var discard byte
			fmt.Scanln(&discard)
		} else {
			chosen := subs[subIndex]
			fmt.Printf("\033[32mYou chose %s! Congratulations!\033[0m\n", chosen.SubFileName)

			dl.Download(chosen)
			break
		}
	}
}

func main() {
	// Init flags
	var filePath string
	var subLangs string
	var dlMode string
	var namePattern string
	var dl osdb.Downloader

	flag.StringVar(&filePath, "movie", "", "Movie file to search for subtitles")
	flag.StringVar(&subLangs, "lang", "pob", "Subtitle language: pob,eng")
	flag.StringVar(&dlMode, "download", "top-rated", "top-rated, interactive or \"pattern\"")
	flag.StringVar(&namePattern, "saveas", "%M.%L.%X", "file name format:\n%M - Movie file name\n%S - Subtitle file name\n%X - Subtitle extension\n%L - Language id (eng,pob)\n")

	flag.Parse()
	if len(filePath) == 0 {
		usage("movie name not defined")
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error: Could not open '%s': %s\n", filePath, err.Error())
	}
	defer file.Close()

	hash, err := osdb.HashFile(file)

	if err != nil {
		fmt.Printf("error: HashFile failed: %s\n", err.Error())
	}

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("\033[31merror: %s\033[0m\n", err.Error())
		os.Exit(-1)
	}

	// Fill in information
	dl.Mode = osdb.ByHash
	dl.MovieHash = hash
	dl.MovieSize = fileInfo.Size()
	dl.Languages = strings.Split(subLangs, ",")
	dl.MovieFile = filePath
	dl.SaveAs = namePattern

	fmt.Printf("\033[32mSearching subtitles for:\033[0m\n")
	fmt.Printf("  movie: %s\n", dl.MovieFile)
	fmt.Printf("   hash: %x\n", dl.MovieHash)
	fmt.Printf("   size: %d\n", dl.MovieSize)
	fmt.Printf("  langs: %s\n", strings.Join(dl.Languages, ""))
	fmt.Printf("save as: %s (pattern)\n\n", dl.SaveAs)

	subtitles, err := dl.GetSubtitles()
	if err != nil {
		fmt.Printf("\033[31m%s\n\033[0m", err.Error())
		os.Exit(-1)
	} else if len(subtitles) == 0 {
		fmt.Printf("\033[33mNo subtitles found.\033[0m\n")
		fmt.Printf("\033[33mTrying to search by file name....\033[0m\n")

		dl.Mode = osdb.ByString
		subtitles, err = dl.GetSubtitles()
	}

	// After last attempt
	if err != nil {
		fmt.Printf("\033[31m%s\n\033[0m", err.Error())
		os.Exit(-1)
	} else if len(subtitles) == 0 {
		fmt.Printf("\033[33mNo subtitles found.\033[0m\n")
		os.Exit(0)
	}

	switch dlMode {
	case "top-rated":
		dl.Download(subtitles[0])
		break

	case "interactive":
		downloadInteractive(subtitles, dl)

	default:
		downloadPattern(dl, subtitles, dlMode)
		break
	}

}
