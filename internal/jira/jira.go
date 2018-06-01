package jira

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/model"

	jira "github.com/andygrunwald/go-jira"
)

const (
	cfgUsername = "username"
	cfgPassword = "password"
	cfgURL      = "url"
	cfgProject  = "project"
)

var prompts = map[string]string{
	cfgUsername: "Jira Username",
	cfgPassword: "Jira Password",
	cfgURL:      "Jira URL",
	cfgProject:  "Jira Project Code",
}

// Prompts are human-readable configuration element names
func (j *jiraPlugin) Prompts() map[string]string {
	return prompts
}

// Register causes the Github plugin to register itself
func Register() {
	model.Register(model.Jira, &jiraPlugin{})
}

type jiraPlugin struct {
	username string
	password string
	url      string
	project  string

	clientMu sync.Mutex
	client   *jira.Client
}

func (j *jiraPlugin) api() *jira.Client {
	j.clientMu.Lock()
	defer j.clientMu.Unlock()

	if j.client == nil {
		tp := jira.BasicAuthTransport{
			Username: j.username,
			Password: j.password,
		}

		client, _ := jira.NewClient(tp.Client(), j.url)
		j.client = client
	}
	return j.client
}

func (j *jiraPlugin) Get(ID string) (*model.Ticket, error) {
	return nil, nil
}

func (j *jiraPlugin) Configured() bool {
	return j.username != "" && j.password != "" && j.url != "" && j.project != ""
}

func (j *jiraPlugin) Links() model.TicketLinks {
	links := model.TicketLinks{}
	// http://localhost:8080/issues/?jql=labels+%3D+comply

	links.AuditAll = fmt.Sprintf("%s/issues?q=is%3Aissue+is%3Aopen+label%3Acomply+label%3Aaudit", j.url)
	links.AuditOpen = fmt.Sprintf("%s/issues?q=is%3Aissue+is%3Aopen+label%3Acomply+label%3Aaudit", j.url)
	links.ProcedureAll = fmt.Sprintf("%s/issues?q=is%3Aissue+label%3Acomply+label%3Aprocedure", j.url)
	links.ProcedureOpen = fmt.Sprintf("%s/issues?q=is%3Aissue+is%3Aopen+label%3Acomply+label%3Aprocedure", j.url)
	return links
}

func (j *jiraPlugin) Configure(cfg map[string]interface{}) error {
	var err error

	if j.username, err = getCfg(cfg, cfgUsername); err != nil {
		return err
	}
	if j.password, err = getCfg(cfg, cfgPassword); err != nil {
		return err
	}
	if j.url, err = getCfg(cfg, cfgURL); err != nil {
		return err
	}
	if j.project, err = getCfg(cfg, cfgProject); err != nil {
		return err
	}

	return nil
}

func getCfg(cfg map[string]interface{}, k string) (string, error) {
	v, ok := cfg[k]
	if !ok {
		return "", errors.New("Missing key: " + k)
	}

	vS, ok := v.(string)
	if !ok {
		return "", errors.New("Malformatted key: " + k)
	}
	return vS, nil
}

func (j *jiraPlugin) FindOpen() ([]*model.Ticket, error) {
	return []*model.Ticket{}, nil
	// issues, _, err := j.api().Issues.ListByRepo(context.Background(), j.username, j.reponame, &github.IssueListByRepoOptions{
	// 	State: "open",
	// })

	// if err != nil {
	// 	return nil, errors.Wrap(err, "error during FindOpen")
	// }

	// return toTickets(issues), nil
}

func (j *jiraPlugin) FindByTag(name, value string) ([]*model.Ticket, error) {
	panic("not implemented")
}

func (j *jiraPlugin) FindByTagName(name string) ([]*model.Ticket, error) {
	issues, _, err := j.api().Issue.Search("labels=comply", &jira.SearchOptions{MaxResults: 1000})
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch Jira issues")
	}
	return toTickets(issues), nil
}

func (j *jiraPlugin) LinkFor(t *model.Ticket) string {
	// return fmt.Sprintf("https://github.com/strongdm/comply/issues/%s", t.ID)
	panic("not implemented")
}

func (j *jiraPlugin) Create(ticket *model.Ticket, labels []string) error {
	i := jira.Issue{
		Fields: &jira.IssueFields{
			Type: jira.IssueType{
				Name: "Task",
			},
			Project: jira.Project{
				Key: j.project,
			},
			Summary:     ticket.Name,
			Description: ticket.Body,
			Labels:      labels,
		},
	}

	_, _, err := j.api().Issue.Create(&i)
	if err != nil {
		return errors.Wrap(err, "unable to create ticket")
	}
	return nil
}

func toTickets(issues []jira.Issue) []*model.Ticket {
	var tickets []*model.Ticket
	for _, i := range issues {
		tickets = append(tickets, toTicket(&i))
	}
	return tickets
}

func toTicket(i *jira.Issue) *model.Ticket {
	t := &model.Ticket{Attributes: make(map[string]interface{})}
	t.ID = i.ID
	t.Name = i.Fields.Description
	t.Body = i.Fields.Summary
	createdAt := time.Time(i.Fields.Created)
	t.CreatedAt = &createdAt
	t.State = toState(i.Fields.Status)

	for _, l := range i.Fields.Labels {
		t.SetBool(l)
	}
	return t
}

func toState(status *jira.Status) model.TicketState {
	switch status.Name {
	case "Closed":
		return model.Closed
	}
	return model.Open
}
