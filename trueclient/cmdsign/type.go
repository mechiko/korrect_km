package cmdsign

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aglyzov/charmap"
)

const modError = "cmdexec"

var paths = []string{
	`C:\Program Files\Crypto Pro\CSP`,
	`C:\Program Files (x86)\Crypto Pro\CSP`,
}

// var FileIn, FileOut string

type cmdexec struct {
	// *exec.Cmd
	FileIn  string
	FileOut string
	Hash    string
	Command []string
}

func setupEnvironment() error {
	path := os.Getenv("PATH")
	path = strings.Join(paths, ";") + ";" + path
	return os.Setenv("PATH", path)
}

func createTempFiles() (string, string, error) {
	fileIn, err := os.CreateTemp("", "lite_")
	if err != nil {
		return "", "", fmt.Errorf("error create IN temp file: %w", err)
	}
	defer fileIn.Close()

	fileOut, err := os.CreateTemp("", "lite_")
	if err != nil {
		return "", "", fmt.Errorf("error create OUT temp file: %w", err)
	}
	defer fileOut.Close()

	return fileIn.Name(), fileOut.Name(), nil
}

func New(hash string) *cmdexec {
	if err := setupEnvironment(); err != nil {
		panic(err.Error())
	}
	fileIn, fileOut, err := createTempFiles()
	if err != nil {
		panic(err.Error())
	}
	command := []string{
		"-sfsign", "-sign",
		"-in", fileIn,
		"-out", fileOut,
		"-my", hash,
		"-base64", "-add", "-addsigtime", "-cades_strict", "-verify",
		"-alg", "GOST12_256",
	}
	return &cmdexec{
		FileIn:  fileIn,
		FileOut: fileOut,
		Command: command,
		Hash:    hash,
	}
}

func (c *cmdexec) Sign(in string) (string, error) {
	err := os.WriteFile(c.FileIn, []byte(in), os.ModePerm)
	if err != nil {
		// panic(fmt.Sprintf("error write IN temp file %s", err.Error()))
		return "", fmt.Errorf("%s error write IN temp file: %w", modError, err)
	}

	cmd := exec.Command("csptest.exe", c.Command...)
	if out, err := cmd.CombinedOutput(); err != nil {
		decodeString := charmap.ToUTF8(&charmap.CP866_UTF8_TABLE, out)
		// fmt.Printf("sign %s", decodeString)
		return "", fmt.Errorf("%s %s Sign() %w", modError, decodeString, err)
	}
	// if err := cmd.Run(); err != nil {
	// 	return "", fmt.Errorf("%s Sign() %w", modError, err)
	// }
	sign, err := os.ReadFile(c.FileOut)
	if err != nil {
		return "", fmt.Errorf("%s Sign() %w", modError, err)
	}
	clearSign := bytes.ReplaceAll(sign, []byte("\r\n"), []byte{})
	// c.FilesRemove()
	return string(clearSign), nil
}

func (c *cmdexec) FilesRemove() {
	os.Remove(c.FileIn)
	os.Remove(c.FileOut)
}
