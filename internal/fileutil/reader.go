package fileutil

import (
	"bufio"
	"os"
)

func ReadToChannel(fileName string) (<-chan string, <-chan error) {
	lines := make(chan string)
	errors := make(chan error, 1) // Buffered channel

	go func() {
		defer close(lines)
		defer close(errors)

		if _, err := os.Stat(fileName); err != nil {
			errors <- err
			return
		}

		file, err := os.Open(fileName)
		if err != nil {
			errors <- err
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lines <- line
		}

		if err := scanner.Err(); err != nil {
			errors <- err
			return
		}

		// Send nil to indicate no error
		errors <- nil
	}()

	return lines, errors
}
