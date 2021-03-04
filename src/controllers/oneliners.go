package controllers

import (
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"SimpleHTTPServer-golang/src/utils"

	rice "github.com/GeertJohan/go.rice"
	"github.com/manifoldco/promptui"
)

type File struct {
	Name string
	Path string
}

// CLI
func CliOnelinersMenu() {
	sliceFiles := []File{}
	templateBox, err := rice.FindBox("../views/")
	utils.Check(err, "oneliners: error while opening rice box")

	// get file contents as string
	templateOneliners, err := templateBox.String("oneliners.tpl")
	utils.Check(err, "oneliners: error while opening oneliners.tpl in rice box")

	templateOneliners = searchAndReaplceOneliners(templateOneliners)

	// Get all the file's name recursively. Store them in sliceFiles
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				sliceFiles = append(sliceFiles, File{Name: info.Name(), Path: path})
			}
			return nil
		})
	utils.Check(err, "oneliners: error while filepath.walk")

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Name | cyan }} ({{ .Path | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Path | red }})",
		Selected: "\U0001F336 {{ .Name | red | cyan }}",
		Details:  templateOneliners}

	searcher := func(input string, index int) bool {
		filesSearcher := sliceFiles[index]
		name := strings.Replace(strings.ToLower(filesSearcher.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "File",
		Items:     sliceFiles,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	utils.Check(err, "oneliners: prompt failed")
	// Tips to define fake faint function which is a color for promptui
	funcMap := template.FuncMap{
		"faint": func(str string) string { return str },
	}

	template, err := template.New("onelinersTemplate").Funcs(funcMap).Parse(templateOneliners)

	utils.Check(err, "oneliners: Error while using text/template")
	// Execute the template with the chosen file
	template.Execute(os.Stdout, File{Name: sliceFiles[i].Name, Path: sliceFiles[i].Path})
}

func searchAndReaplceOneliners(template string) string {
	// Replace IP and Port by the IP, port, ... put in arguments
	template = strings.ReplaceAll(template, "[IP]", *IP)
	template = strings.ReplaceAll(template, "[PORT]", strconv.Itoa(*Bind))
	if *Tls {
		template = strings.ReplaceAll(template, "[PROTO]", "https")
	} else {
		template = strings.ReplaceAll(template, "[PROTO]", "http")
	}
	return template
}
