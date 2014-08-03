net-go
=======
Queries for files and gives meta-data for further interactions

Example:
--------
```go
package main

import (
	"github.com/nateri/net-go"
	"fmt"
	"log"
	)
	
func main() {
	finder, err := netget.NewVideoFinder()
	if err != nil {
		log.Fatal(err)
	}
	err, result = finder.Find("frozen", 720)
	if err != nil {
		log.Fatal(err)
	}
	for {
		result {}
	}
}

func main() {
	watcher, err := netlisten.NewVideoWatcher()
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Watch("barca", 720)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case ev := <-watcher.Event:
			fmt.Println("event:", ev)
		case err := <-watcher.Error:
			fmt.Println("error:", err)
		}
	}
}
```

There is also an option to skip certain folders (like .git for example):

```go
	
```
