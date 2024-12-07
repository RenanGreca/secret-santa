package file

import (
	"os"

	"secret-santa/assert"
)

// ReadFile reads a file and returns its contents as a string.
func ReadFile(fileName string) string {
	b, err := os.ReadFile(fileName)
	assert.NoError(err, "File: Error reading file")

	return string(b)
}

// FileExists checks if a file exists in the filesystem.
func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		return false
	}
	return true
}

func Delete(fileName string) {
	if FileExists(fileName) {
		err := os.Remove(fileName)
		assert.NoError(err, "File: Error deleting file")
	}
}

func Create(fileName string, contents string) *os.File {
	if FileExists(fileName) {
		// f, err := os.OpenFile(fileName, os.O_RDWR, 0644)
		// assert.NoError(err, "File: Error opening file")
		// return f
		Delete(fileName)
	}
	f, err := os.Create(fileName)
	assert.NoError(err, "File: Error creating file")
	err = os.WriteFile(fileName, []byte(contents), 0644)
	return f
}
