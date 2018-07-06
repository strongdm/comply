package gitlab

import (
	// "context"
	"fmt"
	"strconv"
	"sync"

	"github.com/xanzy/go-gitlab"
	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/model"
	// "golang.org/x/oauth2"
)

const (
	cfgDomain   = "domain"
	cfgToken    = "token"
	cfgRepo     = "repo"
)

var prompts = map[string]string{
	cfgDomain:   "Fully Qualified GitLab Domain",
	cfgToken:    "GitLab Token",
	cfgRepo:     "GitLab Repository",
}

// Prompts are human-readable configuration element names
func (g *gitlabPlugin) Prompts() map[string]string {
	return prompts
}

// Register causes the Github plugin to register itself
func Register() {
	fmt.Printf("register %s \n", "gitlab")
	// fmt.Printf("%s -> %s\n", rel, filepath.Join("output", p.OutputFilename))
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
	fmt.Printf("api called %s \n", "gitlab")
	g.clientMu.Lock()
	defer g.clientMu.Unlock()
	fmt.Printf("g.client: %s \n", g)
	if g.client == nil {
		// ts := oauth2.StaticTokenSource(
		// 	&oauth2.Token{AccessToken: g.token},
		// )

		// get go-gitlab client
		gl := gitlab.NewClient(nil, g.token)
		gl.SetBaseURL(g.domain)
		// fmt.Printf("%s+v\n", gl)
		g.client = gl
	}
	return g.client
}

func (g *gitlabPlugin) Get(ID string) (*model.Ticket, error) {
	fmt.Printf("get ID %s \n", ID)
	return nil, nil
}

func (g *gitlabPlugin) Configured() bool {
	return g.reponame != "" && g.token != ""
}

func (g *gitlabPlugin) Links() model.TicketLinks {
	links := model.TicketLinks{}
	links.AuditAll = fmt.Sprintf("%s/%s/issues?scope=all&utf8=✓&state=all&label_name[]=audit", g.domain, g.reponame)
	links.AuditOpen = fmt.Sprintf("%s/%s/issues?scope=all&utf8=✓&state=opened&label_name[]=audit", g.domain, g.reponame)
	links.ProcedureAll = fmt.Sprintf("%s/%s/issues?scope=all&utf8=✓&state=all&label_name[]=procedure", g.domain, g.reponame)
	links.ProcedureOpen = fmt.Sprintf("%s/%s/issues?scope=all&utf8=✓&state=opened&label_name[]=procedure", g.domain, g.reponame)
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
	fmt.Printf("get cfg %s \n", "gitlab")
	fmt.Printf("checking config %s \n", cfg[k])
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

func (g *gitlabPlugin) FindOpen() ([]*model.Ticket, error) {
	fmt.Printf("find open %s \n", "gitlab")
	options := &gitlab.ListProjectIssuesOptions{
		State: gitlab.String("opened"),
	}

	issues, _, err := g.api().Issues.ListProjectIssues(g.reponame, options)

	if err != nil {
		return nil, errors.Wrap(err, "error during FindOpen")
	}

	return toTickets(issues), nil
}

func (g *gitlabPlugin) FindByTag(name, value string) ([]*model.Ticket, error) {
	panic("not implemented")
}

func (g *gitlabPlugin) FindByTagName(name string) ([]*model.Ticket, error) {
	fmt.Printf("> find by tag name %s \n", "gitlab")
	options := &gitlab.ListProjectIssuesOptions{
		State:  gitlab.String("all"),
		Labels: []string{name},
	}

	issues, _, err := g.api().Issues.ListProjectIssues(g.reponame, options)

	if err != nil {
		return nil, errors.Wrap(err, "error during FindOpen")
	}

	return toTickets(issues), nil
}

func (g *gitlabPlugin) LinkFor(t *model.Ticket) string {
	// return fmt.Sprintf("https://github.com/strongdm/comply/issues/%s", t.ID)
	panic("not implemented")
}

func (g *gitlabPlugin) Create(ticket *model.Ticket, labels []string) error {
	options :=  &gitlab.CreateIssueOptions{
		Title:         gitlab.String(ticket.Name),
		Description:   gitlab.String(ticket.Body),
		Labels:        labels,
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
			t.SetBool("audit")
		}
		if l == "procedure" {
			t.SetBool("procedure")
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
