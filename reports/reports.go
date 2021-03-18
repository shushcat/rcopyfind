package reports

import (
	"fmt"
	"github.com/shushcat/rlat/comparator"
	// "os"
	"sort"
	"strings"
)

type Comparator = comparator.Comparator

func limitContext(wordArray []string, indexSlice []int) []string {
	// TODO Allow setting context length from the command line.
	contextWords := 10
	sort.Ints(indexSlice)
	// firstIndex := indexSlice[0]
	// lastIndex := indexSlice[len(indexSlice)-1]
	// beg, end := firstIndex - contextWords, lastIndex - contextWords
	beg := indexSlice[0] - contextWords
	if beg < 0 {
		beg = 0
	}
	end := indexSlice[len(indexSlice)-1] + contextWords
	if end > len(wordArray)-1 {
		end = len(wordArray) - 1
	}
	var limitedWordArray []string
	for i := beg; i <= end; i++ {
		limitedWordArray = append(limitedWordArray, wordArray[i])
	}
	return limitedWordArray
}

func highlight(wordArray []string, indexSlice []int) []string {
	// Duplicate wordArray.
	highlitArray := wordArray
	// TODO Highlight rather than upcasing.  Implementing this may be
	// impractical at this level since imprementing this prior to word
	// wrapping will require making that function significantly more
	// complicated.  Fuck the mail.
	// var pre, post string
	// Color matching words red if the program output is being sent to a terminal, and surround matching words with asterisks otherwise.
	// if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
	// red := "\033[31m"
	// black := "\033[0m"
	// pre = red
	// post = black
	// } else {
	// pre = "**"
	// post = "**"
	// }
	for _, i := range indexSlice {
		// highlitArray[i] = (pre + highlitArray[i] + post)
		highlitArray[i] = strings.ToUpper(highlitArray[i])
	}
	return highlitArray
}

func PrintReport(cmp Comparator) {
	fmt.Println("--------------------------------")
	if len(cmp.SimilarClusters) == 0 {
		fmt.Println("No similar clusters found.")
		fmt.Println("--------------------------------")
	} else {
		highlit := make([]string, len(cmp.Target.WordArray))
		for _, pair := range cmp.SimilarClusters {
			copy(highlit, cmp.Target.WordArray)
			highlit = highlight(highlit, pair[0].Values())
			targetLine := strings.Join(limitContext(highlit, pair[0].Values()), " ")
			fmt.Println(wrapLine(targetLine, "> "))
			fmt.Println("---")
			copy(highlit, cmp.Source.WordArray)
			highlit = highlight(highlit, pair[1].Values())
			sourceLine := strings.Join(limitContext(highlit, pair[1].Values()), " ")
			fmt.Println(wrapLine(sourceLine, "< "))
			fmt.Println("--------------------------------")
		}
	}
}

func wrapLine(line string, pre string) string {
	cols := 80
	wrapped := ""
	var i int
	for len(line[i:]) > cols {
		if string(line[i+cols]) == " " {
			wrapped += pre + strings.TrimSpace(line[i:i+cols]) + "\n"
			i = i + cols
		} else {
			for j := (i + cols) - 1; j > i; j-- {
				if string(line[j]) == " " {
					wrapped += pre + strings.TrimSpace(line[i:j]) + "\n"
					i = j
					break
				}
			}
		}
	}
	wrapped += pre + strings.TrimSpace(line[i:])
	return wrapped
}
