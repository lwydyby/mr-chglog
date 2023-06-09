package config

type MRChLogConfig struct {
	Style         string `yaml:"style"`
	Template      string `yaml:"template"`
	Title         string `yaml:"title"`
	RepositoryURL string `yaml:"repository_url"`
	Token         string `yaml:"token"`
	POEToken      string `yaml:"poe_token"`
}
