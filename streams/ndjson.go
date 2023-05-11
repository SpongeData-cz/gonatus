package streams

import (
	"bufio"
	"os"

	"github.com/SpongeData-cz/gonatus"
)

const (
	FileWrite = iota
	FileAppend
)

type NdjsonInputStreamer interface {
	InputStreamer[gonatus.Conf]
}

type NdjsonOutputStreamer interface {
	OutputStreamer[gonatus.Conf]
}

type ndjsonInputStream struct {
	inputStream[gonatus.Conf]
	file    *os.File
	scanner *bufio.Scanner
}

func NewNdjsonInputStream(path string) NdjsonInputStreamer {

	ego := &ndjsonInputStream{}
	ego.init(ego)

	file, err := os.Open(path)
	check(err)
	ego.file = file
	ego.scanner = bufio.NewScanner(file)

	if !ego.scanner.Scan() {
		ego.file.Close()
		ego.closed = true
		panic("File is empty.")
	}

	return ego

}

func (ego *ndjsonInputStream) get() (value gonatus.Conf, valid bool, err error) {

	if ego.file == nil {
		panic("The file does not exist.")
	}

	newConf := gonatus.NewConf("")
	newConf.Unmarshal([]byte(ego.scanner.Text()))

	if !ego.scanner.Scan() {
		ego.file.Close()
		ego.closed = true
	}

	return newConf, true, nil

}

type ndjsonOutputStream struct {
	outputStream[gonatus.Conf]
	file *os.File
}

func NewNdjsonOutputStream(path string, mode int) NdjsonOutputStreamer {

	if mode != FileAppend && mode != FileWrite {
		panic("Unknown mode.")
	}

	ego := &ndjsonOutputStream{}
	ego.init(ego)

	var flags int
	if mode == FileWrite {
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	} else {
		flags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	}

	file, err := os.OpenFile(path, flags, 0664)
	check(err)

	ego.file = file

	return ego

}

func (ego *ndjsonOutputStream) setSource(s InputStreamer[gonatus.Conf]) {
	ego.source = s
	ego.export()
}

func (ego *ndjsonOutputStream) export() {

	for true {

		value, valid, err := ego.source.get()
		if !valid {
			break
		}

		check(err)
		nd, err := value.Marshal()
		check(err)
		_, err = ego.file.Write(nd)
		check(err)
		_, err = ego.file.WriteString("\n")
		check(err)

		if ego.source.Closed() {
			break
		}

	}

	ego.closed = true
	ego.file.Close()

}
