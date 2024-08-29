package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/eiannone/keyboard"
)

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%d:%02d", m, s)
	} else {
		return fmt.Sprintf("%d", s)
	}
}
func clearScreen() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func main() {
	hideCursor()
	defer showCursor()

	// Set up clean exit to ensure cursor is shown again
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		showCursor()
		os.Exit(1)
	}()

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	start := time.Now()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	paused := false
	var pauseTime time.Time
	var elapsed time.Duration

	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if key == keyboard.KeySpace {
				paused = !paused
				if paused {
					pauseTime = time.Now()
				} else {
					start = start.Add(time.Since(pauseTime))
				}
			} else if char == 'q' || char == 'Q' {
				showCursor()
				os.Exit(0)
			}
		}
	}()

	for range ticker.C {
		clearScreen()
		if !paused {
			elapsed = time.Since(start)
		}
		timeStr := formatDuration(elapsed)
		color := "yellow"
		if paused {
			color = "red"
		}
		myFigure := figure.NewColorFigure(timeStr, "basic", color, true)
		myFigure.Print()

		if paused {
			fmt.Println("\nPAUSED - Press SPACE to resume")
		} else {
			fmt.Println("\nRUNNING - Press SPACE to pause")
		}
	}
}
