package cli

import (
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/strongdm/comply/internal/model"
	"github.com/urfave/cli"
)

var todoCommand = cli.Command{
	Name:   "todo",
	Usage:  "list declared vs satisfied compliance controls",
	Action: todoAction,
	Before: projectMustExist,
}

func todoAction(c *cli.Context) error {
	d, err := model.ReadData()
	if err != nil {
		return err
	}

	w := tablewriter.NewWriter(os.Stdout)
	w.SetHeader([]string{"Standard", "Control", "Satisfied?", "Name"})

	type row struct {
		standard    string
		controlKey  string
		satisfied   string
		controlName string
	}

	satisfied := model.ControlsSatisfied(d)

	var rows []row
	for _, std := range d.Standards {
		for id, c := range std.Controls {
			sat := "NO"
			if _, ok := satisfied[id]; ok {
				sat = color.GreenString("YES")
			}

			rows = append(rows, row{
				standard:    std.Name,
				controlKey:  id,
				satisfied:   sat,
				controlName: c.Name,
			})
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].controlKey < rows[j].controlKey
	})

	w.SetAutoWrapText(false)

	for _, r := range rows {
		w.Append([]string{r.standard, r.controlKey, r.satisfied, r.controlName})
	}

	w.Render()

	return nil
}
