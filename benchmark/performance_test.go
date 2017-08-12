package benchmark

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/crawlcoin/crawlcoin/util"
	"github.com/facebookgo/ensure"
	"github.com/vkrasnov/dictator"
	"gopkg.in/kothar/brotli-go.v0/enc"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

type Stats struct {
	Label               string
	Time                time.Duration
	CompressedSizeBytes int
	start               time.Time
}

func (s *Stats) String(rawSize int) string {
	return s.Label + "\t - took:" + s.Time.String() +
		" compressed size: " + strconv.Itoa(s.CompressedSizeBytes) + " giving compression ratio of: " + util.FormatFloat(float64(s.CompressedSizeBytes)/float64(rawSize))
}

func (s *Stats) Start() {
	s.start = time.Now()
}

func (s *Stats) End() {
	s.Time += time.Since(s.start)
}

func TestBrotliVsZlib(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	files, err := AllFilesIn(usr.HomeDir + "/crawl/http/cnn.com")
	ensure.Nil(t, err)

	zlib1Stats := &Stats{Label: "zlib 1"}
	zlib9Stats := &Stats{Label: "zlib 9"}
	brotliStats := &Stats{Label: "brotli"}

	zlib1DictStats := &Stats{Label: "zlib 1 dict"}
	zlib9DictStats := &Stats{Label: "zlib 9 dict"}
	brotliDictStats := &Stats{Label: "brotli dict"}

	allStats := []*Stats{zlib1Stats, zlib9Stats, brotliStats, zlib1DictStats, zlib9DictStats, brotliDictStats}

	windowSize := 1024 * 16
	dictSize := 1024 * 16

	dictFiles := []string{}
	for _, file := range files {
		b, err := ioutil.ReadFile(file)
		ensure.Nil(t, err)
		if len(b) > 0 {
			dictFiles = append(dictFiles, file)
		}
	}
	table := dictator.GenerateTable(windowSize, dictFiles, 4, nil, runtime.NumCPU())
	dict := []byte(dictator.GenerateDictionary(table, dictSize, int(math.Ceil(float64(len(dictFiles))*0.01))))

	log.Println("Generated dict of length: ", len(dict))

	var totalRaw int

	for _, file := range files {
		b, err := ioutil.ReadFile(file)
		ensure.Nil(t, err)
		totalRaw += len(b)

		zlibSize(b, zlib.BestSpeed, nil, zlib1Stats)
		zlibSize(b, zlib.BestCompression, nil, zlib9Stats)
		brotliSize(b, nil, brotliStats)

		zlibSize(b, zlib.BestSpeed, dict, zlib1DictStats)
		zlibSize(b, zlib.BestCompression, dict, zlib9DictStats)
		brotliSize(b, dict, brotliDictStats)
	}

	for _, s := range allStats {
		fmt.Println(s.String(totalRaw))
	}
}

func zlibSize(b []byte, level int, dict []byte, s *Stats) {
	var zipOut bytes.Buffer
	var err error
	var w *zlib.Writer
	s.Start()
	if dict != nil {
		w, err = zlib.NewWriterLevelDict(&zipOut, level, dict)
	} else {
		w, err = zlib.NewWriterLevel(&zipOut, level)
	}
	if err != nil {
		panic(err)
	}
	w.Write(b)
	w.Close()
	s.End()
	s.CompressedSizeBytes += zipOut.Len()
}

func brotliSize(b []byte, dict []byte, s *Stats) {
	var err error
	var encoded []byte
	s.Start()
	if dict != nil {
		encoded, err = enc.CompressBufferDict(nil, b, dict, nil)
	} else {
		encoded, err = enc.CompressBuffer(nil, b, nil)
	}
	if err != nil {
		panic(err)
	}
	s.End()
	s.CompressedSizeBytes += len(encoded)
}

func AllFilesIn(baseDir string) ([]string, error) {
	ret := []string{}
	fn := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			ret = append(ret, path)
		}
		return nil
	}
	err := filepath.Walk(baseDir, fn)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
