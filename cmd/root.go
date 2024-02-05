package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/patrickap/img-sort/m/v2/internal/config"
	"github.com/patrickap/img-sort/m/v2/internal/exif"
	"github.com/patrickap/img-sort/m/v2/internal/log"
	"github.com/patrickap/img-sort/m/v2/internal/util"
	"github.com/spf13/cobra"
)

var dryRunFlag bool
var modTimeFlag bool

var rootCmd = &cobra.Command{
	Use:     "img-sort <source> <target>",
	Version: "v0.0.8",
	Short:   "Process all images and videos inside a directory and move them to a destination",
	Long:    "Process all images and videos inside a directory and move them to a destination",
	Args:    cobra.ExactArgs(2),
	RunE: func(c *cobra.Command, args []string) error {
		defer exif.Instance().Close()

		sourceArg := args[0]
		targetArg := args[1]
		dryRunFlag := dryRunFlag
		modTimeFlag := modTimeFlag

		log.Info().Msg("Processing files..")
		files := []string{}
		filesErr := filepath.WalkDir(sourceArg, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if entry.IsDir() {
				return nil
			}

			files = append(files, path)
			return nil
		})

		if filesErr != nil {
			log.Error().Msgf("Failed to process files: %v", filesErr)
			return filesErr
		}

		log.Info().Msg("Extracting exif...")
		exifs, exifsErr := exif.Extract(files...)
		if exifsErr != nil {
			log.Error().Msgf("Failed to extract exif: %v", exifsErr)
			return exifsErr
		}

		wg := sync.WaitGroup{}
		filesErrCh := make(chan error, len(files))

		for fileIndex, file := range files {
			file := file
			fileExif := exifs[fileIndex]

			wg.Add(1)
			go func() {
				defer wg.Done()

				if !util.IsFileExtension(config.FILE_EXTENSIONS_SUPPORTED, file) {
					log.Warn().Msgf("Extension %s not supported", filepath.Ext(file))
					return
				}

				fileDate, fileDateErr := exif.ParseDate(config.EXIF_FIELDS_DATE_FORMAT, config.EXIF_FIELDS_DATE_CREATED, fileExif)
				if fileDateErr != nil {
					fileInfo, fileInfoErr := os.Stat(file)

					if modTimeFlag && fileInfoErr == nil {
						fileDate = fileInfo.ModTime()
					} else {
						newPath := filepath.Join(targetArg, "unknown", filepath.Base(file))

						log.Info().Msgf("Moving %s to %s", file, newPath)

						if dryRunFlag {
							return
						}

						moveErr := util.MoveFile(file, newPath, config.DEFAULT_DUPLICATE_FILE_STRATEGY)
						if moveErr != nil {
							filesErrCh <- moveErr
							return
						}
					}
				}

				yearDir := fmt.Sprintf("%d", fileDate.Year())
				monthDir := fmt.Sprintf("%d-%02d", fileDate.Year(), fileDate.Month())
				fileName := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", fileDate.Year(), fileDate.Month(), fileDate.Day(), fileDate.Hour(), fileDate.Minute(), fileDate.Second(), strings.ToLower(filepath.Ext(file)))
				newPath := filepath.Join(targetArg, yearDir, monthDir, fileName)

				log.Info().Msgf("Moving %s to %s", file, newPath)

				if dryRunFlag {
					return
				}

				moveErr := util.MoveFile(file, newPath, config.DEFAULT_DUPLICATE_FILE_STRATEGY)
				if moveErr != nil {
					filesErrCh <- moveErr
					return
				}
			}()
		}

		wg.Wait()
		close(filesErrCh)

		isErr := false
		for fileErr := range filesErrCh {
			if fileErr != nil {
				log.Error().Msgf("Error: %v", fileErr)
				isErr = true
			}
		}

		if isErr {
			log.Info().Msg("Completed with errors")
			return errors.New("completed with errors")
		}

		log.Info().Msg("Completed successfully")
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&dryRunFlag, "dry-run", "D", false, "perform a dry run without modifying data")
	rootCmd.Flags().BoolVarP(&modTimeFlag, "mod-time", "M", false, "use file modification time as fallback")
}

func Execute() error {
	return rootCmd.Execute()
}
