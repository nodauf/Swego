package controllers

// EmbeddedFiles list the embedded files
func EmbeddedFiles() string {
	returnValue := ""
	_, childrenFiles := listEmbeddedFiles()
	for _, value := range childrenFiles {
		returnValue += value + "\n"
	}
	return returnValue
}

// RunEmbeddedBinary Do nothng, only to avoid generating error on Linux
func RunEmbeddedBinary(binary string, arguments string) {

}
