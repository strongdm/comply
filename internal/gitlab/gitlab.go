package gitlab

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/model"
	"github.com/xanzy/go-gitlab"
)

const (
	cfgDomain = "domain"
	cfgToken  = "token"
	cfgRepo   = "repo"
)

var prompts = map[string]string{
	cfgDomain: "Fully Qualified GitLab Domain",
	cfgToken:  "GitLab Token",
	cfgRepo:   "GitLab Repository",
}

// Prompts are human-readable configuration element names
func (g *gitlabPlugin) Prompts() map[string]string {
	return prompts
}

// Register causes the Github plugin to register itself
func Register() {
	model.Register(model.GitLab, &gitlabPlugin{})
}

type gitlabPlugin struct {
	domain   string
	token    string
	reponame string

	clientMu sync.Mutex
	client   *gitlab.Client
}

func (g *gitlabPlugin) api() *gitlab.Client {
	g.clientMu.Lock()
	defer g.clientMu.Unlock()
	if g.client == nil {
		// get go-gitlab client
		gl := gitlab.NewClient(nil, g.token)
		gl.SetBaseURL(g.domain)
		g.client = gl
	}
	return g.client
}

func (g *gitlabPlugin) Get(ID string) (*model.Ticket, error) {
	return nil, nil
}

func (g *gitlabPlugin) Configured() bool {
	return g.reponame != "" && g.token != ""
}

func (g *gitlabPlugin) Links() model.TicketLinks {
	links := model.TicketLinks{}
	links.AuditAll = fmt.Sprintf("%s/%s/issues?scope=all&utf8=✓&state=all&label_name[]=comply-audit", g.domain, g.reponame)
	links.AuditOpen = fmt.Sprintf("%s/%s/issues?scope=all&utf8=✓&state=opened&label_name[]=comply-audit", g.domain, g.reponame)
	links.ProcedureAll = fmt.Sprintf("%s/%s/issues?scope=all&utf8=✓&state=all&label_name[]=comply-procedure", g.domain, g.reponame)
	links.ProcedureOpen = fmt.Sprintf("%s/%s/issues?scope=all&utf8=✓&state=opened&label_name[]=comply-procedure", g.domain, g.reponame)
	return links
}

func (g *gitlabPlugin) Configure(cfg map[string]interface{}) error {
	var err error

	if g.domain, err = getCfg(cfg, cfgDomain); err != nil {
		return err
	}
	if g.token, err = getCfg(cfg, cfgToken); err != nil {
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

func getProjectIssues(g *gitlabPlugin, options *gitlab.ListProjectIssuesOptions) ([]*gitlab.Issue, error) {
	issues := []*gitlab.Issue{}
	options.Page = 1

	for {
		pageIssues, resp, err := g.api().Issues.ListProjectIssues(g.reponame, options)
		if err != nil {
			return nil, errors.Wrap(err, "error retreiving issues from gitlab")
		}

		issues = append(issues, pageIssues...)

		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		options.Page = resp.NextPage
	}

	return issues, nil
}

func (g *gitlabPlugin) FindOpen() ([]*model.Ticket, error) {
	options := &gitlab.ListProjectIssuesOptions{
		State: gitlab.String("opened"),
	}

	issues, err := getProjectIssues(g, options)
	if err != nil {
		return nil, errors.Wrap(err, "error during FindOpen")
	}

	return toTickets(issues), nil
}

func (g *gitlabPlugin) FindByTag(name, value string) ([]*model.Ticket, error) {
	panic("not implemented")
}

func (g *gitlabPlugin) FindByTagName(name string) ([]*model.Ticket, error) {
	options := &gitlab.ListProjectIssuesOptions{
		State:  gitlab.String("all"),
		Labels: []string{name},
	}

	issues, err := getProjectIssues(g, options)
	if err != nil {
		return nil, errors.Wrap(err, "error during FindByTagName")
	}

	return toTickets(issues), nil
}

func (g *gitlabPlugin) LinkFor(t *model.Ticket) string {
	panic("not implemented")
}

func (g *gitlabPlugin) Create(ticket *model.Ticket, labels []string) error {
	l := gitlab.Labels(labels)
	options := &gitlab.CreateIssueOptions{
		Title:       gitlab.String(ticket.Name),
		Description: gitlab.String(ticket.Body),
		Labels:      &l,
	}
	_, _, err := g.api().Issues.CreateIssue(g.reponame, options)
	return err
}

func toTickets(issues []*gitlab.Issue) []*model.Ticket {
	var tickets []*model.Ticket
	for _, i := range issues {
		tickets = append(tickets, toTicket(i))
	}
	return tickets
}

func toTicket(i *gitlab.Issue) *model.Ticket {
	t := &model.Ticket{Attributes: make(map[string]interface{})}
	t.ID = strconv.Itoa(i.ID)
	t.Name = i.Title
	t.Body = i.Description
	t.CreatedAt = i.CreatedAt
	t.State = toState(i.State)

	for _, l := range i.Labels {
		if l == "audit" {
			t.SetBool("comply-audit")
		}
		if l == "procedure" {
			t.SetBool("comply-procedure")
		}

		// seems redundant, but fixes a bug the other two labels introduce
		// whereby open comply tickets aren't properly accounted for in the UI
		if l == "comply-audit" {
			t.SetBool("comply-audit")
		}
		if l == "comply-procedure" {
			t.SetBool("comply-procedure")
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
