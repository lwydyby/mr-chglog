package config

type MRChLogConfig struct {
	Style         string   `yaml:"style"`
	Template      string   `yaml:"template"`
	Title         string   `yaml:"title"`
	RepositoryURL string   `yaml:"repository_url"`
	Token         string   `yaml:"token"`
	POEToken      string   `yaml:"poe_token"`
	NeedRobot     bool     `yaml:"need_robot"`
	AppID         string   `yaml:"app_id"`
	AppSecret     string   `yaml:"app_secret"`
	ChatID        []string `yaml:"chat_id"`
	BotTitle      string   `yaml:"bot_title"`
}
