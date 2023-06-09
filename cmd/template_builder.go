package cmd

const templateTagNameAnchor = ""

// TemplateBuilder ...
type TemplateBuilder interface {
	Builder
}

// TemplateBuilderFactory ...
type TemplateBuilderFactory = func(string) TemplateBuilder

func templateBuilderFactory(template string) TemplateBuilder {
	return NewDefaultTemplateBuilder()
}
