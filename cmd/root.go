package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	Version: "v0.0.7",
	Short:   "Process all images and videos inside a directory and move them to a destination",
	Long:    "Process all images and videos inside a directory and move them to a destination",
	Args:    cobra.ExactArgs(2),
	RunE: func(c *cobra.Command, args []string) error {
		sourceArg := args[0]
		targetArg := args[1]
		dryRunFlag := dryRunFlag
		modTimeFlag := modTimeFlag

		// Recursively read source directory
		processErr := filepath.Walk(sourceArg, func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories, process files only
			if fileInfo.IsDir() {
				return nil
			}

			log.Info().Msgf("Processing %s", path)

			// Allow only specified file extensions
			if !util.IsFileExtension(config.FILE_EXTENSIONS_SUPPORTED, path) {
				log.Warn().Msgf("Extension %s not supported", filepath.Ext(path))
				return nil
			}

			// Decode file exif data and parse create date
			fileExif, fileExifErr := exif.Decode(path)
			fileDate, fileDateErr := exif.ParseDate(config.EXIF_FIELDS_DATE_FORMAT, config.EXIF_FIELDS_DATE_CREATED, fileExif)
			if fileExifErr != nil || fileDateErr != nil {
				if modTimeFlag {
					// Use file modtime as fallback
					fileDate = fileInfo.ModTime()
				} else {
					// Move file to unknown
					newPath := filepath.Join(targetArg, "unknown", filepath.Base(path))
					log.Warn().Msg("Failed to parse date (no modtime fallback)")
					log.Info().Msgf("Moving to %s", newPath)

					if dryRunFlag {
						return nil
					}

					return util.MoveFile(path, newPath, config.DEFAULT_DUPLICATE_FILE_STRATEGY)
				}
			}

			// Move file to destination
			yearDir := fmt.Sprintf("%d", fileDate.Year())
			monthDir := fmt.Sprintf("%d-%02d", fileDate.Year(), fileDate.Month())
			fileName := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", fileDate.Year(), fileDate.Month(), fileDate.Day(), fileDate.Hour(), fileDate.Minute(), fileDate.Second(), strings.ToLower(filepath.Ext(path)))
			newPath := filepath.Join(targetArg, yearDir, monthDir, fileName)
			log.Info().Msgf("Moving to %s", newPath)

			if dryRunFlag {
				return nil
			}

			return util.MoveFile(path, newPath, config.DEFAULT_DUPLICATE_FILE_STRATEGY)
		})

		if processErr != nil {
			log.Error().Msgf("%v", processErr)
			log.Error().Msg("View log output above")
			return processErr
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
