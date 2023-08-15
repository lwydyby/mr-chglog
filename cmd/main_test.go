package cmd

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func getMockAppAction(t *testing.T) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		assert.Equal(t, "c.yml", c.String("config"))
		assert.Equal(t, "o.md", c.String("output"))
		assert.Equal(t, "v5", c.String("next-tag"))
		return nil
	}
}

func TestCreateApp(t *testing.T) {
	app := CreateApp(getMockAppAction(t))
	args := []string{
		"mr-changelog ",
		"--config", "c.yml",
		"--output", "o.md",
		"--next-tag", "v5",
	}
	err := app.Run(args)
	if err != nil {
		log.Fatal(err)
	}
}
