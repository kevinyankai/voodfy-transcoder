package main

import (
	"bufio"
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
	app.Author = "Voodfy"
	app.Email = "contact@voodfy.com"
	app.Version = "0.0.2"
}

func main() {
	var err error

	server, err := startServer()
	if err != nil {
		log.Fatal(err)
	}

	app.Commands = []cli.Command{
		{
			Name:    "signup",
			Aliases: []string{"s"},
			Usage:   "setup an account at Voodfy",
			Action: func(c *cli.Context) error {
				pwd := task.ManagerSetupAccountVoodfy()
				log.Println(pwd)
				return nil
			},
		},
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "login at Voodfy",
			Action: func(c *cli.Context) error {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter secret hash: ")
				pwd, _ := reader.ReadString('\n')
				log.Println(task.ManagerLoginVoodfy(pwd))
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
					c.Args().Get(0), c.Args().Get(1), server)
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

				log.Println("Directory ID:", directory.ID)
				log.Println("Directory CID:", directory.CID)

				for _, r := range directory.Resources {
					log.Println("Resource ID:", r.ID)
					log.Println("Resource Jid:", r.Jid)
					log.Println("Resource Name:", r.Name)
					log.Println("Resource CID:", r.CID)
				}

				return nil
			},
		},
		{
			Name:    "store_config",
			Aliases: []string{"sc"},
			Usage:   "show the default config at Filecoin",
			Action: func(c *cli.Context) error {
				_, token, _ := powergate.FFSCreate()
				addr := "127.0.0.1:5002"
				powergate.FFSDefaultConfig(token, addr)
				return nil
			},
		},
		{
			Name:    "store",
			Aliases: []string{"st"},
			Usage:   "store the resources on Filecoin",
			Action: func(c *cli.Context) error {
				task.ManagerPowergate(c.Args().Get(0))
				return nil
			},
		},
		{
			Name:    "store",
			Aliases: []string{"st"},
			Usage:   "store the resources on Filecoin",
			Action: func(c *cli.Context) error {
				task.ManagerPowergate(c.Args().Get(0))
				return nil
			},
		},
		{
			Name:    "embed",
			Aliases: []string{"eb"},
			Usage:   "retrieve an embed from Voodfy",
			Action: func(c *cli.Context) error {
				log.Println(task.ManagerEmbedByVoodfy(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2)))
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
