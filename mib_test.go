package mibparser_test

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/smarkm/mibparser"
)

var valuesMap map[string]*mibparser.MibValue

func TestParser(t *testing.T) {
	valuesMap = make(map[string]*mibparser.MibValue)
	println("===============")
	filename := "mibs/ietf/SNMPv2-MIB"
	fileHandle, _ := os.Open(filename)
	defer fileHandle.Close()
	var content []string
	fileScanner := bufio.NewScanner(fileHandle)
	var block mibparser.Block
	sequenceBlock := false
	block = mibparser.Block{}
	for fileScanner.Scan() {
		line := fileScanner.Text()

		if strings.HasPrefix(line, "--") || len(line) < 1 {
			continue
		}

		block.Lines = append(block.Lines, line)
		if strings.Contains(line, "::= SEQUENCE {") {
			sequenceBlock = true
		}
		if sequenceBlock {
			if strings.Contains(line, "}") {
				printlines(block)
				block = mibparser.Block{}
				sequenceBlock = false
			}
			continue
		}
		if strings.Contains(line, "::=") {
			printlines(block)
			block = mibparser.Block{}
		}
		content = append(content, line)
	}

}

func printlines(b mibparser.Block) {
	for _, line := range b.Lines {
		log.Printf(line)
	}
	obj := b.Recognize()
	if mibv, ok := obj.(mibparser.MibValue); ok {
		valuesMap[mibv.Name] = &mibv
		if parent, ok := valuesMap[mibv.Parent]; ok {
			parent.SubValues = append(parent.SubValues, mibv)
			log.Println("--------------", parent)
		}
	}
	log.Printf("===============IsDefinition: , IsSequence:,IsValue: \n", b.IsDefinition, b.IsSequence, b.IsValue, obj, len(valuesMap))
}

func TestReg(t *testing.T) {
	regNubmer, _ := regexp.Compile("[0-9]+")
	number := regNubmer.FindAllString("{Hell2 89 }", -1)
	log.Println("====", number[len(number)-1])
}
