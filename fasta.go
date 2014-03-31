package fasta

import (
	"bufio"
	"io/ioutil"
	"strings"
)

type Dna struct {
	Sequence []byte
}

type FastaFile struct {
	DnaSeqs map[string]Dna
}

func NewFastaFileWithStream(stream []byte) *FastaFile {
	var (
		text       string
		dna        []byte
		identifier string
		dnaMap     = make(map[string]Dna)
	)

	stringStream := string(stream)

	for _, line := range strings.Split(stringStream, "\n") {
		if len(line) == 0 {
			continue
		}

		text = strings.TrimSpace(line)

		// Hit an identifier line
		if text[0] == '>' {
			// If we stored a previous identifier, get the DNA string and map to the
			// identifier and clear the string
			if identifier != "" {
				dnaMap[identifier] = Dna{dna}
				dna = make([]byte, 0)
				identifier = ""
			}

			// Standard FASTA identifiers look like: ">foo_<id>"
			identifier = strings.Split(text, ">")[1]
		} else {
			// Append here since multi-line DNA strings are possible
			dna = append(dna, []byte(text)...)
		}
	}

	// EOF, there's one last identifier to store
	dnaMap[identifier] = Dna{dna}

	f := FastaFile{DnaSeqs: dnaMap}
	return &f
}

func NewFastaFileWithPath(filename string) *FastaFile {
	contents, err := ioutil.ReadFile(filename)

	if err != nil {
		panic("Not a valid file.")
	}

	return NewFastaFileWithStream(contents)
}
