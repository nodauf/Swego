package controllers

func EmbeddedFiles() string {
	returnValue := ""
	_, childrenFiles := listEmbeddedFiles()
	for _, value := range childrenFiles {
		returnValue += value + "\n"
	}
	return returnValue
}

// To avoid generating error on Linux
func RunEmbeddedBinary(binary string, arguments string) {

}
