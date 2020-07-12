package tasks

import (
	"context"
	"fmt"
	"log"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/Voodfy/voodfy-transcoder/internal/logging"
	"github.com/opentracing/opentracing-go"
)

// Ping send ping to queue
func Ping(server machinery.Server) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()

	longRunningTask := tasks.Signature{
		Name: "long_running_task",
	}

	_, err := server.SendTaskWithContext(ctx, &longRunningTask)
	if err != nil {
		logging.Info(fmt.Sprintf("Could not send task: %s", err.Error()))
	}
}

// LivepeerChain send the chunks to livepeer
func LivepeerChain(resourceID, resourceName, directory, tracker string, server *machinery.Server) {
	removeAudioTask := tasks.Signature{
		Name: "removeAudioFromMp4Task",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s", resourceName),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	extractAudioTask := tasks.Signature{
		Name: "extractAudioFromMp4Task",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s", resourceName),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	generateImageFromFrameVideoTask := tasks.Signature{
		Name: "generateImageFromFrameVideoTask",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s_without_audio.mp4", resourceID),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	thumbsPreviewTask := tasks.Signature{
		Name: "thumbsPreviewGeneratorTask",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s", resourceName),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	livepeerTask := tasks.Signature{
		Name: "renditionTask",
		Args: []tasks.Arg{
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "profile",
				Type:  "string",
				Value: "ultra.json",
			},
		},
	}

	longRunningTask := tasks.Signature{
		Name: "sendDirToIPFSTask",
		Args: []tasks.Arg{
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
		},
	}

	chain, _ := tasks.NewChain(
		&removeAudioTask, &extractAudioTask,
		&generateImageFromFrameVideoTask, &thumbsPreviewTask, &livepeerTask, &longRunningTask)

	_, err := server.SendChain(chain)

	if err != nil {
		log.Fatal(err)
	}
}

// Local task to use ffmpeg
func Local(resourceID, resourceName, directory, tracker string, server *machinery.Server) {
	removeAudioTask := tasks.Signature{
		Name: "removeAudioFromMp4Task",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s", resourceName),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	extractAudioTask := tasks.Signature{
		Name: "extractAudioFromMp4Task",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s", resourceName),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	generateImageFromFrameVideoTask := tasks.Signature{
		Name: "generateImageFromFrameVideoTask",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s_without_audio.mp4", resourceID),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	thumbsPreviewTask := tasks.Signature{
		Name: "thumbsPreviewGeneratorTask",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s", resourceName),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	lowRenditionTask := tasks.Signature{
		Name: "lowRenditionTask",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s_without_audio.mp4", resourceID),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	standardRenditionTask := tasks.Signature{
		Name: "standardRenditionTask",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s_without_audio.mp4", resourceID),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	midRenditionTask := tasks.Signature{
		Name: "midRenditionTask",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s_without_audio.mp4", resourceID),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	hdRenditionTask := tasks.Signature{
		Name: "hdRenditionTask",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s_without_audio.mp4", resourceID),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	ultraHdRenditionTask := tasks.Signature{
		Name: "ultraRenditionTask",
		Args: []tasks.Arg{
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
			{
				Name:  "filename",
				Type:  "string",
				Value: fmt.Sprintf("%s_without_audio.mp4", resourceID),
			},
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
		},
	}

	longRunningTask := tasks.Signature{
		Name: "sendDirToIPFSTask",
		Args: []tasks.Arg{
			{
				Name:  "directory",
				Type:  "string",
				Value: directory,
			},
			{
				Name:  "tracker",
				Type:  "string",
				Value: tracker,
			},
			{
				Name:  "id",
				Type:  "string",
				Value: resourceID,
			},
		},
	}

	chain, err := tasks.NewChain(
		&removeAudioTask, &extractAudioTask,
		&generateImageFromFrameVideoTask, &thumbsPreviewTask, &lowRenditionTask,
		&standardRenditionTask, &midRenditionTask, &hdRenditionTask, &ultraHdRenditionTask,
		&longRunningTask)

	if err != nil {
		log.Panic(err)
	}

	tasks, err := server.SendChain(chain)

	if err != nil {
		log.Panic(err)
	}

	log.Println("tasks", tasks)

}
