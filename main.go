package main

import (
	"io/ioutil"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
)

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func getAsm(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	gocodes, ok := r.Form["gocode"]
	if !ok {
		fmt.Fprintln(w, "Error: missing go code")
		return
	}
	gocode := []byte(gocodes[0])
	// write gocode to temp file
	tmpfile, err := ioutil.TempFile("", "goasm-")
	if err != nil {
		fmt.Fprintln(w, "Error: couldn't open temp file")
		return
	}
	defer os.Remove(tmpfile.Name()) // clean up
	if _, err := tmpfile.Write(gocode); err != nil {
		fmt.Fprintln(w, "Error: couldn't write temp file")
		return
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Fprintln(w, "Error: couldn't close temp file")
		return
	}
	// compile
	cmd := exec.Command("go", "tool", "compile", "-o", "/dev/null", "-S",
		tmpfile.Name())
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(w, "Error compiling")
		return
	}
	// replace real paths with fake paths
	re := regexp.MustCompile(`\(/.*?:`) 
	fmt.Fprintln(w, re.ReplaceAllLiteralString(out.String(), "(main.go:"))
}

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/asm", getAsm)
	log.Printf("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
