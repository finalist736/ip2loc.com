package config

import (
	"bufio"
	"os"
	"strings"
)

// FileProvider describes a file based loader which loads the configuration
// from a file listed.
type FileProvider struct {
	Filename string
}

func NewFileProvider(name *string) *FileProvider {
	return &FileProvider{Filename: *name}
}

// Provide implements the Provider interface.
func (fp FileProvider) Provide() (map[string]string, error) {
	var configMap = make(map[string]string)

	file, err := os.Open(fp.Filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) < 3 {
			// the line doesn't have enough data
			continue
		}

		if line[0] == '#' || line[0] == ';'  {
			// the line starts with a comment character
			continue
		}

		// find the first equals sign
		index := strings.Index(line, "=")

		// if we couldn't find one
		if index <= 0 {
			// the line is invalid
			continue
		}

		if index == len(line) {
			// the line is invalid
			continue
		}

		// add the item to the config
		configMap[line[:index]] = line[index+1:]
	}

	return configMap, nil
}