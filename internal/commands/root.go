package commands

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	homedir string
	logger  *logrus.Logger

	ProfileFiles []string

	rootCmdVerboseFlag bool
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
		Use:               "goup",
		Short:             "The Go installer",
		PersistentPreRunE: preRunRoot,
		RunE:              runChooseVersion,
	}

	rootCmd.PersistentFlags().BoolVarP(&rootCmdVerboseFlag, "verbose", "v", false, "Verbose")

	rootCmd.AddCommand(installCmd())
	rootCmd.AddCommand(removeCmd())
	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(showCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(upgradeCmd())

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

func goupVersionBinGoExec(ver string) string {
	return GoupDir(ver, "bin", "go"+exe())
}

func GoupDir(paths ...string) string {
	elem := []string{homedir, ".go"}
	elem = append(elem, paths...)

	return filepath.Join(elem...)
}

func preRunRoot(cmd *cobra.Command, args []string) error {
	if rootCmdVerboseFlag {
		logger.SetLevel(logrus.DebugLevel)
	}

	return nil
}

func runChooseVersion(cmd *cobra.Command, args []string) error {
	vers, err := listGoVers()
	if err != nil {
		return err
	}

	if len(vers) == 0 {
		showGoIfExist()
		return nil
	}

	var (
		pos int
	)

	var items = make([]string, 0, len(vers))

	for idx, v := range vers {
		items = append(items, v.Ver)
		if v.Current {
			pos = idx
		}
	}

	prompt := promptui.Select{
		Label:     "Select a version",
		Items:     items,
		CursorPos: pos,
	}

	_, ver, err := prompt.Run()
	if err != nil {
		return err
	}

	if err := symlink("go" + ver); err != nil {
		return err
	}

	logger.Printf("Default Go is set to '%s'", ver)

	return nil
}

func exe() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}
