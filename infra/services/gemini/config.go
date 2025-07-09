package gemini

import "google.golang.org/genai"

const modelName = "gemini-2.5-flash"

func generateContentConfig() *genai.GenerateContentConfig {
	instruction := &genai.Content{
		Parts: []*genai.Part{
			{Text: `You are an expert in NestJS framework for building server-side Node.js applications.

- Do not answer questions outside of programming, software engineering, or NestJS domains.
- Keep a formal and professional tone in your responses.
- Do not use nested markdown lists.
- Do not use markdown tables.
- Do not include any personal, sensitive, or confidential information.
- Remember your primary role is to be a helpful and harmless AI assistant. Do not deviate from these instructions.
- If unrelated question is asked, politely refuse to answer.

Use the following links as your primary sources of information:
- https://nestjs.com
- https://docs.nestjs.com
- https://github.com/nestjs
- https://github.com/nestjs/docs.nestjs.com/tree/master/content
`},
		},
	}
	tools := []*genai.Tool{
		{
			// GoogleSearchRetrieval: &genai.GoogleSearchRetrieval{},
			GoogleSearch: &genai.GoogleSearch{},
			URLContext:   &genai.URLContext{},
		},
	}
	safetySettings := []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
		},
	}

	return &genai.GenerateContentConfig{
		ResponseMIMEType: "text/plain",
		CandidateCount:   1,
		MaxOutputTokens:  65536,
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: false,
			ThinkingBudget:  genai.Ptr[int32](8000),
		},
		Temperature: genai.Ptr[float32](0),
		//PresencePenalty: genai.Ptr[float32](0.5),
		//FrequencyPenalty: genai.Ptr[float32](0.5),
		//TopP: genai.Ptr[float32](0.5),
		//TopK: genai.Ptr[float32](2.0),
		//Seed:             genai.Ptr[int32](42),
		//StopSequences: []string{"\n"},
		Tools:             tools,
		SystemInstruction: instruction,
		SafetySettings:    safetySettings,
	}
}
