package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

const (
	defaultScanPath = "/media"
	rsync           = "/usr/local/bin/rsync"
)

var (
	fileExts    = []string{"CR2"}
	targetOwner string
	targetPath  string
)

func main() {
	scanpath := os.Getenv("PHINGESTER_SCANPATH")
	if len(scanpath) == 0 {
		scanpath = defaultScanPath
	}

	targetOwner = os.Getenv("PHINGESTER_OWNER")

	tp := os.Getenv("PHINGESTER_DEST")
	if len(tp) != 0 {
		targetPath = tp
	} else {
		homedir := os.Getenv("HOME")
		if len(homedir) == 0 {
			log.Fatal("neither PHINGESTER_DEST nor HOME are set, please define one of those to define copy destination")
		}
		targetPath = filepath.Join(homedir, "phingester_media")
	}
	os.MkdirAll(targetPath, 0777)

	for {
		scan(scanpath)
		time.Sleep(1 * time.Minute)
	}
}

func scan(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		dir, ok := isPhotoMedium(path, f.Name())
		if !ok {
			continue
		}

		fls := mediaFiles(dir)
		transferFiles(fls)
	}
}

func isPhotoMedium(basepath, dirname string) (string, bool) {
	p := path.Join(basepath, dirname, "DCIM")
	_, err := os.Stat(p)
	if err != nil {
		return "", false
	}
	return p, true
}

func mediaFiles(basepath string) []string {
	var files []string

	filepath.Walk(basepath, func(path string, info os.FileInfo, err error) error {
		for _, ext := range fileExts {
			if "."+ext == filepath.Ext(path) {
				files = append(files, path)
			}
		}
		return nil
	})
	return files
}

func transferFiles(files []string) {
	for _, file := range files {
		tfp := targetFilePath(file)
		if _, err := os.Stat(tfp); err == nil {
			// Target file already exists, skip.
			continue
		}
		cmd := exec.Command(rsync, "-a", file, tfp)
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		if targetOwner != "" {
			exec.Command("chown", targetOwner, tfp).Run()
		}
	}
}

func targetFilePath(origin string) string {
	_, fn := filepath.Split(origin)
	fi, err := os.Stat(origin)
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(targetPath, fi.ModTime().Format("2006-01-02")+"_"+fn)
}
