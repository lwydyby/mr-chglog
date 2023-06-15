package cmd

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"time"

	"github.com/fatih/color"

	"github.com/lwydyby/mr-chglog/config"
)

type CLI struct {
	ctx          *CLIContext
	fs           FileSystem
	configLoader ConfigLoader
	generator    Generator
}

func NewCLI(
	ctx *CLIContext, fs FileSystem,
	configLoader ConfigLoader,
	generator Generator,
) *CLI {
	return &CLI{
		ctx:          ctx,
		fs:           fs,
		configLoader: configLoader,
		generator:    generator,
	}
}

func (c *CLI) Run() int {
	start := time.Now()

	if c.ctx.NoColor {
		color.NoColor = true
	}

	log.Println(":watch: Generating changelog ...")

	config, err := c.prepareConfig()
	if err != nil {
		panic(err)
	}

	w, err := c.createOutputWriter()
	if err != nil {
		panic(err)
	}

	err = c.generator.Generate(w, c.ctx, config)
	if err != nil {
		panic(err)
	}

	fmt.Printf(":sparkles: Generate of %s is completed! (%s)",
		color.GreenString("\""+c.ctx.OutputPath+"\""),
		color.New(color.Bold).SprintFunc()(time.Since(start).String()),
	)

	return 1
}

func (c *CLI) prepareConfig() (*config.MRChLogConfig, error) {
	cfg, err := c.configLoader.Load(c.ctx.ConfigPath)
	if err != nil {
		return nil, err
	}
	if len(c.ctx.AppID) != 0 {
		cfg.AppID = c.ctx.AppID
	}
	if len(c.ctx.AppSecret) != 0 {
		cfg.AppSecret = c.ctx.AppSecret
	}
	if len(c.ctx.ChatID) != 0 && len(c.ctx.ChatID[0]) != 0 {
		cfg.ChatID = c.ctx.ChatID
	}
	if len(c.ctx.BotTitle) != 0 {
		cfg.BotTitle = c.ctx.BotTitle
	}
	if !filepath.IsAbs(cfg.Template) {
		cfg.Template = filepath.Join(filepath.Dir(c.ctx.ConfigPath), cfg.Template)
	}
	if len(c.ctx.Template) != 0 {
		cfg.Template = c.ctx.Template
	}
	if len(c.ctx.RepositoryURL) != 0 {
		cfg.RepositoryURL = c.ctx.RepositoryURL
	}
	if len(c.ctx.Token) != 0 {
		cfg.Token = c.ctx.Token
	}

	return cfg, nil
}

func (c *CLI) createOutputWriter() (io.Writer, error) {
	if c.ctx.OutputPath == "" {
		return c.ctx.Stdout, nil
	}

	out := c.ctx.OutputPath
	dir := filepath.Dir(out)
	err := c.fs.MkdirP(dir)
	if err != nil {
		return nil, err
	}

	file, err := c.fs.Create(out)
	if err != nil {
		return nil, err
	}

	return file, nil
}
