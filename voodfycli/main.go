package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/Voodfy/voodfy-transcoder/internal/task"
	"github.com/Voodfy/voodfy-transcoder/voodfycli/tasks"
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
	server, err := startServer()
	if err != nil {
		log.Fatal(err)
	}

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a video to transcode",
			Action: func(c *cli.Context) error {
				if c.Args().Get(0) == "livepeer" {
					tasks.LivepeerChain(c.Args().Get(1), c.Args().Get(2), c.Args().Get(3), c.Args().Get(4), server)
				}
				if c.Args().Get(0) == "local" {
					tasks.Local(c.Args().Get(1), c.Args().Get(2), c.Args().Get(3), c.Args().Get(4), server)
				}
				fmt.Println("added video: ", c.Args().First())
				return nil
			},
		},
		{
			Name:    "ping",
			Aliases: []string{"p"},
			Usage:   "ping the queue",
			Action: func(c *cli.Context) error {
				tasks.Ping(*server)
				return nil
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
