package cmd

import (
	"fmt"

	"path/filepath"
	"sync"
	"time"
	"vsmod/internal/config"
	"vsmod/internal/files"

	"github.com/Masterminds/semver/v3"
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
	PreRun: func(cmd *cobra.Command, args []string) {
		toggleDebug(cmd, args)
		if err := conf.Hooks["download"].Pre_Run.Run(conf); err != nil {
			log.Errorf("error running pre-run hook: %v", err)
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		if err := conf.Hooks["download"].Post_Run.Run(conf); err != nil {
			log.Errorf("error running post-run hook: %v", err)
		}
	},
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

func downloadMod(configMod config.ConfigFileMod, gameVersion semver.Constraints, forceCheck bool) error {
	mod, err := modAPI.GetMod(configMod.ID)

	if err != nil {
		return err
	}

	version, err := mod.Release(configMod.Version)
	if err != nil {
		return fmt.Errorf("mod %s version %s not found: %w", configMod.ID, configMod.Version, err)
	}

	log.Debugf("compatCheck: %v, forceCheck: %v, compatibleWith: %v", configMod.CompatCheck, forceCheck, version.CompatibleWith(gameVersion))
	if (configMod.CompatCheck || forceCheck) && !version.CompatibleWith(gameVersion) {
		log.Debugf("Checking mod %s version %s tags: %v against game version %s", configMod.ID, version.Version, version.Tags, gameVersion)
		log.Debugf("releases: %v", mod.Releases)
		return fmt.Errorf("mod %s version %s is not compatible with game version %s", configMod.ID, version.Version, gameVersion)
	}

	log.Debug("Found compatible releases: ", version)

	log.Infof("Downloading mod %s version %s\n", configMod.ID, version.Version)
	log.Debugf("Download URL: %s", version.DownloadURL())
	data, err := files.DownloadFile(version.DownloadURL())
	if err != nil {
		return err
	}

	if err := files.WriteFile(AppFs, filepath.Join("Mods", version.FileName), data); err != nil {
		return err
	}

	return nil
}

func downloadMods(configFileMods []config.ConfigFileMod, gameVersion semver.Constraints, forceCheck bool) error {
	start := time.Now()

	var wg sync.WaitGroup
	errCh := make(chan error, len(configFileMods))

	for _, configFileMod := range configFileMods {
		wg.Add(1)
		go func(mod config.ConfigFileMod) {
			defer wg.Done()
			if err := downloadMod(mod, gameVersion, forceCheck); err != nil {
				errCh <- err
			}
		}(configFileMod)
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
