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
}

func todoAction(c *cli.Context) error {
	d, err := model.ReadData()
	if err != nil {
		panic(err)
	}

	w := tablewriter.NewWriter(os.Stdout)
	w.SetHeader([]string{"Standard", "Control", "Satisfied?"})

	type row struct {
		standard   string
		controlKey string
		satisfied  string
	}

	satisfied := make(map[string]bool)
	for _, n := range d.Narratives {
		for _, controlKeys := range n.Satisfies {
			for _, key := range controlKeys {
				satisfied[key] = true
			}
		}
	}

	var rows []row
	for _, std := range d.Standards {
		for id, _ := range std.Controls {
			sat := "NO"
			if _, ok := satisfied[id]; ok {
				sat = color.GreenString("YES")
			}

			rows = append(rows, row{
				standard:   std.Name,
				controlKey: id,
				satisfied:  sat,
			})
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].controlKey < rows[j].controlKey
	})

	for _, r := range rows {
		w.Append([]string{r.standard, r.controlKey, r.satisfied})
	}

	w.Render()

	return nil
}
