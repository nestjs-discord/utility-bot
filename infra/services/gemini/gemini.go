package gemini

import (
	"context"
	"fmt"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"google.golang.org/genai"
	"strings"
	"sync"
)

type Gemini struct {
	clients          []*genai.Client
	currentClientIdx int
	lock             sync.RWMutex
	genContentConfig *genai.GenerateContentConfig
}

func NewGemini(cfg *env.GeminiConfig) (*Gemini, error) {
	g := &Gemini{
		genContentConfig: generateContentConfig(),
	}

	err := g.generateClients(cfg)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Gemini) generateClients(cfg *env.GeminiConfig) error {
	ctx := context.Background()
	for _, apiKey := range cfg.APIKeys {
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  apiKey,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			return fmt.Errorf("failed to create gemini client: %w", err)
		}
		g.clients = append(g.clients, client)
	}
	return nil
}

type GenerateResponse struct {
	Content           string
	Sources           []Source
	SearchSuggestions []string
}

type Source struct {
	Text string
	URL  string
}

func (g *Gemini) getCurrentClient() *genai.Client {
	g.lock.RLock()
	defer g.lock.RUnlock()

	client := g.clients[g.currentClientIdx]
	g.currentClientIdx = (g.currentClientIdx + 1) % len(g.clients)
	return client
}

func (g *Gemini) Generate(ctx context.Context, prompts ...string) (*GenerateResponse, error) {
	client := g.getCurrentClient()
	contents := genai.Text(strings.Join(prompts, "\n"))
	result, err := client.Models.GenerateContent(ctx, modelName, contents, g.genContentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content with gemini client: %w", err)
	}

	if len(result.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates returned from gemini client")
	}

	candidate := result.Candidates[0]

	if candidate.FinishReason != genai.FinishReasonStop {
		return nil, fmt.Errorf("got an unexpected finish reason: %s", candidate.FinishReason)
	}

	builder := strings.Builder{}
	for _, part := range candidate.Content.Parts {
		if part.Text == "" {
			continue
		}

		builder.WriteString("\n")
		text := strings.TrimSpace(part.Text)
		//text = strings.ReplaceAll(text, "    ```", "```") // remove leading spaces from code blocks
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimPrefix(line, "    ")
			if line == "" {
				continue // skip empty lines
			}
			builder.WriteString(line)
			builder.WriteString("\n")
		}
	}
	res := &GenerateResponse{
		Content: strings.TrimSpace(builder.String()),
		Sources: g.mapGroundingMetadataToSources(candidate.GroundingMetadata),
	}

	if candidate.GroundingMetadata != nil {
		res.SearchSuggestions = candidate.GroundingMetadata.WebSearchQueries
	}

	return res, nil
}

func (g *Gemini) mapGroundingMetadataToSources(metadata *genai.GroundingMetadata) []Source {
	if metadata == nil || len(metadata.GroundingChunks) == 0 {
		return nil
	}

	var sources []Source
	for _, chunk := range metadata.GroundingChunks {
		if chunk.Web == nil || chunk.Web.Title == "" || chunk.Web.URI == "" {
			continue
		}
		sources = append(sources, Source{
			Text: chunk.Web.Title,
			URL:  chunk.Web.URI,
		})
	}
	return sources
}
