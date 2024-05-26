package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func FirstPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("test.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func AsciiArtWeb(w http.ResponseWriter, r *http.Request) {
	tmpl , err := template.ParseFiles("ascii-art.html")
	if err != nil {
		log.Fatalln(err)
	}
	text := r.FormValue("string")
	banner := r.FormValue("banner") + ".txt"
	

	if text == "" || banner == "" {
		http.Error(w, "Missing form fields", http.StatusBadRequest)
		return
	}

	slice, slicedArgs := FormatofBanner(banner, text)
	Result := DrawAsciiFS(slice, slicedArgs)

	tmpl.Execute(w,Result)
}

func FormatofBanner(InputFile, text string) ([]string, []string) {
	var sep string
	if InputFile == "standard.txt" || InputFile == "shadow.txt" {
		sep = "\n"
	} else {
		sep = "\r\n"
	}
	// Read the file
	data, err := os.ReadFile(InputFile)
	if err != nil {
		log.Fatalln(err)
	}
	slice := RemoveEmptyStrings(strings.Split(string(data), sep))
	slicedArgs := strings.Split(text, `\n`)
	return slice, slicedArgs
}

func DrawAsciiFS(slice, slicedArgs []string) string {
	Result := ""
	for _, word := range slicedArgs {
		if word != "" {
			for i := 0; i < 8; i++ {
				for _, char := range word {
					if char < 32 || char > 126 {
						log.Fatalln("You did entered an inprintabale character !!!")
					} else {
						start := int(char-32)*8 + i
						Result += slice[start]
					}
				}
				Result += "\n"
			}
		} else {
			Result += "\n"
		}
	}
	if IsAllNewLines(Result) {
		Result = Result[1:]
	}
	return Result
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

func main() {
	http.HandleFunc("/", FirstPage)
	http.HandleFunc("/ascii-art", AsciiArtWeb)
	err := http.ListenAndServe(":1515", nil)
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}
