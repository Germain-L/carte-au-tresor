package fileutil

import (
	"bufio"
	"fmt"
	"os"
)

func WriteLinesToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write line to file: %w", err)
		}
	}

	return nil
}

func WriteToChannel(filename string, lines <-chan string) <-chan error {
	errors := make(chan error, 1)

	go func() {
		defer close(errors)

		file, err := os.Create(filename)
		if err != nil {
			errors <- fmt.Errorf("failed to create file %s: %w", filename, err)
			return
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		defer writer.Flush()

		for line := range lines {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				errors <- fmt.Errorf("failed to write line to file: %w", err)
				return
			}
		}
	}()

	return errors
}
