package llm

type Config struct {
	Gemini GenericConfig `mapstructure:"gemini"`
}

type CachedPrompt struct {
	Identifier string `mapstructure:"identifier"`
	Prompt     string `mapstructure:"prompt"`
	Attachment string `mapstructure:"attachment"`
}

type PredefinedAction struct {
	Identifier string `mapstructure:"identifier"`
	Reference  string `mapstructure:"reference"`
}

type RedactionThreshold struct {
	Harrasement      string `mapstructure:"harrasement"`
	HateSpeech       string `mapstructure:"hate_speech"`
	SexuallyExplicit string `mapstructure:"sexually_explicit"`
	DangerousContent string `mapstructure:"dangerous_content"`
	CivicIntegrity   string `mapstructure:"civic_integrity"`
}

type GenericConfig struct {
	Enabled            bool               `mapstructure:"enabled"`
	APIKey             string             `mapstructure:"api_key"`
	LLMModel           string             `mapstructure:"llm_model"`
	EmbeddingModel     string             `mapstructure:"embedding_model"`
	ResponseFormat     string             `mapstructure:"response_format"`
	SystemPrompt       string             `mapstructure:"system_prompt"`
	CachedPrompts      []CachedPrompt     `mapstructure:"cached_prompts"`
	PredefinedActions  []PredefinedAction `mapstructure:"predefined_actions"`
	RedactionThreshold RedactionThreshold `mapstructure:"redaction_threshold"`
}
