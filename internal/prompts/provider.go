package prompts

type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
	Content     string           `json:"content"`
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
	var allPrompts []Prompt
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
