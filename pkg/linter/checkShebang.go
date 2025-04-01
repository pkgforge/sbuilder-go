package linter

import (
	"bufio"
	"bytes"
	"strings"
)

func (v *Validator) checkShebang() ([]byte, error, string) {
	scanner := bufio.NewScanner(bytes.NewReader(v.raw))
	var buffer bytes.Buffer

	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
		buffer.WriteByte('\n')
	}

	if buffer.Len() > 0 && strings.HasPrefix(buffer.String(), "#!/SBUILD") {
		return nil, nil, ""
	}

	newContent := bytes.NewBufferString("#!/SBUILD\n")
	newContent.Write(buffer.Bytes())

	v.raw = newContent.Bytes()

	return newContent.Bytes(), nil, "didn't find a shebang in the SBUILD file, automatically fixed"
}
