package preview

import (
	"bufio"
	"io"
	"os"
	"strings"
)
const MaxPreviewLines = 200


func TextPreview(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bin, err := isBinary(f)
	if err != nil {
		return "", err
	}
	if bin {
		return "[binary file]", nil
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	var b strings.Builder
	scanner := bufio.NewScanner(f)
	lines := 0

	for scanner.Scan() && lines < MaxPreviewLines {
		b.WriteString(scanner.Text() + "\n")
		lines++
	}

	return b.String(), scanner.Err()
}

func isBinary(f *os.File) (bool, error) {
	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return true, nil
		}
	}
	return false, nil
}
