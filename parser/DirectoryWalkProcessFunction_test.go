package parser

import (
	"parser/types"
	"parser/constants"
	"strings"
	"testing"
)

const blankLine string = " "
const commentLineContent string = "This is a comment line"
const sourceCodeLine string = "int x = 12;"
const testSourceCodeFilePath string = "C:\\users\\TestData\\folder1\\folder2\\File.c"
const testTextFilePath string = "C:\\users\\TestData\\folder1\\folder2\\shakespeare.txt"
const testFileName string = "File"
const testSupportedLanguageFileExtension string = "py"
const dummyProgrammingLanguageFileExtension string = "test"

/*
	Due to the way the program is structured, these tests are fairly inefficient.
	For example, testing if a line is a comment requires looping over all the supported
	language.
	Because of that, the running time complexity for the majority of these tests is something like O(n^4)
	which will limit the amount of language we can support and still have a good test suite.
*/

func TestBlankLineIsDetected(t *testing.T) {
	if !isBlank(blankLine) {
		t.Fail()
	}
}

func TestNotBlankLikeIsDetected(t *testing.T) {
	if isBlank(commentLineContent) {
		t.Fail()
	}
}

func TestFileNameIsExtracted(t *testing.T) {
	for _, supportedLanguage := range supportedLanguages {
		sourceFileExtension := supportedLanguage.SourceFileExtensionName
		testFileNameWithExtension := testFileName + "." + sourceFileExtension

		fileExtension, err := getFileExtension(testFileNameWithExtension)
		if err != nil || strings.Compare(fileExtension, sourceFileExtension) != 0 {
			t.Fail()
		}
	}
}

func TestStringContainingStartCommentIsDetected(t *testing.T) {
	for _, supportedLanguage := range supportedLanguages {
		commentSymbols := supportedLanguage.CommentSymbols
		for _, commentSymbol := range commentSymbols {
			testLineWithComment := commentSymbol.StartSymbol + commentLineContent
			if !isStartComment(testLineWithComment, commentSymbols) {
				t.Fail()
			}
		}
	}
}

func TestStringContainingEndCommentIsDetected(t *testing.T) {
	for _, supportedLanguage := range supportedLanguages {
		commentSymbols := supportedLanguage.CommentSymbols
		for _, commentSymbol := range commentSymbols {
			if commentSymbol.EndSymbol != "" {
				testLineWithComment := commentLineContent + commentSymbol.EndSymbol
				if !isEndComment(testLineWithComment, commentSymbols) {
					t.Fail()
				}
			}
		}
	}
}

func TestLineTypeIsCommentWhenInsideMultilineComment(t *testing.T) {
	for _, commentSymbols := range supportedLanguages {
		if getLineType(commentLineContent, commentSymbols.CommentSymbols, true) != constants.COMMENT {
			t.Fail()
		}
	}
}

func TestSourceCodeLineAreDetected(t *testing.T) {
	for _, commentSymbols := range supportedLanguages {
		if getLineType(sourceCodeLine, commentSymbols.CommentSymbols, false) != constants.SOURCE_CODE {
			t.Fail()
		}
	}
}

func TestSourceCodeLineInsideCommentsAreDetected(t *testing.T) {
	for _, commentSymbols := range supportedLanguages {
		if getLineType(sourceCodeLine, commentSymbols.CommentSymbols, true) != constants.COMMENT {
			t.Fail()
		}
	}
}

func TestLineTypeRoutineChecksForBlankLine(t *testing.T) {
	if getLineType(blankLine, supportedLanguages[testSupportedLanguageFileExtension].CommentSymbols, false) != constants.BLANK {
		t.Fail()
	}
}

func TestLineTypeRoutineChecksForComments(t *testing.T) {
	commentSymbols := supportedLanguages[testSupportedLanguageFileExtension].CommentSymbols
	lineWithComment := commentSymbols[0].StartSymbol + commentLineContent
	if getLineType(lineWithComment, commentSymbols, false) != constants.COMMENT {
		t.Fail()
	}
}

func TestSupportedLangugeIsDetected(t *testing.T) {
	if !isSupportedLanguage(testSupportedLanguageFileExtension) {
		t.Fail()
	}

	if isSupportedLanguage(dummyProgrammingLanguageFileExtension) {
		t.Fail()
	}
}

func TestMultiLineCommentsAreDetected(t *testing.T) {
	for _, supportedLanguage := range supportedLanguages {
		for _, commentSymbols := range supportedLanguage.CommentSymbols {
			lineWithComment := commentSymbols.StartSymbol + commentLineContent
			canSpan := commentSymbols.CanSpanMultipleLines

			if canSpanMultipleLines(lineWithComment, supportedLanguage.CommentSymbols) != canSpan {
				t.Fail()
			}
		}
	}
}

func TestFileAnalysisResultsAreUpdatedCorrectly(t *testing.T) {
	var fileAnalysisResults types.FileAnalysisResult
	updateFileAnalysisResults(&fileAnalysisResults, constants.BLANK, false)
	updateFileAnalysisResults(&fileAnalysisResults, constants.COMMENT, false)
	updateFileAnalysisResults(&fileAnalysisResults, constants.SOURCE_CODE, false)
	updateFileAnalysisResults(&fileAnalysisResults, constants.SOURCE_CODE, true)

	if fileAnalysisResults.BlankLinesCount != 1 ||
		fileAnalysisResults.CommentLinesCount != 1 ||
		fileAnalysisResults.SourceCodeLinesCount != 1 {
		t.Fail()
	}
}
