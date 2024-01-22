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

var modtimeFlag bool

var rootCmd = &cobra.Command{
	Use:          "img-sort <source> <target>",
	Version:      "v0.0.6",
	Short:        "Process all images and videos inside <source> and move them to <target>",
	Long:         "Process all images and videos inside <source> using exif information and move them to <target>",
	Args:         cobra.ExactArgs(2),
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {
		sourceArg := args[0]
		targetArg := args[1]
		modtimeFlag := modtimeFlag

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
			if !util.IsExtension(path, config.FILE_EXTENSIONS_ALLOWED) {
				log.Warn().Msgf("Extension %s not supported", filepath.Ext(path))
				return nil
			}

			// Decode file exif data and parse create date
			var fileDate time.Time
			var fileError error
			fileExif, fileError := exif.Decode(path)
			fileDate, fileError = exif.ParseDate(fileExif, config.EXIF_FIELDS_DATE_CREATED, config.EXIF_FIELDS_DATE_FORMAT)
			if fileError != nil {
				if !modtimeFlag {
					// Move file to unknown
					newPath := filepath.Join(targetArg, "unknown", filepath.Base(path))
					log.Warn().Msg("Failed to parse date (no modtime fallback)")
					log.Info().Msgf("Moving to %s", newPath)
					return util.Move(path, newPath)
				}

				// Set file modtime as fallback
				fileDate = fileInfo.ModTime()
			}

			// Move file to destination
			yearDir := fmt.Sprintf("%d", fileDate.Year())
			monthDir := fmt.Sprintf("%d-%02d", fileDate.Year(), fileDate.Month())
			fileName := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", fileDate.Year(), fileDate.Month(), fileDate.Day(), fileDate.Hour(), fileDate.Minute(), fileDate.Second(), strings.ToLower(filepath.Ext(path)))
			newPath := filepath.Join(targetArg, yearDir, monthDir, fileName)
			log.Info().Msgf("Moving to %s", newPath)
			return util.Move(path, newPath)
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
	rootCmd.Flags().BoolVarP(&modtimeFlag, "modtime", "M", false, "Use the modification time as fallback when there is no exif information")
}

func Execute() error {
	return rootCmd.Execute()
}
