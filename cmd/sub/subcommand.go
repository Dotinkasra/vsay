package sub

import "github.com/urfave/cli/v2"

type Cmd interface {
	Action(c *cli.Context) error
	GetFlags() []cli.Flag
}
