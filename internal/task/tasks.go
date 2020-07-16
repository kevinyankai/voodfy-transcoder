package task

import (
	"context"
	"fmt"
	"log"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends/result"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/Voodfy/voodfy-transcoder/internal/logging"
	"github.com/opentracing/opentracing-go"
)

// Get return the tasks
func Get() map[string]interface{} {
	return map[string]interface{}{
		"long_running_task":               LongRunningTask,
		"extractAudioFromMp4Task":         ExtractAudioFromMp4Task,
		"removeAudioFromMp4Task":          RemoveAudioFromMp4Task,
		"thumbsPreviewGeneratorTask":      ThumbsPreviewGeneratorTask,
		"generateImageFromFrameVideoTask": GenerateImageFromFrameVideoTask,
		"ultraRenditionTask":              FullHDFallbackTask,
		"hdRenditionTask":                 HDFallBackTask,
		"midRenditionTask":                MidDefinitionFallbackTask,
		"standardRenditionTask":           StandardFallbackTask,
		"lowRenditionTask":                LowDefinitionTask,
		"renditionTask":                   RenditionTask,
		"sendDirToIPFSTask":               SendDirToIPFSTask,
	}
}

// AsyncResultArray slice of AsyncResult
type AsyncResultArray []result.AsyncResult

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

	chain, _ := tasks.NewChain(
		&removeAudioTask, &extractAudioTask,
		&generateImageFromFrameVideoTask, &thumbsPreviewTask, &livepeerTask)

	_, err := server.SendChain(chain)

	if err != nil {
		log.Fatal(err)
	}
}

// Local task to use ffmpeg
func Local(resourceID, resourceName, directory, tracker string, server *machinery.Server) AsyncResultArray {
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

	chain, err := tasks.NewChain(
		&removeAudioTask, &extractAudioTask,
		&generateImageFromFrameVideoTask, &thumbsPreviewTask, &lowRenditionTask,
		&standardRenditionTask, &midRenditionTask, &hdRenditionTask, &ultraHdRenditionTask)

	if err != nil {
		log.Panic(err)
	}

	_, err = server.SendChain(chain)

	if err != nil {
		log.Panic(err)
	}

	ipfs := result.NewAsyncResult(&ultraHdRenditionTask, server.GetBackend())

	var a AsyncResultArray
	a = append(a, *ipfs)

	if err != nil {
		log.Println(fmt.Sprintf("Could not send task: %s", err.Error()))
	}

	return a
}

// IPFSAddDir send the directory to ipfs
func IPFSAddDir(resourceID, directory, tracker string, server *machinery.Server) *tasks.TaskState {
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
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()
	asyncResult, err := server.SendTaskWithContext(ctx, &longRunningTask)
	if err != nil {
		log.Println(fmt.Sprintf("Could not send task: %s", err.Error()))
	}
	return asyncResult.GetState()
}
