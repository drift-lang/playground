// The Drift Playground
// GPL v3.0 license to bingxio@qq.com
package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

const (
	port = ":8080" // Protocol exposed port
)

// Configuration
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Accept request and callback data
func recv(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // Config
	if err != nil {
		log.Println("recv:", err)
		return
	}
	defer conn.Close()
	for { // Loop request
		mt, message, err := conn.ReadMessage() // Request data
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Println("successfully:", conn.RemoteAddr())
		out, err := eval(string(message)) // Evaluate the source code
		if err != nil {
			log.Println("file:", err)
		}
		err = conn.WriteMessage(mt, []byte(out)) // Callback to app
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

// Create temporary drift program file
func tempFile(code string) (*os.File, string, error) {
	p := "./tmp/" + time.Now().Format("20060102150405") + ".ft"
	os.Mkdir("./tmp", os.ModePerm) // Path
	f, err := os.Create(p)
	if err != nil {
		return nil, "", err
	}
	return f, p, nil
}

// Write source code to temporary file
func writeFile(code string) (string, error) {
	f, p, err := tempFile(code)
	if err != nil {
		return "", err
	}
	_, err = f.WriteString(code) // Write
	if err != nil {
		return "", err
	}
	return p, nil
}

// Use the drift command to evaluate
func eval(code string) (string, error) {
	p, err := writeFile(code)
	if err != nil {
		return "", err
	}
	cmd := exec.Command("drift", p)  // Command
	out, err := cmd.CombinedOutput() // Stdout and stderr
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func main() {
	http.HandleFunc("/", recv) // An interface
	log.Println("Listen and serve websocket on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
