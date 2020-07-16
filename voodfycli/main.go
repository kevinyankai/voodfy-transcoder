package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/Voodfy/voodfy-transcoder/internal/models"
	"github.com/Voodfy/voodfy-transcoder/internal/task"
	"github.com/Voodfy/voodfy-transcoder/pkg/powergate"
	"github.com/urfave/cli"
)

var (
	app *cli.App
)

// get return the tasks
func get() map[string]interface{} {
	return map[string]interface{}{}
}

func startServer() (*machinery.Server, error) {
	var cnf = &config.Config{
		Broker:        fmt.Sprintf("redis://%s/0", os.Getenv("REDIS_BROKER")),
		DefaultQueue:  "transcoder_tasks",
		ResultBackend: fmt.Sprintf("redis://%s/1", os.Getenv("REDIS_RESULT")),
	}
	server, _ := machinery.NewServer(cnf)

	// Register tasks
	tasks := task.Get()

	return server, server.RegisterTasks(tasks)
}

func init() {
	// Initialise a CLI app
	app = cli.NewApp()
	app.Name = "voodfycli"
	app.Usage = "voodfycli it is the command line interface to add task on voodfy transcoder"
	app.Author = "Leandro Barbosa"
	app.Email = "contact@voodfy.com"
	app.Version = "0.0.1"
}

func main() {
	var id string
	var token string
	var err error

	server, err := startServer()
	if err != nil {
		log.Fatal(err)
	}

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"a"},
			Usage:   "init the default configurations",
			Action: func(c *cli.Context) error {
				id, token, err = powergate.FFSCreate()
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a video to transcode",
			Action: func(c *cli.Context) error {
				task.ManagerTranscoder(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2),
					c.Args().Get(3), c.Args().Get(4), server)
				return nil
			},
		},
		{
			Name:    "ipfs",
			Aliases: []string{"ipfs"},
			Usage:   "send the result video transcoded to IPFS",
			Action: func(c *cli.Context) error {
				task.ManagerIPFS(
					c.Args().Get(0),
					c.Args().Get(1), c.Args().Get(2), server)
				return nil
			},
		},
		{
			Name:    "directory",
			Aliases: []string{"dt"},
			Usage:   "get a directory giving the resource id",
			Action: func(c *cli.Context) error {
				directory := models.Directory{
					ID: c.Args().Get(0),
				}
				directory.Get()
				log.Println("Directory:", directory.ID)
				for _, r := range directory.Resources {
					log.Println("Resource:", r)
				}

				return nil
			},
		},
		{
			Name:    "store_config",
			Aliases: []string{"sc"},
			Usage:   "show the default config at Filecoin",
			Action: func(c *cli.Context) error {
				id, token, err = powergate.FFSCreate()
				powergate.FFSDefaultConfig(token)
				return nil
			},
		},

		{
			Name:    "store",
			Aliases: []string{"st"},
			Usage:   "store the resources on Filecoin",
			Action: func(c *cli.Context) error {
				id, token, err = powergate.FFSCreate()

				directory := models.Directory{
					ID: c.Args().Get(0),
				}
				directory.Get()

				for idx, r := range directory.Resources {
					jid := powergate.FFSPush(r.CID, token)
					r.Jid = jid
					directory.Resources[idx] = r
				}

				directory.Save()

				return nil
			},
		},
		{
			Name:    "ping",
			Aliases: []string{"p"},
			Usage:   "ping the queue",
			Action: func(c *cli.Context) error {
				task.Ping(*server)
				return nil
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
