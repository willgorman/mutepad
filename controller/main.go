package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/term"
	"github.com/vrecan/death/v3"
)

var mutex = &sync.Mutex{}

func main() {
	term, err := term.Open("/dev/cu.usbmodem141203", term.Speed(115200))
	if err != nil {
		panic(err)
	}
	err = term.SetReadTimeout(1 * time.Second)
	if err != nil {
		panic(err)
	}

	// scan the input for a button press
	go func() {
		for {
			b := make([]byte, 6)
			// log.Println("reading")
			_, err = term.Read(b)
			// log.Println("read")
			if errors.Is(err, io.EOF) {
				continue
			}
			if err != nil {
				panic(err)
			}
			mutex.Lock()
			if isMuted() {
				_ = unmute()
				_, err = term.Write([]byte("on"))
				if err != nil {
					panic(err)
				}
			} else {
				_ = mute()
				_, err = term.Write([]byte("off"))
				if err != nil {
					panic(err)
				}
			}
			mutex.Unlock()
			fmt.Println(string(b))
		}
	}()

	// send the current state of the mic
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			msg := "on"
			mutex.Lock()
			if isMuted() {
				msg = "off"
			}
			log.Println(msg)
			_, err = term.Write([]byte(msg))
			if err != nil {
				panic(err)
			}
			mutex.Unlock()
		}
	}()

	log.Println("waiting")
	reaper := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	err = reaper.WaitForDeath()
	log.Println("death is here")
	if err != nil {
		panic(err)
	}
	log.Println("mutex")
	mutex.Lock() // check the lock to make sure we don't halt mid write
	log.Println("close")
	err = term.Close()
	if err != nil {
		panic(err)
	}
	log.Println("bye!")
	mutex.Unlock()
}

func isMuted() bool {
	output, err := exec.Command("osascript", "-e", "input volume of (get volume settings)").Output()
	if err != nil {
		log.Println(err)
	}
	log.Printf("level " + string(output) + "\n")
	log.Printf("muted %t\n", string(output) == "0")
	return strings.TrimSpace(string(output)) == "0"
}

func mute() error {
	return setLevel(0)
}

// TODO: (willgorman) might want to track the previous level and
// restore instead of setting to 100
func unmute() error {
	return setLevel(100)
}

func setLevel(level int) error {
	if level < 0 {
		level = 0
	}
	if level > 100 {
		level = 100
	}
	log.Printf("setting level %d\n", level)
	command := fmt.Sprintf("set volume input volume %d", level)
	_, err := exec.Command("osascript", "-e", command).Output()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
