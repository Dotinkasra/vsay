package sub

import "github.com/urfave/cli/v2"

type SubCommand interface {
	Action(c *cli.Context) error
	GetFlags() []cli.Flag
}
