package github

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/model"
	"golang.org/x/oauth2"
)

const (
	cfgToken    = "token"
	cfgUsername = "username"
	cfgRepo     = "repo"
)

var prompts = map[string]string{
	cfgToken:    "GitHub Token",
	cfgUsername: "GitHub Username",
	cfgRepo:     "GitHub Repository",
}

// Prompts are human-readable configuration element names
func (g *githubPlugin) Prompts() map[string]string {
	return prompts
}

// Register causes the Github plugin to register itself
func Register() {
	model.Register(model.GitHub, &githubPlugin{})
}

type githubPlugin struct {
	token    string
	username string
	reponame string

	clientMu sync.Mutex
	client   *github.Client
}

func (g *githubPlugin) api() *github.Client {
	g.clientMu.Lock()
	defer g.clientMu.Unlock()

	if g.client == nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: g.token},
		)

		// get go-github client
		gh := github.NewClient(oauth2.NewClient(context.Background(), ts))
		g.client = gh
	}
	return g.client
}

func (g *githubPlugin) Get(ID string) (*model.Ticket, error) {
	return nil, nil
}

func (g *githubPlugin) Configured() bool {
	return g.username != "" && g.reponame != "" && g.token != ""
}

func (g *githubPlugin) Links() model.TicketLinks {
	links := model.TicketLinks{}
	links.AuditAll = fmt.Sprintf("https://github.com/%s/%s/issues?q=is%%3Aissue+is%%3Aopen+label%%3Acomply+label%%3Aaudit", g.username, g.reponame)
	links.AuditOpen = fmt.Sprintf("https://github.com/%s/%s/issues?q=is%%3Aissue+is%%3Aopen+label%%3Acomply+label%%3Aaudit", g.username, g.reponame)
	links.ProcedureAll = fmt.Sprintf("https://github.com/%s/%s/issues?q=is%%3Aissue+label%%3Acomply+label%%3Acomply-procedure", g.username, g.reponame)
	links.ProcedureOpen = fmt.Sprintf("https://github.com/%s/%s/issues?q=is%%3Aissue+is%%3Aopen+label%%3Acomply+label%%3Acomply-procedure", g.username, g.reponame)
	return links
}

func (g *githubPlugin) Configure(cfg map[string]interface{}) error {
	var err error

	if g.token, err = getCfg(cfg, cfgToken); err != nil {
		return err
	}
	if g.username, err = getCfg(cfg, cfgUsername); err != nil {
		return err
	}
	if g.reponame, err = getCfg(cfg, cfgRepo); err != nil {
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

func (g *githubPlugin) FindOpen() ([]*model.Ticket, error) {
	issues, _, err := g.api().Issues.ListByRepo(context.Background(), g.username, g.reponame, &github.IssueListByRepoOptions{
		State: "open",
	})

	if err != nil {
		return nil, errors.Wrap(err, "error during FindOpen")
	}

	return toTickets(issues), nil
}

func (g *githubPlugin) FindByTag(name, value string) ([]*model.Ticket, error) {
	panic("not implemented")
}

func (g *githubPlugin) FindByTagName(name string) ([]*model.Ticket, error) {
	issues, _, err := g.api().Issues.ListByRepo(context.Background(), g.username, g.reponame, &github.IssueListByRepoOptions{
		State:  "all",
		Labels: []string{name},
	})

	if err != nil {
		return nil, errors.Wrap(err, "error during FindOpen")
	}

	return toTickets(issues), nil
}

func (g *githubPlugin) LinkFor(t *model.Ticket) string {
	// return fmt.Sprintf("https://github.com/strongdm/comply/issues/%s", t.ID)
	panic("not implemented")
}

func (g *githubPlugin) Create(ticket *model.Ticket, labels []string) error {
	_, _, err := g.api().Issues.Create(context.Background(), g.username, g.reponame, &github.IssueRequest{
		Title:  &ticket.Name,
		Body:   &ticket.Body,
		Labels: &labels,
	})
	return err
}

func toTickets(issues []*github.Issue) []*model.Ticket {
	var tickets []*model.Ticket
	for _, i := range issues {
		tickets = append(tickets, toTicket(i))
	}
	return tickets
}

func toTicket(i *github.Issue) *model.Ticket {
	t := &model.Ticket{Attributes: make(map[string]interface{})}
	t.ID = strconv.Itoa(*i.Number)
	t.Name = ss(i.Title)
	t.Body = ss(i.Body)
	t.CreatedAt = i.CreatedAt
	t.State = toState(ss(i.State))

	for _, l := range i.Labels {
		if l.Name != nil {
			t.SetBool(*l.Name)
		}
	}
	return t
}

func toState(state string) model.TicketState {
	switch state {
	case "closed":
		return model.Closed
	}
	return model.Open
}

func ss(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
