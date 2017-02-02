package main

import (
  "io/ioutil"
  "bytes"
  "fmt"
  "log"
  "net/http"
  "os"
  "os/exec"
)

func determineListenAddress() (string, error) {
  port := os.Getenv("PORT")
  if port == "" {
    return "", fmt.Errorf("$PORT not set")
  }
  return ":" + port, nil
}

func getAsm(w http.ResponseWriter, r *http.Request) {
  // write asm to temp file
  content := []byte(`package main

import (
  "fmt"
)

func main() {
  fmt.Println("Hello, world")
}
`)
  tmpfile, err := ioutil.TempFile("", "goasm-")
  if err != nil {
    fmt.Fprintln(w, "Error writing to file")
  }
  defer os.Remove(tmpfile.Name()) // clean up
  if _, err := tmpfile.Write(content); err != nil {
    log.Fatal(err)
  }
  if err := tmpfile.Close(); err != nil {
    log.Fatal(err)
  }

  // compile
  cmd := exec.Command("go", "tool", "compile", "-o", "/dev/null", "-S",
                      tmpfile.Name())
  var out bytes.Buffer
  cmd.Stdout = &out
  err = cmd.Run()
  if err != nil {
    fmt.Fprintln(w, "Error compiling")
  }
  fmt.Fprintln(w, out.String())
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
