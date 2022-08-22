package main

import (
	"bufio"
	"fmt"
	"github.com/vova616/screenshot"
	"image/png"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	controllerPath = "/Users/khalil/IdeaProjects/MX-26049/MM/app/Controller/MT/LessonsController.php"
	apiURL         = "https://khalil.testmoby.com/MM/MM/lessons/"
	fkINTeacherID  = "3299943"
	fkINStudentID  = "64684033"
)

func main() {

	file, err := os.Open(controllerPath)
	if err != nil {
		log.Fatal(err)
	}
	counter := 0
	line := ""
	url := ""
	functionName := ""
	argArray := []string{}
	varName := ""
	varType := ""
	defaultValue := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "function MT_") {
			counter++
			functionName = line[strings.Index(line, "MT_")+3 : strings.Index(line, "(")]
			url = apiURL + functionName + "/"
			args := strings.Split(line[strings.Index(line, "(")+1:strings.Index(line, ")")], ",")

			for _, arg := range args {
				arg = strings.TrimSpace(arg)
				if strings.Contains(arg, "=") {
					argArray = strings.Split(arg, " ")
					if len(argArray) == 4 {
						defaultValue = argArray[3]
						argArray = argArray[:2]
					} else if len(argArray) == 3 {
						defaultValue = argArray[2]
						argArray = argArray[:1]
					}
				} else {
					argArray = strings.Split(arg, " ")
					defaultValue = ""
				}
				if len(argArray) == 1 {
					varName = argArray[0]
					varType = ""
				} else if len(argArray) == 2 {
					varType = argArray[0]
					varName = argArray[1]
				}

				if varName == "$fkINTeacherID" {
					url += fkINTeacherID + "/"
				} else if varName == "$fkINStudentID" {
					url += fkINStudentID + "/"
				} else if defaultValue != "" {
					defaultValue = strings.ReplaceAll(defaultValue, `"`, "")
					defaultValue = strings.ReplaceAll(defaultValue, `'`, "")
					defaultValue = strings.ReplaceAll(defaultValue, `null`, "0")
					url += defaultValue + "/"
				} else if varType == "int" {
					url += "0/"
				} else if varType == "string" {
					url += "0/"
				} else if varType == "?string" {
					url += "0/"
				}
			}

			fmt.Println(counter, url)
			openURL(url)
			time.Sleep(5 * time.Second)
			captureScreen(url)
			time.Sleep(1 * time.Second)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func captureScreen(url string) {
	if _, err := os.Stat("./testImages"); os.IsNotExist(err) {
		os.Mkdir("./testImages", os.FileMode(0777))
	}
	img, err := screenshot.CaptureScreen()
	if err != nil {
		panic(err)
	}
	f, err := os.Create("./testImages/" + strings.ReplaceAll(url, "/", ".") + ".png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}

func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
