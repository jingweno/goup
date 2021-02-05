package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show installed Go",
		Long:  "Show installed Go versions.",
		RunE:  runShow,
	}
}

func runShow(cmd *cobra.Command, args []string) error {
	vers, err := listGoVers()
	if err != nil {
		return err
	}

	if len(vers) == 0 {
		showGoIfExist()
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Version", "Active"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	for _, ver := range vers {
		if ver.Current {
			table.Append([]string{ver.Ver, "*"})
		} else {
			table.Append([]string{ver.Ver, ""})
		}
	}

	table.Render()

	return nil
}

func showGoIfExist() {
	goBin, err := exec.LookPath("go")
	if err == nil {
		fmt.Printf("No Go is installed by Goup. Using system Go %q.\n", goBin)
	} else {
		fmt.Println("No Go is installed by Goup.")
	}
}

type goVer struct {
	Ver     string
	Current bool
}

func listGoVers() ([]goVer, error) {
	files, err := ioutil.ReadDir(GoupDir())
	if err != nil {
		return nil, err
	}

	current, err := currentGoVersion()
	if err != nil {
		return nil, err
	}

	var vers []goVer
	for _, file := range files {
		if filename := file.Name(); strings.HasPrefix(filename, "go") {
			err := exec.Command(goupVersionBinGoExec(filename), "version").Run()
			if err != nil {
				continue
			}
			vers = append(vers, goVer{
				Ver:     strings.TrimPrefix(filename, "go"),
				Current: current == filename,
			})
		}
	}

	return vers, nil
}

func currentGoVersion() (string, error) {
	current := GoupCurrentDir()
	goroot, err := os.Readlink(current)
	if err != nil {
		return "", err
	}

	return filepath.Base(goroot), nil
}
