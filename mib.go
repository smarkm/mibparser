package mibparser

import (
	"log"
	"regexp"
	"strings"
)

type Mib struct {
	Name    string
	SMIVer  string
	Valuels []MibValue
}

type MibValue struct {
	Name      string
	Number    string
	Parent    string
	Path      string
	SubValues []MibValue
}

type ObjectIdentifire struct{}
type ObjectType struct{}
type Definition struct{}

type Block struct {
	Lines        []string
	IsDefinition bool
	IsSequence   bool
	IsValue      bool
}

func (b *Block) Recognize() (v interface{}) {
	header := b.Lines[0]
	footer := b.Lines[len(b.Lines)-1]
	regParent, _ := regexp.Compile("[a-zA-Z][a-zA-Z0-9]+")
	regNubmer, _ := regexp.Compile("[0-9]+")
	switch {
	case strings.Contains(header, "DEFINITIONS ::= BEGIN"):
		b.IsDefinition = true
		v = Mib{Name: header[:strings.Index(header, " ")]}
	case strings.Contains(header, "::= SEQUENCE {"):
		b.IsSequence = true
	case strings.Contains(header, "OBJECT-TYPE") || strings.Contains(header, "OBJECT IDENTIFIER") ||
		strings.Contains(header, "OBJECT-GROUP"):
		b.IsValue = true
		parent := regParent.Find([]byte(footer))
		numbers := regNubmer.FindAllString(footer, -1)
		number := numbers[len(numbers)-1]
		v = MibValue{Name: header[:strings.Index(header, " ")], Parent: string(parent), Number: number}
		log.Println(string(parent))
	}
	//
	return v
}
