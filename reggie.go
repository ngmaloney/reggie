package main

import "os"
import "fmt"
import "path/filepath"
import "regexp"
import "io"
import "flag"

func main() {
	regex := os.Args[1]
	src := os.Args[2]
	dest := os.Args[3]
	dryRun := flag.Bool("dry-run", false, "Dry run, don't copy files.")

	flag.Parse()

	if(len(regex) == 0) {
		fmt.Println("Valid regex argument is required")
		os.Exit(1)
	}

  src_fi, src_err := os.Stat(src)
	if (len(src) == 0 || src_err != nil || !src_fi.IsDir()) {
		fmt.Println("Valid source directory is required")
		os.Exit(1)
  }

  dest_fi, dest_err := os.Stat(dest)
	if (len(dest) == 0 || dest_err != nil || !dest_fi.IsDir()) {
		fmt.Println("Valid destination directory is required")
		os.Exit(1)
	}

	callback := func(path string, fi os.FileInfo, err error) error {
		match, err := regexp.MatchString(regex, path)
		if(match) {
			fmt.Println("Copying: " + path)
			if(!*dryRun) {
				err := CopyFile(path, dest)
				if(err != nil) {
					fmt.Println(err)
				}
			}
		}
		return err
  }

	filepath.Walk(src, callback)
	if(*dryRun) {
		fmt.Println("Dry run complete")
	}
}

func CopyFile(csrc string, cdest string) error {
	s, err := os.Open(csrc)
	if err != nil {
		return err
	}

	fd_base := filepath.Base(csrc)
	fd_dest := cdest + "/" + fd_base

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
