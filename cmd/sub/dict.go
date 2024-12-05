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
	SubCommand
}

type DictDelete struct {
	SubCommand
}

func (scmd *DictAdd) GetFlags() []cli.Flag {
	lsFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "surface",
			Aliases:  []string{"w"},
			Usage:    "show speakers",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "pronunciation",
			Aliases:  []string{"y"},
			Usage:    "show dictionaries",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "accent",
			Aliases:  []string{"a"},
			Usage:    "",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "type",
			Aliases:  []string{"t"},
			Usage:    "",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "priority",
			Usage:    "",
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
	word_type := c.String("type")
	var convert_word_type dictionary.WordType
	if word_type != "" {
		switch strings.ToUpper(word_type) {
		case "PROPER_NOUN":
			convert_word_type = dictionary.PROPER_NOUN
		case "COMMON_NOUN":
			convert_word_type = dictionary.COMMON_NOUN
		case "VERB":
			convert_word_type = dictionary.VERB
		case "ADJECTIVE":
			convert_word_type = dictionary.ADJECTIVE
		case "SUFFIX":
			convert_word_type = dictionary.SUFFIX
		}
	}
	dictRequest.WordType = &convert_word_type

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
