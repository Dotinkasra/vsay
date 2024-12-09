package sub

import (
	"vsay/pkg/engine"
	"vsay/pkg/engine/speaker"

	"github.com/urfave/cli/v2"
)

type Install struct {
	Cmd
}

func (scmd Install) GetFlags() []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "path",
			Aliases:  []string{"i"},
			Usage:    "The surface form of the `word`.",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "link",
			Aliases:  []string{"l"},
			Usage:    "Pronunciation of words (`katakana`)",
			Required: false,
		},
	}
	return flags
}

func (scmd Install) Action(c *cli.Context) error {
	e := engine.Engine{Host: c.String("host"), Port: c.Int("port")}
	_, err := speaker.InstallAivmModels(e.MyHost(), c.String("path"), c.String("link"))
	if err != nil {
		return err
	}
	return nil
}
