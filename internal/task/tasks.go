package task

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
