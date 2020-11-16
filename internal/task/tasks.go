package task

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends/result"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/Voodfy/voodfy-transcoder/pkg/logging"
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
		"fallbackRenditionTask":           FallbackRenditionTask,
		"renditionTask":                   RenditionTask,
		"sendDirToIPFSTask":               SendDirToIPFSTask,
		"sendDirToFilecoinTask":           SendDirToFilecoinTask,
		"ffprobeTask":                     FFprobeTask,
		"convertToMp4Task":                ConvertToMp4Task,
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

// Local task to use ffmpeg
func Local(resourceID, resourceName, directory, tracker string, server *machinery.Server) AsyncResultArray {
	src := fmt.Sprintf("%s/%s/", directory, tracker)
	dstFiles := fmt.Sprintf("%s%s_ipfs/", src, resourceID)
	os.MkdirAll(dstFiles, 0777)

	log.Println("input: ------->", src)
	log.Println("oputput: ----->", dstFiles)

	removeAudioTask := tasks.Signature{
		Name: "removeAudioFromMp4Task",
		Args: []tasks.Arg{
			{
				Name:  "input",
				Type:  "string",
				Value: fmt.Sprintf("%s%s", src, resourceName),
			},
			{
				Name:  "output",
				Type:  "string",
				Value: fmt.Sprintf("%s%s_without_audio.mp4", src, resourceID),
			},
		},
	}

	extractAudioTask := tasks.Signature{
		Name: "extractAudioFromMp4Task",
		Args: []tasks.Arg{
			{
				Name:  "input",
				Type:  "string",
				Value: fmt.Sprintf("%s%s", src, resourceName),
			},
			{
				Name:  "output",
				Type:  "string",
				Value: fmt.Sprintf("%s%s_ipfs/%s_a1.m4a", src, resourceID, resourceID),
			},
		},
	}

	generateImageFromFrameVideoTask := tasks.Signature{
		Name: "generateImageFromFrameVideoTask",
		Args: []tasks.Arg{
			{
				Name:  "input",
				Type:  "string",
				Value: fmt.Sprintf("%s%s_without_audio.mp4", src, resourceID),
			},
			{
				Name:  "output",
				Type:  "string",
				Value: dstFiles,
			},
		},
	}

	thumbsPreviewTask := tasks.Signature{
		Name: "thumbsPreviewGeneratorTask",
		Args: []tasks.Arg{
			{
				Name:  "input",
				Type:  "string",
				Value: fmt.Sprintf("%s%s", src, resourceName),
			},
			{
				Name:  "output",
				Type:  "string",
				Value: dstFiles,
			},
		},
	}

	lowRenditionTask := tasks.Signature{
		Name: "fallbackRenditionTask",
		Args: []tasks.Arg{
			{
				Name:  "input",
				Type:  "string",
				Value: fmt.Sprintf("%s%s", src, resourceName),
			},
			{
				Name:  "output",
				Type:  "string",
				Value: fmt.Sprintf("%s%s_v3.mp4", dstFiles, resourceID),
			},
			{
				Name:  "fnc",
				Type:  "string",
				Value: "240p",
			},
		},
	}

	standardRenditionTask := tasks.Signature{
		Name: "fallbackRenditionTask",
		Args: []tasks.Arg{
			{
				Name:  "input",
				Type:  "string",
				Value: fmt.Sprintf("%s%s", src, resourceName),
			},
			{
				Name:  "output",
				Type:  "string",
				Value: fmt.Sprintf("%s%s_v4.mp4", dstFiles, resourceID),
			},
			{
				Name:  "fnc",
				Type:  "string",
				Value: "360p",
			},
		},
	}

	midRenditionTask := tasks.Signature{
		Name: "fallbackRenditionTask",
		Args: []tasks.Arg{
			{
				Name:  "input",
				Type:  "string",
				Value: fmt.Sprintf("%s%s", src, resourceName),
			},
			{
				Name:  "output",
				Type:  "string",
				Value: fmt.Sprintf("%s%s_v5.mp4", dstFiles, resourceID),
			},
			{
				Name:  "fnc",
				Type:  "string",
				Value: "480p",
			},
		},
	}

	hdRenditionTask := tasks.Signature{
		Name: "fallbackRenditionTask",
		Args: []tasks.Arg{
			{
				Name:  "input",
				Type:  "string",
				Value: fmt.Sprintf("%s%s", src, resourceName),
			},
			{
				Name:  "output",
				Type:  "string",
				Value: fmt.Sprintf("%s%s_v6.mp4", dstFiles, resourceID),
			},
			{
				Name:  "fnc",
				Type:  "string",
				Value: "720p",
			},
		},
	}

	ultraHdRenditionTask := tasks.Signature{
		Name: "fallbackRenditionTask",
		Args: []tasks.Arg{
			{
				Name:  "input",
				Type:  "string",
				Value: fmt.Sprintf("%s%s", src, resourceName),
			},
			{
				Name:  "output",
				Type:  "string",
				Value: fmt.Sprintf("%s%s_v7.mp4", dstFiles, resourceID),
			},
			{
				Name:  "fnc",
				Type:  "string",
				Value: "1080p",
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
func IPFSAddDir(directory, resourceID string, server *machinery.Server) *tasks.TaskState {
	longRunningTask := tasks.Signature{
		Name: "sendDirToIPFSTask",
		Args: []tasks.Arg{
			{
				Name:  "output",
				Type:  "string",
				Value: directory,
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
