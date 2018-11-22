package main

import (
	"fmt"
	"strings"

	ownbleve "github.com/bsinou/gopocs/bleve"
)

var (
	dStore, withMapping, boltStore, scorchStore *ownbleve.TestServer
)

const (
	// Enter the path to your audit.bleve db that you want to populate with dummy data
	workingDir = "/home/bsinou/tmp/test/bleve/2018-11-21"
	// workingDir = ""
	// Raise this if you want more events: +1 adds 1k new events, but it takes longer to index.
	loopNb = 100
	// set this flag to true to reduce the size of the sample data set
	useSmallSample = false
)

func main() {
	if workingDir == "" {
		fmt.Println("[Error] No path to working directory is defined for the various bleve server instances, aborting tests...")
		return
	}

	err := initialisesStores()
	if err != nil {
		fmt.Println("[Error] Could not initialise stores: " + err.Error())
		return
	}

	lines := ownbleve.GenerateDummyData(loopNb, useSmallSample)

	fmt.Printf("### [Info] Launching index on %s with %dk lines.\n", dStore.GetName(), loopNb)
	fullPut(dStore, dStore.GetName(), lines)
	// fmt.Printf("### [Info] Launching index on %s with %dk lines.\n", boltStore.GetName(), loopNb)
	// fullPut(boltStore, boltStore.GetName(), lines)
	fmt.Printf("### [Info] Launching index on %s with %dk lines.\n", scorchStore.GetName(), loopNb)
	fullPut(scorchStore, scorchStore.GetName(), lines)
}

func initialisesStores() error {

	var err error
	sname := "defaultStore"
	path := fmt.Sprintf("%s/test-%dk-%s.bleve", workingDir, loopNb, sname)
	dStore, err = ownbleve.NewDefaultServer(path, sname)
	if err != nil {
		return fmt.Errorf("fail to create %s index, cause: %s", sname, err.Error())
	}

	sname = "boltStore"
	path = fmt.Sprintf("%s/test-%dk-%s.bleve", workingDir, loopNb, sname)
	boltStore, err = ownbleve.NewServerOnBolt(path, sname)
	if err != nil {
		return fmt.Errorf("fail to create %s index, cause: %s", sname, err.Error())
	}

	sname = "scorchStore"
	path = fmt.Sprintf("%s/test-%dk-%s.bleve", workingDir, loopNb, sname)
	scorchStore, err = ownbleve.NewServerOnScorch(path, sname)
	if err != nil {
		return fmt.Errorf("fail to create %s index, cause: %s", sname, err.Error())
	}

	return nil
}

func fullPut(index *ownbleve.TestServer, sname string, lines []string) {
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		err := index.PutLog(ownbleve.Json2map(line))
		if err != nil {
			fmt.Printf("[ERROR] failed to put line[%d] - %s : %s \n", i, line, err.Error())
		}
		if (i+1)%100 == 0 {
			fmt.Printf("Indexing events in %s, %d done\n", sname, i+1)
		}
	}
	fmt.Printf("Index done in %s, added %d messages\n", sname, len(lines))
}
