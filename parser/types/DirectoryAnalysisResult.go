package types

type DirectoryAnalysisResult struct {
	ProcessedFilesCount        uint32
	SkippedFilesCount          uint32
	FailedToAnalyzeFilesCount  uint32
	ResultsPerDetectedLanguage map[string]FileAnalysisResult
}
