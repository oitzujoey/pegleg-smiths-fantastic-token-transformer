package main

import "fmt"
import "bufio"
import "os"
import "strings"
import "io"
import "encoding/json"
import "io/fs"
import "path/filepath"


type buildConfiguration struct {
	fileSystem fs.FS
	files []string
	macros []string
	machine string
}


func exit(returnValue int, message string) {
	fmt.Println(message)
	os.Exit(returnValue)
}


func main() {
	var buildConfiguration buildConfiguration
	consoleReader := bufio.NewReader(os.Stdin)
	const buildjson_file_name = "build.json"
	var error error

	// Read build configuration.

	workingDirectory, error := os.Getwd()
	if error != nil {
		exit(1, "Failed to get working directory.")
	}
	var buildjson_directory_path string
	if (len(os.Args) >= 2) {
		buildjson_directory_path = filepath.Join(workingDirectory, os.Args[1])
	} else {
		buildjson_directory_path = workingDirectory
		if error != nil {
			exit(1, "Failed to get working directory.")
		}
	}
	fileSystem := os.DirFS(buildjson_directory_path)
	build_json, error := fs.ReadFile(fileSystem, buildjson_file_name)
	if error != nil {
		exit(1, "Could not read build.json.")
	}
	buildConfiguration.fileSystem = fileSystem
	var build_configuration_map map[string]interface{}
	error = json.Unmarshal(build_json, &build_configuration_map)
	if error != nil {
		panic(error)
	}
	files_json, ok := build_configuration_map["source files"].([]interface{})
	if !ok {exit(1, "No \"source files\" entry found in build.json")}
	buildConfiguration.files = make([]string, len(files_json))
	for index := range len(files_json) {
		buildConfiguration.files[index], ok = files_json[index].(string)
		if !ok {exit(1, "\"source files\" should be an array of strings.")}
	}
	macros_json, ok := build_configuration_map["macros"].([]interface{})
	if !ok {exit(1, "No \"macros\" entry found in build.json")}
	buildConfiguration.macros = make([]string, len(macros_json))
	for index := range len(macros_json) {
		buildConfiguration.macros[index], ok = macros_json[index].(string)
		if !ok {exit(1, "\"macros\" should be an array of strings.")}
	}
	machine_json, ok := build_configuration_map["machine"].(string)
	if !ok {exit(1, "No \"machine\" entry found in build.json")}
	buildConfiguration.machine = machine_json

	if len(os.Args) >= 3 {
		for _, arg := range os.Args[2:] {
			if command_run(arg, &buildConfiguration) {break}
		}
	} else {
		// Start REPL

		quit := false
		for !quit {
			fmt.Print("> ")
			text, e := consoleReader.ReadString('\n')
			if e == io.EOF {
				quit = true
			} else if e != nil {
				fmt.Printf("Error: %v", e);
				return
			}
			quit = command_run(strings.TrimSuffix(text, "\n"), &buildConfiguration)
		}
	}
}

func command_run(command string, buildConfiguration *buildConfiguration) bool {
	quit := false
	if command == "quit" || command == "exit" {quit = true}
	if command == "parse" {
		parse(buildConfiguration)
	} else {
		fmt.Println(command)
	}
	return quit
}
