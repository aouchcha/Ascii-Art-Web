package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func AscciWeb(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("test.html"))
	tmpl.Execute(w, nil)
	switch r.URL.Path {
	case "/":
		fmt.Fprint(w, Result)
	default:
		fmt.Fprint(w, "<h1>Error 404 !</h>")
	}
}

var Result string

func main() {
	if len(os.Args[1:]) == 2 || len(os.Args[1:]) == 1 {
		if os.Args[1] == "" {
			return
		}
		InputFile, todraw := CheckFormat()
		// Traite the filex
		data, err := os.ReadFile(strings.ReplaceAll(InputFile, ".txt", "") + ".txt")
		if err != nil {
			log.Fatalln("open " + InputFile + " no such file or directory")
		}
		if len(data) == 0 {
			log.Fatalln("The file is empty")
		}
		slice := RemoveEmptyStrings(strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n"))
		if len(slice) != 760 {
			log.Fatalln("You have changed the Input file !!!!")
		}
		Result = ascii(slice, todraw)
		if IsAllNewLines(Result) {
			Result = Result[1:]
		}
		fmt.Print(Result)
		Result = strings.ReplaceAll(Result, "\n", "<br>")
		Result = strings.ReplaceAll(Result, " ", "&nbsp;")
	} else {
		log.Fatalln("Usage: go run . [STRING] [BANNER] \nEX: go run . \"something\"")
	}
	http.HandleFunc("/", AscciWeb)
	log.Fatal(http.ListenAndServe(":1303", nil))
}

func ascii(template []string, todraw string) string {
	var Result string
	slicedArgs := strings.Split(todraw, `\n`)
	for _, word := range slicedArgs {
		if word != "" {
			for i := 0; i < 8; i++ {
				for _, char := range word {
					if char < 32 || char > 126 {
						log.Fatalln("You did entered an inprintabale character !!!")
					} else {
						start := int(char-32)*8 + i

						Result += template[start]

					}
				}
				Result += "\n"
			}
		} else {
			Result += "\n"
		}
	}
	return Result
}

func CheckFormat() (string, string) {
	InputFile, todraw := "standard.txt", os.Args[1]
	if len(os.Args[1:]) == 2 {

		InputFile, todraw = os.Args[1], os.Args[2]
		if IsBanner(os.Args[2]) {
			InputFile, todraw = os.Args[2], os.Args[1]
		}
	}
	return InputFile, todraw
}

func RemoveEmptyStrings(slice []string) []string {
	var temp []string
	for i := range slice {
		if slice[i] != "" {
			temp = append(temp, slice[i])
		}
	}
	return temp
}

func IsAllNewLines(str string) bool {
	for _, char := range str {
		if char != '\n' {
			return false
		}
	}
	return true
}

func IsBanner(s string) bool {
	slice := []string{"standard", "shadow", "thinkertoy", "standard.txt", "shadow.txt", "thinkertoy.txt"}
	for i := 0; i < len(slice); i++ {
		if strings.HasSuffix(s, slice[i]) {
			return true
		}
	}
	return false
}
