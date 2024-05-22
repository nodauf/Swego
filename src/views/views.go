package views

import "embed"

//go:embed *.tpl
var views embed.FS

func GetViews(filename string) ([]byte, error) {
	return views.ReadFile(filename)
}
