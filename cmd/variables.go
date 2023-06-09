package cmd

// Previewable ...
type Previewable interface {
	Display() string
	Preview() string
}

// Defaults
var (
	defaultConfigDir        = ".chglog"
	defaultConfigFilename   = "mr_config.yml"
	defaultTemplateFilename = "MRCHANGELOG.tpl.md"
)

// Styles
var (
	styleGitHub    = "github"
	styleGitLab    = "gitlab"
	styleBitbucket = "bitbucket"
	styleNone      = "none"
	styles         = []string{
		styleGitHub,
		styleGitLab,
		styleBitbucket,
		styleNone,
	}
)

// TemplateStyleFormat ...
type TemplateStyleFormat struct {
	preview string
	display string
}

// Display ...
func (t *TemplateStyleFormat) Display() string {
	return t.display
}

// Preview ...
func (t *TemplateStyleFormat) Preview() string {
	return t.preview
}

// Templates
var (
	tplDefaultChangelog = &TemplateStyleFormat{
		display: "default-changelog",
	}
	templates = []Previewable{
		tplDefaultChangelog,
	}
)
