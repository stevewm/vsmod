package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all mods in the specified config file",
	Long: `List all mods in the specified config file. 

	This command will show each mod's ID along with its current version and the latest available version.`,
	Example: `vsmod list --file mods.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := listMods(); err != nil {
			log.Fatalf("Error listing mods %v", err)

		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listMods() error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Current Version", "Latest Version"})

	for _, mod := range conf.Mods {
		modDetails, err := modAPI.GetMod(mod.ID)
		if err != nil {
			return err
		}

		latestVersion, err := modDetails.LatestRelease()

		if err != nil {
			t.AppendRow(table.Row{mod.ID, mod.Version, "Error fetching latest version"})
			continue
		}
		t.AppendRow(table.Row{mod.ID, mod.Version, latestVersion.Version})
	}
	t.Render()
	return nil
}
