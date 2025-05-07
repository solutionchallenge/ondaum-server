package llm

type Config struct {
	Gemini GenericConfig `mapstructure:"gemini"`
}

type PromptType string

const (
	PromptTypeSystemInstruction PromptType = "system_instruction"
	PromptTypeActionPrompt      PromptType = "action_prompt"
)

type PreparedPrompt struct {
	Identifier     string     `mapstructure:"identifier"`
	PromptType     PromptType `mapstructure:"prompt_type"`
	PromptFile     string     `mapstructure:"prompt_file"`
	AttachmentFile string     `mapstructure:"attachment_file"`
	AttachmentMime string     `mapstructure:"attachment_mime"`
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
	PreparedPrompts    []PreparedPrompt   `mapstructure:"prepared_prompts"`
	RedactionThreshold RedactionThreshold `mapstructure:"redaction_threshold"`
}
