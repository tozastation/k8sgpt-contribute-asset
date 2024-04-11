package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
	"golang.org/x/exp/slog"
)

func main() {
	llm, err := ollama.New(ollama.WithModel("llama2"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	qdrantURL, err := url.Parse("http://localhost:6333")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	store, err := qdrant.New(
		qdrant.WithURL(*qdrantURL),
		qdrant.WithCollectionName("kubernetes_changelog"),
		qdrant.WithEmbedder(e),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = 300   // size of the chunk is number of characters
	split.ChunkOverlap = 30 // overlap is the number of characters that the chunks overlap
	err = filepath.Walk("./files", func(path string, info os.FileInfo, err error) error {
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
				slog.Error(err.Error())
				os.Exit(1)
			}
			docs, err := documentloaders.NewText(dat).LoadAndSplit(context.Background(), split)
			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
			fmt.Println(docs)
			_, err = store.AddDocuments(context.Background(), docs)
			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
		}
		return nil
	})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
