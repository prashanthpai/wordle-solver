package wordle

import (
	"bufio"
	"embed"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

const wordsFile = "dict_wikipedia_freq.txt"

//go:embed dict_wikipedia_freq.txt
var efs embed.FS

func openFile() (fs.File, error) {
	return os.Open(wordsFile)
}

func openEmbed() (fs.File, error) {
	return efs.Open(wordsFile)
}

func loadDict() (map[string]int, error) {
	f, err := openEmbed()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dict := make(map[string]int)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		processLine(dict, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dict, nil
}

func processLine(dict map[string]int, line string) {

	word := strings.ToUpper(strings.TrimSpace(line))
	l := strings.Split(word, " ")

	if len(l) == 1 {
		dict[word] = 0
	}

	if len(l) == 2 {
		i, err := strconv.Atoi(l[1])
		if err != nil {
			return
		}
		dict[l[0]] = i
	}
}
