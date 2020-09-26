package queue

import (
	"fmt"
	"time"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/Voodfy/voodfy-transcoder/internal/influxdbclient"
	"github.com/Voodfy/voodfy-transcoder/internal/settings"
	"github.com/Voodfy/voodfy-transcoder/internal/task"
	"github.com/Voodfy/voodfy-transcoder/internal/utils"
	"github.com/Voodfy/voodfy-transcoder/pkg/logging"
)

func startServer() (*machinery.Server, error) {
	var cnf = &config.Config{
		Broker:        fmt.Sprintf("redis://%s/0", settings.RedisSetting.Host),
		DefaultQueue:  "transcoder_tasks",
		ResultBackend: fmt.Sprintf("redis://%s/1", settings.RedisSetting.Host),
	}
	server, _ := machinery.NewServer(cnf)

	// Register tasks
	tasks := task.Get()

	return server, server.RegisterTasks(tasks)
}

// NewWorker return a instance of a worker
func NewWorker() *machinery.Worker {
	var start time.Time
	var finished float64

	consumerTag := settings.AppSetting.Tag

	server, err := startServer()
	if err != nil {
		utils.SendError("startServer", err)
	}

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker(consumerTag, 0)
	influx := influxdbclient.NewClient()

	// Here we inject some custom code for error handling,
	// start and end of task hooks, useful for metrics for example.
	errorhandler := func(err error) {
		utils.SendError("I am an error handler:", err)
	}

	pretaskhandler := func(signature *tasks.Signature) {
		start = time.Now()
		logging.Info(fmt.Sprintf("I am a start of task handler for: %s", signature.Name))
	}

	posttaskhandler := func(signature *tasks.Signature) {
		finished = time.Since(start).Seconds()

		if len(signature.Args) != 0 {
			for _, arg := range signature.Args {
				if arg.Name == "id" {
					influx.Send(arg.Value, signature.Name, fmt.Sprintf("%f", finished))
				}
			}
		}

		logging.Info(fmt.Sprintf("I am an end of task handler for: %s", signature.Name))
	}

	worker.SetPostTaskHandler(posttaskhandler)
	worker.SetErrorHandler(errorhandler)
	worker.SetPreTaskHandler(pretaskhandler)

	return worker
}
