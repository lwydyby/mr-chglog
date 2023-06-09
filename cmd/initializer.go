package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
)

// Initializer ...
type Initializer struct {
	ctx                    *InitContext
	fs                     FileSystem
	questioner             Questioner
	configBuilder          ConfigBuilder
	templateBuilderFactory TemplateBuilderFactory
}

// NewInitializer ...
func NewInitializer(
	ctx *InitContext,
	fs FileSystem,
	questioner Questioner,
	configBuilder ConfigBuilder,
	tplBuilderFactory TemplateBuilderFactory) *Initializer {
	return &Initializer{
		ctx:                    ctx,
		fs:                     fs,
		questioner:             questioner,
		configBuilder:          configBuilder,
		templateBuilderFactory: tplBuilderFactory,
	}
}

// Run ...
func (init *Initializer) Run() int {
	ans, err := init.questioner.Ask()
	if err != nil {
		panic(err)
	}

	if err = init.fs.MkdirP(filepath.Join(init.ctx.WorkingDir, ans.ConfigDir)); err != nil {
		panic(err)
	}

	if err = init.generateConfig(ans); err != nil {
		panic(err)
	}

	if err = init.generateTemplate(ans); err != nil {
		panic(err)
	}

	success := color.CyanString("✔")
	fmt.Printf(`
:sparkles: %s
  %s %s
  %s %s
`,
		color.GreenString("Configuration file and template generation completed!"),
		success,
		filepath.Join(ans.ConfigDir, defaultConfigFilename),
		success,
		filepath.Join(ans.ConfigDir, defaultTemplateFilename),
	)

	return 0
}

func (init *Initializer) generateConfig(ans *Answer) error {
	s, err := init.configBuilder.Build(ans)
	if err != nil {
		return err
	}

	return init.fs.WriteFile(filepath.Join(init.ctx.WorkingDir, ans.ConfigDir, defaultConfigFilename), []byte(s))
}

func (init *Initializer) generateTemplate(ans *Answer) error {
	templateBuilder := init.templateBuilderFactory(ans.Template)
	s, err := templateBuilder.Build(ans)
	if err != nil {
		return err
	}

	return init.fs.WriteFile(filepath.Join(init.ctx.WorkingDir, ans.ConfigDir, defaultTemplateFilename), []byte(s))
}
