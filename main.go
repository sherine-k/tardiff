package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	var startingDate string
	var inPath string
	var outPath string

	rootCmd := &cobra.Command{
		Use:   "tardiff",
		Short: "tardiff archives all files modified after a given date",
		Long:  "tardiff archives all files modified after a given date",
		Run: func(cmd *cobra.Command, args []string) {
			root := inPath

			// Parse startingDate into time
			t, err := time.Parse("2006-01-02", startingDate)
			if err != nil {
				fmt.Printf("Error parsing time: %v\n", err)
				return
			}
			fmt.Printf("Parsed time: %v\n", t)

			// Create the tar file
			outFile, err := os.Create(outPath)
			if err != nil {
				fmt.Printf("Error creating tar file: %v\n", err)
				return
			}
			defer outFile.Close()

			tarWriter := tar.NewWriter(outFile)
			defer tarWriter.Close()

			err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					fmt.Printf("Directory: %s\n", path)
				} else if info.Mode().IsRegular() {
					// Check if file was modified after t
					if info.ModTime().After(t) {
						fmt.Printf("File: %s was modified after %v\n", path, t)
						// Add this file to the tar file
						// Open the file
						f, err := os.Open(path)
						if err != nil {
							return fmt.Errorf("error opening file %s: %v", path, err)

						}
						defer f.Close()

						// Get the file info
						info, err := f.Stat()
						if err != nil {
							return fmt.Errorf("error getting file info for %s: %v", path, err)

						}

						// Create a new header for the file
						header := &tar.Header{
							Name:    path,
							Size:    info.Size(),
							Mode:    int64(info.Mode()),
							ModTime: info.ModTime(),
						}

						// Write the header to the tar file
						if err := tarWriter.WriteHeader(header); err != nil {
							return fmt.Errorf("error writing header for %s: %v", path, err)

						}

						// Copy the file contents to the tar file
						if _, err := io.Copy(tarWriter, f); err != nil {
							return fmt.Errorf("error copying file contents for %s: %v", path, err)

						}

						fmt.Printf("Added %s to tar file\n", path)
					}
				} else if info.Mode()&os.ModeSymlink != 0 {
					fmt.Printf("Symlink: %s\n", path)
				} else {
					fmt.Printf("Other: %s\n", path)
				}

				return nil
			})

			if err != nil {
				fmt.Printf("Error walking the path %q: %v\n", root, err)
				return
			}

		},
	}

	rootCmd.Flags().StringVar(&startingDate, "starting", "", "Starting date in format yyyy-mm-dd")
	rootCmd.Flags().StringVar(&inPath, "in", "", "Input directory path")
	rootCmd.Flags().StringVar(&outPath, "out", "", "Output file path")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
