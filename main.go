package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func readDate(dirname string, filename os.FileInfo) time.Time {
	fpath := path.Join(dirname, filename.Name())
	f, err := os.Open(filepath.Clean(fpath))

	if err != nil {
		panic(err)
	}

	x, err := exif.Decode(f)
	if err != nil {
		panic(err)
	}

	tm, _ := x.DateTime()

	_ = f.Close()

	return tm
}

func formatDate(d time.Time) string {
	return d.Format("2006-01-02T15:04:05")
}

type by []os.FileInfo

func (a by) Len() int           { return len(a) }
func (a by) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a by) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func main() {
	miniumArguments := 2
	args := os.Args

	if len(args) < miniumArguments {
		panic("You have to specify a filename")
	}

	home := os.Getenv("HOME")
	destDirname := home + "/txt/photos/"

	dirname := filepath.Clean(args[1])

	files, _ := ioutil.ReadDir(dirname)
	for _, file := range files {
		dt := readDate(dirname, file)
		to := strings.Split(formatDate(dt), "+")[0]
		to += ".jpeg"

		_ = os.Rename(path.Join(dirname, file.Name()), path.Join(dirname, to))
	}

	files, _ = ioutil.ReadDir(dirname)
	sort.Sort(by(files))

	data := make(map[string][]os.FileInfo)

	for _, file := range files {
		dt := readDate(dirname, file)
		yr := dt.Format("2006-01-02")

		if data[yr] == nil {
			data[yr] = make([]os.FileInfo, 0)
		}

		data[yr] = append(data[yr], file)
	}

	for year, data := range data {
		fmt.Printf("* %s\n", year)
		fmt.Printf("  :PROPERTIES:\n")
		fmt.Printf("  :CREATED: %s\n", strings.ReplaceAll(year, "-", ""))
		fmt.Printf("  :END:\n")

		for _, file := range data {
			dt := readDate(dirname, file)
			tt := dt.Format("15:04:05")
			fmt.Printf("** %s\n", tt)
			fmt.Printf("[[../../photos/%s]]\n", file.Name())
		}
	}

	for _, file := range files {
		_ = os.Rename(path.Join(dirname, file.Name()), path.Join(destDirname, file.Name()))
	}
}
