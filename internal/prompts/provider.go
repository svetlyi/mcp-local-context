package prompts

type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
	Content     string           `json:"content"`
	Language    string           `json:"language,omitempty"`
}

type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type Provider interface {
	GetPrompts() []Prompt
}

type Registry struct {
	providers []Provider
}

func NewRegistry() *Registry {
	return &Registry{
		providers: make([]Provider, 0),
	}
}

func (r *Registry) Register(provider Provider) {
	r.providers = append(r.providers, provider)
}

func (r *Registry) GetAllPrompts() []Prompt {
	allPrompts := make([]Prompt, 0, 10)
	for _, provider := range r.providers {
		allPrompts = append(allPrompts, provider.GetPrompts()...)
	}
	return allPrompts
}

func (r *Registry) GetPrompt(name string) *Prompt {
	for _, prompt := range r.GetAllPrompts() {
		if prompt.Name == name {
			return &prompt
		}
	}
	return nil
}

func (r *Registry) GetSupportedLanguages() []string {
	languageMap := make(map[string]bool)
	for _, prompt := range r.GetAllPrompts() {
		if prompt.Language != "" {
			languageMap[prompt.Language] = true
		}
	}

	languages := make([]string, 0, len(languageMap))
	for lang := range languageMap {
		languages = append(languages, lang)
	}
	return languages
}

func (r *Registry) GetPromptsByLanguage(language string) []Prompt {
	var result []Prompt
	for _, prompt := range r.GetAllPrompts() {
		if prompt.Language == language {
			result = append(result, prompt)
		}
	}
	return result
}
