package parser

import (
	"sync"
	"parser/types"
	"path/filepath"
)

const defaultFileAnalysisResultsSize int = 100

var goroutinesWaitGroup sync.WaitGroup
var cLikeCommentSymbols []types.Comment = []types.Comment{
	types.Comment{StartSymbol: "//", EndSymbol: "", CanSpanMultipleLines: false},
	types.Comment{StartSymbol: "/*", EndSymbol: "*/", CanSpanMultipleLines: true},
}
var analysisResults types.DirectoryAnalysisResult
var perFileAnalysisResults []types.FileAnalysisResult
var supportedLanguages map[string]types.ProgrammingLanguageFeatures = map[string]types.ProgrammingLanguageFeatures{
	"py": types.ProgrammingLanguageFeatures{
		LanguageName:            "Python",
		CommentSymbols: []types.Comment{types.Comment{StartSymbol: "#", EndSymbol : "", CanSpanMultipleLines: false}}},
	"go": types.ProgrammingLanguageFeatures{
		LanguageName:            "Go",
		CommentSymbols: cLikeCommentSymbols},
	"c": types.ProgrammingLanguageFeatures{
		LanguageName:            "C",
		CommentSymbols: cLikeCommentSymbols},
	"java": types.ProgrammingLanguageFeatures{
		LanguageName:            "Java",
		CommentSymbols: cLikeCommentSymbols},
}

func init() {
	perFileAnalysisResults = make([]types.FileAnalysisResult, defaultFileAnalysisResultsSize)
	analysisResults.ResultsPerDetectedLanguage = make(map[string]types.FileAnalysisResult)
}

func DirectoryAnalyzer(inputPath string) {
	filepath.Walk(inputPath, ProcessDirectoryEntry)
	goroutinesWaitGroup.Wait()
}
