package main

import (
	"os"
	"fmt"
	"io"
	"flag"
	"io/ioutil"
	"path/filepath"
	"bufio"
)

var (
	flagPackage     = flag.String("package", "main", "Package name where resource will be used")
	flagDestination = flag.String("dest", "data.go", "Path to destination file")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\nOne non-flag argument: path to source file or directory. Also flags supported:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func convertToFile(srcPath string, destFile *os.File) error {
	const BytesPerLine = 16

	if _, err := fmt.Fprintf(destFile, "\"%s\": {\n\t", srcPath); err != nil {
		return err
	}
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	br := bufio.NewReader(srcFile)
	bytesOnLine := 0
	for {
		b, err := br.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		if bytesOnLine == BytesPerLine {
			bytesOnLine = 0
			if _, err = fmt.Fprint(destFile, "\n\t"); err != nil {
				return err
			}
		}
		if _, err = fmt.Fprintf(destFile, "0x%02x,", b); err != nil {
			return err
		}
		bytesOnLine++
	}

	if _, err = fmt.Fprint(destFile, "\n},\n"); err != nil {
		return err
	}
	return nil
}

func main() {
	packageName := *flagPackage
	destPath := *flagDestination

	fmt.Printf("Package name: %s\n", packageName)
	fmt.Printf("Destination: %s\n", destPath)

	srcPath := flag.Arg(0)
	if srcPath == "" {
		fmt.Println("Source path is not defined!")
		os.Exit(1)
		return
	}

	fmt.Printf("Source: %s\n", srcPath)

	stat, err := os.Stat(srcPath)
	if err != nil {
		fmt.Printf("Source data is not available: %s\n", err)
		os.Exit(1)
		return
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		fmt.Printf("Unable to create file: %s", err)
		os.Exit(1)
		return
	}
	defer destFile.Close()

	_, err = fmt.Fprintf(destFile, "package %s\n\nvar filesByteData map[string][]byte = {\n", packageName)
	if err != nil {
		fmt.Printf("Unable to write to file: %s", err)
		os.Exit(1)
		return
	}

	if stat.IsDir() {
		fmt.Println("Source is directory. Walk over...")
		files, err := ioutil.ReadDir(srcPath)
		if err != nil {
			fmt.Printf("Can't list directory: %s", err)
			os.Exit(1)
			return
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			err = convertToFile(filepath.Join(srcPath, file.Name()), destFile)
			if err != nil {
				fmt.Printf("Can't convert file %s: %s", file.Name(), err)
				os.Exit(1)
				return
			}
		}
	} else {
		err = convertToFile(srcPath, destFile)
		if err != nil {
			fmt.Printf("Can't convert file %s: %s", srcPath, err)
			os.Exit(1)
			return
		}
	}

	if _, err = fmt.Fprint(destFile, "\n}\n"); err != nil {
		fmt.Printf("Can't convert file %s: %s", srcPath, err)
		os.Exit(1)
		return
	}
	fmt.Println("Success!")
}
