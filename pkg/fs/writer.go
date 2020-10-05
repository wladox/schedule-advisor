package fs

import (
	"bufio"
	"os"
	"strings"
)

type FileWriter struct {
}

// writeLines writes the lines to the given file.
func (w *FileWriter) Save(data []string, path string) error {
	s := strings.Split(path, "/")
	if len(s) > 0 {
		if _, err := os.Stat(s[0]); os.IsNotExist(err) {
			_ = os.Mkdir(s[0], 0766)
		}
	}

	if len(data) > 0 {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		defer writer.Flush()

		for _, stats := range data {
			_, err = writer.WriteString(stats + "\n")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
