package main

import "flag"
import "fmt"
import "io"
import "os"
import "path/filepath"
import "regexp"
import "strconv"

func main() {
	dryRun := flag.Bool("dry-run", false, "Dry run, don't actually copy files.")
	verbose := flag.Bool("verbose", false, "Verbose output, print list of files copied.")
	fileCount := 0

	flag.Parse()

	regex := flag.Arg(0)
	src := flag.Arg(1)
	dest := flag.Arg(2)

	if len(regex) == 0 {
		fmt.Println("Invalid regex param")
		os.Exit(1)
	}

	src_fi, src_err := os.Stat(src)
	if len(src) == 0 || src_err != nil || !src_fi.IsDir() {
		fmt.Println("Reggie: Invalid source directory.")
		os.Exit(1)
	}

	dest_fi, dest_err := os.Stat(dest)
	if len(dest) == 0 || dest_err != nil || !dest_fi.IsDir() {
		fmt.Println("Reggie: Invalid destination directory.")
		os.Exit(1)
	}

	callback := func(path string, fi os.FileInfo, err error) error {
		match, err := regexp.MatchString(regex, path)
		if match {
			if *verbose {
				fmt.Println("Copying: " + path)
			}
			if !*dryRun {
				err := CopyFile(path, dest)
				if err != nil {
					fmt.Println(err)
				}
			}
			fileCount++
		}
		return err
	}

	filepath.Walk(src, callback)
	fmt.Printf("Copied %v files\n", fileCount)
	if *dryRun {
		fmt.Println("Dry run complete")
	}
}

func FileExists(fpath string) bool {
	exists := false
	if _, err := os.Stat(fpath); err == nil {
		exists = true
	}
	return exists
}

func NewFileName(fpath string) string {
	idx := 1
	for FileExists(fpath) {
		dir, file := filepath.Split(fpath)

		//TODO: Try and consolidate these as one regex replacement
		match, err := regexp.MatchString(`\-[\d]+\.[\w]{1,3}$`, file)

		if match && err == nil {
			re := regexp.MustCompile(`\-(\d)+(\.[\w]{1,3})$`)
			file = re.ReplaceAllString(file, "-"+strconv.Itoa(idx)+"$2")
		} else {
			re := regexp.MustCompile(`(\.[\w]{1,3})$`)
			file = re.ReplaceAllString(file, "-"+strconv.Itoa(idx)+"$1")
		}
		fpath = dir + "/" + file
		idx++
	}
	return fpath
}

func CopyFile(csrc string, cdest string) error {
	s, err := os.Open(csrc)
	if err != nil {
		return err
	}

	fd_base := filepath.Base(csrc)
	fd_dest := cdest + "/" + fd_base

	if FileExists(fd_dest) {
		fd_dest = NewFileName(fd_dest)
	}

	d, err := os.Create(fd_dest)
	if err != nil {
		return err
	}

	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
