package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var filesPath []string = parseArgs()
var properties, keys = getProperties(filesPath)
func parseArgs() []string {
	if len(os.Args) < 2 {
		log.Fatalln("Not enough arguments: \nUsage: go run main.go <file path>...")
	}
	return os.Args[1:]
}
func main() {
	app := app.New()

	// Create a new window and set its default size
	w := app.NewWindow("lke - .ltx Properties Viewer")
	w.Resize(fyne.NewSize(600, 400))


	container := container.NewVBox()
	properties, keys := getProperties(filesPath)
	for _, key := range keys {
		value := properties[key]
		entry := widget.NewEntry()
		entry.SetText(value)
		container.Add(widget.NewForm(&widget.FormItem{
			Text:   key,
			Widget: entry,
		}))
	
		saveButton := widget.NewButton("Save", func() {
			properties[key] = entry.Text
			found := false

			for _, k := range keys {
				if k == key {
					found = true
					break
				}
			}

			if !found {
				keys = append(keys, key)
			}
    		// Save changes to file
			ltxFilePath := os.Args[1:]
    		// Overwrite file
			file, err := os.OpenFile(ltxFilePath[0], os.O_WRONLY|os.O_TRUNC, 0644)
    		if err != nil {
    		    log.Fatal(err)
    		}
    		defer file.Close()
    		for _, key := range keys {
    		    _, err := file.WriteString(key + "=" + properties[key] + "\n")
				println(key + "=" + properties[key] + "\n")
    		    if err != nil {
    		        log.Fatal(err)
    		    }
    		}
		})
	
		undoButton := widget.NewButton("Undo", func() {
			entry.SetText(properties[key])
		})
		
		container.Add(saveButton)
		container.Add(undoButton)

		container.Add(widget.NewSeparator())
		
	}

	// footer
container.Add(widget.NewLabel("ltx key-value editor by @editor99 (github.com/gteditor99/ltx)"))
w.SetContent(container)
w.ShowAndRun()
	
}


// Modify getProperties to return a slice of keys in addition to the map
func getProperties(filesPath []string) (map[string]string, []string) {
	properties := make(map[string]string)
	var keys []string
	for _, filePath := range filesPath {
		file, err := os.Open(filePath)
		if err != nil {
		   log.Fatal(err)
		}
		defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
			line := scanner.Text()
			split := strings.SplitN(line, "=", 2)
			if len(split) == 2 {
				properties[split[0]] = split[1]
				keys = append(keys, split[0])
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return properties, keys
}