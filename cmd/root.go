package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	Version: "v0.0.10",
	Short:   "process images and videos inside <source> and move them to <target>",
	Long:    "process images and videos inside <source> and move them to <target>",
	Args:    cobra.ExactArgs(2),
	RunE: func(c *cobra.Command, args []string) error {
		defer exif.Instance().Close()

		sourceArg := args[0]
		targetArg := args[1]
		dryRunFlag := dryRunFlag
		modTimeFlag := modTimeFlag

		log.Info().Msg("Reading files...")
		files, err := util.ReadFiles(sourceArg, config.FILE_EXTENSIONS_SUPPORTED)
		if err != nil {
			log.Error().Msgf("Failed to read files: %v", err)
			return err
		}

		log.Info().Msg("Extracting exif...")
		exifs := exif.ExtractData(files...)

		for index, file := range files {
			file := file
			fileExif := exifs[index]

			fileDate, err := exif.ParseDate(fileExif, config.EXIF_FIELDS_DATE_CREATED)
			if fileExif.Err != nil || err != nil {
				if modTimeFlag {
					fileInfo, fileInfoErr := os.Stat(file)
					if fileInfoErr != nil {
						err := moveFileToUnknown(file, targetArg, dryRunFlag)
						if err != nil {
							return err
						}
					} else {
						err := moveFileToTarget(file, fileInfo.ModTime(), targetArg, dryRunFlag)
						if err != nil {
							return err
						}
					}
				} else {
					err := moveFileToUnknown(file, targetArg, dryRunFlag)
					if err != nil {
						return err
					}
				}
			} else {
				err := moveFileToTarget(file, fileDate, targetArg, dryRunFlag)
				if err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&dryRunFlag, "dry", "d", false, "perform a dry run without changes")
	rootCmd.Flags().BoolVarP(&modTimeFlag, "modtime", "m", false, "use file modification time as fallback")
}

func Execute() error {
	return rootCmd.Execute()
}

func moveFileToUnknown(file string, targetArg string, dryRunFlag bool) error {
	newPath := filepath.Join(targetArg, "unknown", filepath.Base(file))

	log.Info().Msgf("Moving %s to %s", file, newPath)
	if !dryRunFlag {
		return util.MoveFile(file, newPath)
	}

	return nil
}

func moveFileToTarget(file string, fileDate time.Time, targetArg string, dryRunFlag bool) error {
	yearDir := fmt.Sprintf("%d", fileDate.Year())
	monthDir := fmt.Sprintf("%d-%02d", fileDate.Year(), fileDate.Month())
	fileName := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", fileDate.Year(), fileDate.Month(), fileDate.Day(), fileDate.Hour(), fileDate.Minute(), fileDate.Second(), strings.ToLower(filepath.Ext(file)))
	newPath := filepath.Join(targetArg, yearDir, monthDir, fileName)

	log.Info().Msgf("Moving %s to %s", file, newPath)
	if !dryRunFlag {
		return util.MoveFile(file, newPath)
	}

	return nil
}
