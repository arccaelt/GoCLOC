package parser

import (
	"bufio"
	"errors"
	"io"
	"os"
	"parser/constants"
	"parser/types"
	"strings"
	"sync"
	"sync/atomic"
)

const fileExtensionSeparator string = "."
const tokensArrayFileExtensionIndex int = 1
const eofDelim byte = '\n'

var mutex sync.Mutex

func ProcessDirectoryEntry(path string, info os.FileInfo, err error) error {
	if err != nil {
		analysisResults.FailedToAnalyzeFilesCount++
		return err
	}

	if info.IsDir() {
		return nil
	}

	fileName := info.Name()
	fileExtension, err := getFileExtension(fileName)

	if err != nil {
		analysisResults.FailedToAnalyzeFilesCount++
		return nil
	}

	if !isSupportedLanguage(fileExtension) {
		analysisResults.SkippedFilesCount++
		return nil
	}

	startFileParsing(path, fileExtension)
	return nil
}

func startFileParsing(path string, fileExtension string) {
	languageFeatures := supportedLanguages[fileExtension]
	startGoroutineForFileAnalysis(path, languageFeatures)
}

func startGoroutineForFileAnalysis(path string, languageFeatures types.ProgrammingLanguageFeatures) {
	go analyzeFile(path, languageFeatures)
	goroutinesWaitGroup.Add(1)
}

func analyzeFile(path string, languageFeatures types.ProgrammingLanguageFeatures) {
	fileDescriptor, err := os.Open(path)

	if err != nil {
		atomic.AddUint32(&analysisResults.FailedToAnalyzeFilesCount, 1)
		return
	}

	defer fileDescriptor.Close()

	commentSymbols := languageFeatures.CommentSymbols

	isCurrentlyInComment := false
	fileAnalysisResult := types.FileAnalysisResult{LanguageName: languageFeatures.LanguageName}
	bufferedFileReader := bufio.NewReader(fileDescriptor)

	for {
		line, err := bufferedFileReader.ReadString(eofDelim)
		if err != nil {
			if err != io.EOF {
				atomic.AddUint32(&analysisResults.FailedToAnalyzeFilesCount, 1)
				return
			}
			break
		}

		line = strings.TrimSpace(line)
		lineType := getLineType(line, commentSymbols, isCurrentlyInComment)
		updateCurrentLineState(line, lineType, &isCurrentlyInComment, commentSymbols)
		updateFileAnalysisResults(&fileAnalysisResult, lineType, isCurrentlyInComment)
	}

	updateGlobalAnalysisResults(fileAnalysisResult)
	goroutinesWaitGroup.Done()
}

func updateCurrentLineState(line string, lineType int, isCurrentlyInComment *bool, commentSymbols []types.Comment) {
	checkAndUpdateIfIsTheEndOfAComment(line, commentSymbols, isCurrentlyInComment)
	checkAndUpdateIfIsTheStartOfAMultilineAComment(line, commentSymbols, isCurrentlyInComment)
}

func checkAndUpdateIfIsTheEndOfAComment(line string, commentSymbols []types.Comment, isCurrentlyInComment *bool) {
	if *isCurrentlyInComment && isEndComment(line, commentSymbols) {
		*isCurrentlyInComment = false
	}
}

func checkAndUpdateIfIsTheStartOfAMultilineAComment(line string, commentSymbols []types.Comment, isCurrentlyInComment *bool) {
	if isStartComment(line, commentSymbols) && canSpanMultipleLines(line, commentSymbols) {
		*isCurrentlyInComment = true
	}
}

func updateFileAnalysisResults(fileAnalysisResults *types.FileAnalysisResult, lineType int, isCurrentlyInComment bool) {
	switch lineType {
	case constants.COMMENT:
		(*fileAnalysisResults).CommentLinesCount++
	case constants.BLANK:
		(*fileAnalysisResults).BlankLinesCount++
	case constants.SOURCE_CODE:
		if !isCurrentlyInComment {
			(*fileAnalysisResults).SourceCodeLinesCount++
		}
	}
}

func getLineType(line string, commentSymbols []types.Comment, isCurrentlyInComment bool) int {
	if isCurrentlyInComment {
		return constants.COMMENT
	}
	return searchForLineType(line, commentSymbols)
}

func searchForLineType(line string, commentSymbols []types.Comment) int {
	if isBlank(line) {
		return constants.BLANK
	} else if isComment(line, commentSymbols) {
		return constants.COMMENT
	} else {
		return constants.SOURCE_CODE
	}
}

func updateGlobalAnalysisResults(analyzedFileResults types.FileAnalysisResult) {
	addAnalyzedFileResultsToGlobalList(analyzedFileResults)
	atomic.AddUint32(&analysisResults.ProcessedFilesCount, 1)
}

func addAnalyzedFileResultsToGlobalList(analyzedFileResults types.FileAnalysisResult) {
	mutex.Lock()
	perFileAnalysisResults = append(perFileAnalysisResults, analyzedFileResults)
	mutex.Unlock()
}

func canSpanMultipleLines(lineWithComment string, commentSymbols []types.Comment) bool {
	for _, commentSymbol := range commentSymbols {
		if strings.HasPrefix(lineWithComment, commentSymbol.StartSymbol) ||
			(len(commentSymbol.EndSymbol) > 0 && strings.HasSuffix(lineWithComment, commentSymbol.EndSymbol)) {
			return commentSymbol.CanSpanMultipleLines
		}
	}
	return false
}

func isComment(line string, languageCommentSymbols []types.Comment) bool {
	return isStartComment(line, languageCommentSymbols) || isEndComment(line, languageCommentSymbols)
}

func isStartComment(line string, languageCommentSymbols []types.Comment) bool {
	for _, commentSymbols := range languageCommentSymbols {
		if strings.HasPrefix(line, commentSymbols.StartSymbol) {
			return true
		}
	}
	return false
}

func isEndComment(line string, languageCommentSymbols []types.Comment) bool {
	for _, commentSymbols := range languageCommentSymbols {
		if commentSymbols.EndSymbol != "" && strings.HasSuffix(line, commentSymbols.EndSymbol) {
			return true
		}
	}
	return false
}

func isBlank(line string) bool {
	return len(line) == 0 || line[0] == ' '
}

func isSupportedLanguage(fileExtension string) bool {
	_, ok := supportedLanguages[fileExtension]
	return ok
}

func getFileExtension(fileName string) (string, error) {
	if !strings.Contains(fileName, fileExtensionSeparator) {
		return "", errors.New("The given file is malformed")
	}

	tokens := strings.Split(fileName, fileExtensionSeparator)
	return tokens[tokensArrayFileExtensionIndex], nil
}
