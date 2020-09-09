package parser

import (
	"fmt"
	"parser/types"
)

func DisplayResultsSummary() {
	addUpPartialFileAnalysisResults()
	displayTopLevelStatsInformation()
	displayPerLanguageStatsInformation()
}

func displayTopLevelStatsInformation() {
	fmt.Printf("Analyzed files: %d\n", analysisResults.ProcessedFilesCount)
	fmt.Printf("Skipped(unrecoginised) files: %d\n", analysisResults.ProcessedFilesCount)
	fmt.Printf("Failed to analyze files: %d\n", analysisResults.FailedToAnalyzeFilesCount)
}

func displayPerLanguageStatsInformation() {
	for languageName, information := range analysisResults.ResultsPerDetectedLanguage {
		if len(languageName) > 0 {
			println(languageName)
			fmt.Printf("\tSource code lines: %d\n", information.SourceCodeLinesCount)
			fmt.Printf("\tBlank lines lines: %d\n", information.BlankLinesCount)
			fmt.Printf("\tComment lines: %d\n", information.CommentLinesCount)
		}
	}
}

func addUpPartialFileAnalysisResults() {
	for _, fileResults := range perFileAnalysisResults {
		fileLanguageName := fileResults.LanguageName
		languageSummaryInformation, ok := analysisResults.ResultsPerDetectedLanguage[fileLanguageName]

		if !ok {
			analysisResults.ResultsPerDetectedLanguage[fileLanguageName] = fileResults
		} else {
			analysisResults.ResultsPerDetectedLanguage[fileLanguageName] = mergeFileResultsCounters(languageSummaryInformation, fileResults)
		}
	}
}

func mergeFileResultsCounters(firstFileResults types.FileAnalysisResult, secondFileResults types.FileAnalysisResult) types.FileAnalysisResult {
	var mergedFileResults types.FileAnalysisResult
	mergedFileResults.SourceCodeLinesCount = firstFileResults.SourceCodeLinesCount + secondFileResults.SourceCodeLinesCount
	mergedFileResults.BlankLinesCount = firstFileResults.BlankLinesCount + secondFileResults.BlankLinesCount
	mergedFileResults.CommentLinesCount = firstFileResults.CommentLinesCount + secondFileResults.CommentLinesCount
	return mergedFileResults
}
