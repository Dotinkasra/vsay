package sub

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unsafe"
	"vsay/pkg/engine"
	"vsay/pkg/engine/dictionary"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type Dict struct {
	DictAdd
	DictDelete
}

type DictAdd struct {
	Cmd
}

type DictDelete struct {
	Cmd
}

func (scmd *DictAdd) GetFlags() []cli.Flag {
	lsFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "surface",
			Aliases:  []string{"w"},
			Usage:    "The surface form of the `word`.",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "pronunciation",
			Aliases:  []string{"y"},
			Usage:    "Pronunciation of words (`katakana`)",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "accent",
			Aliases:  []string{"a"},
			Usage:    "Accented type (refers to where the sound goes down)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "type",
			Aliases:  []string{"t"},
			Usage:    "One of the following: PROPER_NOUN, COMMON_NOUN, VERB, ADJECTIVE,SUFFIX",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "priority",
			Usage:    "Word priority (integer from `0 to 10`).",
			Required: false,
		},
	}
	return lsFlags
}

func (scmd *DictAdd) Action(c *cli.Context) error {
	e := engine.Engine{Host: c.String("host"), Port: c.Int("port")}
	dictRequest := dictionary.DictRequest{}
	dictRequest.Surface = c.String("surface")
	dictRequest.Pronunciation = c.String("pronunciation")
	dictRequest.AccentType = c.Int("accent")
	wordType := c.String("type")
	var convertWordType dictionary.WordType
	if wordType != "" {
		switch strings.ToUpper(wordType) {
		case "PROPER_NOUN":
			convertWordType = dictionary.PROPERNOUN
		case "COMMON_NOUN":
			convertWordType = dictionary.COMMONNOUN
		case "VERB":
			convertWordType = dictionary.VERB
		case "ADJECTIVE":
			convertWordType = dictionary.ADJECTIVE
		case "SUFFIX":
			convertWordType = dictionary.SUFFIX
		}
	}
	dictRequest.WordType = &convertWordType

	if c.Int("priority") == 0 {
		dictRequest.Priority = nil
	} else {
		p := c.Int("priority")
		dictRequest.Priority = &p
	}
	result, err := dictRequest.RegisterUserDict(e.MyHost())
	if err != nil {
		color.Red(fmt.Sprintln("Error"))
		log.Panic(err)
	}
	color.Green(fmt.Sprintln("Success"))
	fmt.Println(result)
	return nil
}

func (scmd *DictDelete) GetFlags() []cli.Flag {
	lsFlags := []cli.Flag{}
	return lsFlags
}

func (scmd *DictDelete) Action(c *cli.Context) error {
	e := engine.Engine{Host: c.String("host"), Port: c.Int("port")}
	uuid := c.Args().First()
	if uuid == "" {
		stdin := os.Stdin
		s, err := io.ReadAll(stdin)
		if err != nil {
			log.Panic(err)
		}
		uuid = *(*string)(unsafe.Pointer(&s))
	}
	uuid = strings.TrimSpace(uuid)
	uuid = strings.TrimSuffix(uuid, "\n")
	err := e.DeleteDict(uuid)
	if err != nil {
		log.Panic(err)
	}
	color.Green(fmt.Sprintln("Success"))

	return nil
}
