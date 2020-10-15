package commands

import (
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	homedir string
	logger  *logrus.Logger

	ProfileFiles []string
)

func init() {
	logger = logrus.New()

	var err error
	homedir, err = os.UserHomeDir()
	if err != nil {
		logger.Fatal(err)
	}

	ProfileFiles = []string{
		filepath.Join(homedir, ".profile"),
		filepath.Join(homedir, ".zprofile"),
		filepath.Join(homedir, ".bash_profile"),
	}
}

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "goup",
		Short: "The Go installer",
		RunE:  choiceVersion,
	}

	rootCmd.AddCommand(installCmd())
	rootCmd.AddCommand(removeCmd())
	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(showCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(versionCmd())

	return rootCmd
}

func GoupBinDir() string {
	return GoupDir("bin")
}

func GoupCurrentDir() string {
	return GoupDir("current")
}

func GoupEnvFile() string {
	return GoupDir("env")
}

func GoupCurrentBinDir() string {
	return GoupDir("current", "bin")
}

func goupVersionDir(ver string) string {
	return GoupDir(ver)
}

func GoupDir(paths ...string) string {
	elem := []string{homedir, ".go"}
	elem = append(elem, paths...)

	return filepath.Join(elem...)
}

func choiceVersion(cmd *cobra.Command, args []string) error {
	vers, err := listGoVers()
	if err != nil {
		return err
	}

	if len(vers) == 0 {
		showGoIfExist()
		return nil
	}

	var curVer string
	var pos int

	var items = make([]string, 0, len(vers))
	for idx, v := range vers {
		items = append(items, v.Ver)
		if v.Current {
			curVer = v.Ver
			pos = idx
		}
	}

	prompt := promptui.Select{
		Label:     "select a version",
		Items:     items,
		CursorPos: pos,
	}

	_, ver, err := prompt.Run()
	if err != nil {
		return err
	}

	if ver == curVer {
		return nil
	}

	if err := symlink("go" + ver); err != nil {
		return err
	}
	return nil
}
