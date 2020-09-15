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
	Usage:  "list declared vs satisfied compliance criteria",
	Action: todoAction,
	Before: projectMustExist,
}

func todoAction(c *cli.Context) error {
	d, err := model.ReadData()
	if err != nil {
		return err
	}

	w := tablewriter.NewWriter(os.Stdout)
	w.SetHeader([]string{"Framework", "Criterion", "Satisfied?", "Name"})

	type row struct {
		framework    string
		criterionKey  string
		satisfied   string
		criterionName string
	}

	satisfied := model.CriteriaSatisfied(d)

	var rows []row
	for _, std := range d.Frameworks {
		for id, c := range std.Criteria{
			sat := "NO"
			if _, ok := satisfied[id]; ok {
				sat = color.GreenString("YES")
			}

			rows = append(rows, row{
				framework:    std.Name,
				criterionKey:  id,
				satisfied:   sat,
				criterionName: c.Name,
			})
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].criterionKey < rows[j].criterionKey
	})

	w.SetAutoWrapText(false)

	for _, r := range rows {
		w.Append([]string{r.framework, r.criterionKey, r.satisfied, r.criterionName})
	}

	w.Render()

	return nil
}
