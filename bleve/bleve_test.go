package bleve

import (
	"fmt"
	"strings"
	"testing"

	"github.com/blevesearch/bleve"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	// Enter the path to your audit.bleve db that you want to populate with dummy data
	workingDir = "/home/bsinou/tmp/test/bleve/2018-11-22"
	// workingDir = ""
	// Use "en" or "fr" to generate randome messages in English or French
	lang = "en"
	// Raise this if you want more events: +1 adds 1k new events, but it takes longer to index.
	loopNb = 20
	// set this flag to true to reduce the size of the sample data set
	useSmallSample = false
)

func TestMain(m *testing.M) {

	if workingDir == "" {
		fmt.Println("\n#### No path to working directory is defined for the various bleve server instances, aborting tests...")
		return
	}

	m.Run()
}

type dummyMsg struct {
	Name string
	Msg  string
}

// TestFillBleveDB simply fills a bleve index with some dummy logs.
// Control size of the test corpus with the 2 global parameters defined above: loopNb & useSmallSample.
// NOTE: you probably want to run this test witout timeout with bigger corpus (typically if loopNb >= 20)
// <code>go test -v -timeout 0 --run TestFillBleveDB</code>
func TestFillBleveDB(t *testing.T) {

	fmt.Printf("\n#### Initializing environment at %s\n", workingDir)

	Convey(fmt.Sprintf("Given a test corpus of %dk log messages", loopNb), t, func() {
		lines := GenerateDummyData(lang, loopNb, useSmallSample)

		Convey("Index in the default store\n", func() {
			var err error
			sname := "defaultStore"
			path := fmt.Sprintf("%s/test-%dk-%s.bleve", workingDir, loopNb, sname)
			server, err := NewDefaultServer(path, sname)
			if err != nil {
				panic("Failed to create or open server at " + path)
			}
			fullPut(server, sname, lines)
		})

		Convey("Index in a default store with mapping\n", func() {
			var err error
			sname := "defaultStoreWithMapping"
			path := fmt.Sprintf("%s/test-%dk-%s.bleve", workingDir, loopNb, sname)
			server, err := NewDefaultServer(path, sname)
			if err != nil {
				panic("Failed to create or open server at " + path)
			}
			fullPut(server, sname, lines)
		})

		Convey("Index in a bolt store\n", func() {
			var err error
			sname := "boltStore"
			path := fmt.Sprintf("%s/test-%dk-%s.bleve", workingDir, loopNb, sname)
			server, err := NewDefaultServer(path, sname)
			if err != nil {
				panic("Failed to create or open server at " + path)
			}
			fullPut(server, sname, lines)
		})

		Convey("Index in a scorch store\n", func() {
			var err error
			sname := "scorchStore"
			path := fmt.Sprintf("%s/test-%dk-%s.bleve", workingDir, loopNb, sname)
			server, err := NewDefaultServer(path, sname)
			if err != nil {
				panic("Failed to create or open server at " + path)
			}
			fullPut(server, sname, lines)
		})
	})
}

func fullPut(index *TestServer, sname string, lines []string) {
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		err := index.PutLog(Json2map(line))
		if err != nil {
			fmt.Printf("[ERROR] failed to put line[%d] - %s : %s \n", i, line, err.Error())
		}
		if (i+1)%100 == 0 {
			fmt.Printf("Indexing events in %s, %d done\n", sname, i+1)
		}
	}
	fmt.Printf("Index done in %s, added %d messages\n", sname, len(lines))
}

// TestSimpleBleve runs a canonical basic test on a new index.
func TestSimpleBleve(t *testing.T) {

	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(workingDir+"/example.bleve", mapping)
	if err != nil {
		fmt.Println(err)
		return
	}

	msgs := []dummyMsg{
		dummyMsg{
			Name: "john",
			Msg:  "is a guy",
		},
		dummyMsg{
			Name: "lily",
			Msg:  "is much more clever",
		},
	}
	// index some data
	index.Index("first", msgs[0])
	index.Index("second", msgs[1])

	// search for some text
	query := bleve.NewMatchQuery("john")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)
}
