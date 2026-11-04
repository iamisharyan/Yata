package preview

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
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
func Highlightfile(path string)(string ,error){
	c,err:=os.ReadFile(path)
	if err!=nil{
		return "",err
	}
	for _,b:=range c{
		if b==0{
			return "",fmt.Errorf("binary file cant be highlighted")
		}
	}
	lexer:=lexers.Match(path)
	if lexer == nil{
		lexer=lexers.Fallback
	}
	it,err:=lexer.Tokenise(nil,string(c))

	if err!=nil{
		return "",err
	}

	formatter:=formatters.Get("terminal256")
	if formatter==nil{
		return "",fmt.Errorf("terminal formatter unavailable")

	}
	style :=styles.Get("monkai")

	var output bytes.Buffer
	err=formatter.Format(&output,style,it)
	if err!=nil{
		return "",err
	}
	return output.String(),nil
}
