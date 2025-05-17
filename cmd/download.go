package cmd

import (
	"fmt"

	"path/filepath"
	"sync"
	"time"
	"vsmod/internal/config"
	"vsmod/internal/files"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var AppFs afero.Fs

var downloadCmd = &cobra.Command{
	Use:     "download",
	Aliases: []string{"dl"},
	Short:   "Download mods defined in a config file",
	Long:    `Download mods defined in a config file. This command will download each mod to the directory set in mods_dir.`,
	Example: `vsmod download --file mods.yaml`,
	PreRun:  toggleDebug,
	Run: func(cmd *cobra.Command, args []string) {
		forceCheck, _ := cmd.Flags().GetBool("force-compatibility-check")
		if err := downloadMods(conf.Mods, conf.GameVersion, forceCheck); err != nil {
			log.Fatalf("error downloading mods: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.PersistentFlags().Bool("force-compatibility-check", false, "Force check for mod compatibility")
	AppFs = afero.NewBasePathFs(afero.NewOsFs(), conf.Dir())
}

func downloadMod(mod config.ConfigFileMod, gameVersion string, forceCheck bool) error {
	modDetails, err := modAPI.GetMod(mod.ID)

	if err != nil {
		return err
	}

	if mod.CompatCheck || forceCheck {
		log.Debugf("Checking compatibility between %s %s and game %s", mod.ID, mod.Version, gameVersion)
		log.Debugf("Force check: %t", forceCheck)

		compatible := false
		for _, release := range modDetails.Releases {
			if release.CompatibleWith(gameVersion) {
				compatible = true
				break
			}
		}
		if !compatible {
			return fmt.Errorf("%s %v not compatible with game version %s", mod.ID, mod.Version, gameVersion)
		}
	}

	release, err := modDetails.Release(mod.Version)
	if err != nil {
		return err
	}

	log.Infof("Downloading mod %s version %s\n", mod.ID, release.Version)
	log.Debugf("Download URL: %s", release.DownloadURL())
	data, err := files.DownloadFile(release.DownloadURL())
	if err != nil {
		return err
	}

	if err := files.WriteFile(AppFs, filepath.Join("Mods", release.FileName), data); err != nil {
		return err
	}

	return nil
}

func downloadMods(mods []config.ConfigFileMod, gameVersion string, forceCheck bool) error {
	start := time.Now()

	var wg sync.WaitGroup
	errCh := make(chan error, len(mods))

	for _, mod := range mods {
		wg.Add(1)
		go func(mod config.ConfigFileMod) {
			defer wg.Done()
			if err := downloadMod(mod, gameVersion, forceCheck); err != nil {
				errCh <- err
			}
		}(mod)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		if len(errs) > 0 {
			for _, err := range errs {
				log.Errorf("%v\n", err)
			}
			return fmt.Errorf("failed to download some mods")
		}
	}

	elapsed := time.Since(start).Seconds()
	log.Infof("Done: %.2fs\n", elapsed)

	return nil
}
