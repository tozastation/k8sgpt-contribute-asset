package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

func main() {
	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = 300   // size of the chunk is number of characters
	split.ChunkOverlap = 30 // overlap is the number of characters that the chunks overlap
	err := filepath.Walk("./files", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("dir: %v name: %s\n", info.IsDir(), path)
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, ".md") {
			dat, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
			}
			docs, err := documentloaders.NewText(dat).LoadAndSplit(context.Background(), split)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(docs)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
