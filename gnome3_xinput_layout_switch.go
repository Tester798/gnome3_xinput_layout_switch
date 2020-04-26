//go:generate env CGO_ENABLED=0 go build -o ./gnome-xinput-layout-switch -ldflags "-s -w"

package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os/exec"
	"strings"
)

type EventData struct {
	Type   string
	Detail string
}

var debug = flag.Bool("debug", false, "Dump keyboard events")
var key1 = flag.String("key1", "37,105", "First key code to monitor (default: 37,105 [Ctrl])")
var key2 = flag.String("key2", "50,62", "Second key code to monitor (default: 50,62 [Shift])")
var key1Set = make(map[string]string)
var key2Set = make(map[string]string)

func dumpEvent(event *EventData) {
	if *debug {
		log.Printf("EVENT: %v\n", event)
	}
}

var key1Pressed = false
var key2Pressed = false
var canSwitch = false

func processEvent(event EventData) {
	if len(event.Type) == 0 {
		return
	}
	if event.Type == "RawKeyPress" {
		dumpEvent(&event)
		if len(key1Set[event.Detail]) > 0 {
			key1Pressed = true
			canSwitch = key2Pressed
			return
		}
		if len(key2Set[event.Detail]) > 0 {
			key2Pressed = true
			canSwitch = key1Pressed
			return
		}
		canSwitch = false
	}
	if event.Type == "RawKeyRelease" {
		dumpEvent(&event)
		doSwitch := false
		if len(key1Set[event.Detail]) > 0 {
			if key2Pressed {
				doSwitch = true
			}
			key1Pressed = false
		}
		if len(key2Set[event.Detail]) > 0 {
			if key1Pressed {
				doSwitch = true
			}
			key2Pressed = false
		}

		if canSwitch && doSwitch {
			if *debug {
				log.Printf("Change layout\n")
			}
			layoutCmd := exec.Command("gdbus",
				"call",
				"--session",
				"--dest", "org.gnome.Shell",
				"--object-path", "/org/gnome/Shell",
				"--method", "org.gnome.Shell.Eval",
				"imports.ui.status.keyboard.getInputSourceManager()._mruSources[1].activate()",
			)
			err := layoutCmd.Run()
			if err != nil {
				log.Printf("Unable to swithc layout:\n%v\n", err)
			} else {
				_ = layoutCmd.Wait()
			}
		}
	}
}
func main() {
	flag.Parse()
	cmd := exec.Command("xinput", "test-xi2", "--root")
	for _, k := range strings.Split(*key1, ",") {
		key1Set[k] = k
	}
	for _, k := range strings.Split(*key2, ",") {
		key2Set[k] = k
	}
	reader, writer := io.Pipe()
	cmd.Stdout = writer
	go func() {

		scanner := bufio.NewScanner(reader)
		event := EventData{}
		for scanner.Scan() {
			eoe := false
			text := strings.TrimSpace(scanner.Text())
			attrs := strings.Split(text, " ")
			if strings.HasPrefix(text, "EVENT ") {
				processEvent(event)
				event = EventData{
					Type: strings.Trim(attrs[3], "()"),
				}
			}
			if strings.HasPrefix(text, "detail:") {
				event.Detail = attrs[1]
			}
			if len(text) == 0 || strings.HasPrefix(attrs[0], "valuators:") {
				eoe = true
			}
			//log.Printf("%d: %v\n", len(text), text) // the line
			if eoe {
				processEvent(event)
				event = EventData{}
			}
		}
	}()
	err := cmd.Run()
	_ = writer.Close()
	_ = reader.Close()
	if err != nil {
		log.Fatalf("Unable to run xinput %v", err)
	}
	_ = cmd.Wait()

}
