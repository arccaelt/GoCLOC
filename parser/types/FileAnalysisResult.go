package types

type FileAnalysisResult struct {
	LanguageName         string
	SourceCodeLinesCount int
	BlankLinesCount      int
	CommentLinesCount    int
}
