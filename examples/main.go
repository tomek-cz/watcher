package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/radovskyb/watcher"
)

func main() {
	w := watcher.New()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case event := <-w.Event:
				// Print the event type.
				fmt.Println(event)

				// Print out the file name with a message
				// based on the event type.
				switch event.EventType {
				case watcher.EventFileModified:
					fmt.Println("Modified file:", event.Name())
				case watcher.EventFileAdded:
					fmt.Println("Added file:", event.Name())
				case watcher.EventFileDeleted:
					fmt.Println("Deleted file:", event.Name())
				}
			case err := <-w.Error:
				log.Fatalln(err)
			}
		}
	}()

	// Watch this file for changes.
	if err := w.Add("main.go"); err != nil {
		log.Fatalln(err)
	}

	// Watch test_folder recursively for changes.
	if err := w.Add("test_folder"); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched.
	for _, f := range w.Files {
		fmt.Println(f.Name())
	}

	// Trigger an event after 500 milliseconds.
	go func() {
		time.Sleep(time.Millisecond * 500)
		w.Trigger(watcher.EventFileAdded, nil)
	}()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(100); err != nil {
		log.Fatalln(err)
	}

	wg.Wait()
}
