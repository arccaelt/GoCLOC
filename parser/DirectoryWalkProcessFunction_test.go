package parser

import "testing"
import "strings"

const blankLine string = " "
const commentLineContent string = "This is a comment line"
const testSourceCodeFilePath string = "C:\\users\\TestData\\folder1\\folder2\\File.c"
const testTextFilePath string = "C:\\users\\TestData\\folder1\\folder2\\shakespeare.txt"
const testFileName string = "File"

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
	// for _, supportedLanguage := range supportedLanguages {
	// 	commentSymbols := supportedLanguage.CommentSymbols
	// }
}

func TestStringContainingACommentIsDetected(t *testing.T) {

}