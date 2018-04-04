// package taskbot

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strconv"
// 	"strings"
// 	"sync"

// 	// "github.com/Sirupsen/logrus"
// 	// "github.com/google/go-github/github"
// 	// "github.com/robfig/cron"
// 	// "github.com/strongdm/web/pkg/errors"
// 	// "github.com/strongdm/web/pkg/tardis"
// )

// var configPath string

// type runConfig struct {
// 	sync.Mutex
// 	LastRepoUpdate       tardis.Time `json:"lastRepoUpdate,omitempty"`
// 	LastIssueCreated     tardis.Time `json:"lastIssueCreated,omitempty"`
// 	LastMilestoneRefresh tardis.Time `json:"lastMilestoneRefresh,omitempty"`
// }

// func (r *runConfig) SetLastRepoUpdate() {
// 	r.Lock()
// 	defer r.Unlock()

// 	r.LastRepoUpdate = tardis.Now()
// }

// func (r *runConfig) SetLastIssueCreated() {
// 	r.Lock()
// 	defer r.Unlock()

// 	r.LastIssueCreated = tardis.Now()
// }

// func (r *runConfig) SetLastMilestoneRefresh(t tardis.Time) {
// 	r.Lock()
// 	defer r.Unlock()

// 	r.LastMilestoneRefresh = t
// }

// func (r *runConfig) Save() error {
// 	f, err := os.Create(filepath.Join(configPath, "config.json"))
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	r.Lock()
// 	defer r.Unlock()

// 	b, err := json.Marshal(r)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = f.Write(b)
// 	return err
// }

// func (d *daemon) Run() error {
// 	runCfg, err := initConfigDir()
// 	if err != nil {
// 		return errors.E(err, "unable to initialize config dir")
// 	}
// 	d.runCfg = runCfg

// 	err = updateTasks.Run()
// 	if err != nil {
// 		return errors.E(err, "unable sync tasks from github on startup")
// 	}
// 	runCfg.SetLastRepoUpdate()
// 	err = runCfg.Save()
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("🤖    updating task definitions")
// 	// scan the repo for changes or new tasks
// 	initialTasks, err := loadTasks(repoPath)
// 	if err != nil {
// 		return errors.E(err, "unable to load tasks on startup")
// 	}

// 	d.tasksMu.Lock()
// 	d.tasks = initialTasks
// 	fmt.Println("🤖    starting scheduler")
// 	d.sched = startCron(d.cl, d.runCfg, d.tasks)
// 	d.tasksMu.Unlock()

// 	// start periodically scanning for task changes
// 	ctx, cancel = context.WithCancel(context.Background())
// 	defer cancel()

// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		t := tardis.Tick(TaskUpdateFrequency)
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case <-t:
// 				cmd := exec.Command("git", "pull")
// 				cmd.Dir = repoPath
// 				err := cmd.Run()
// 				if err != nil {
// 					logrus.WithError(err).Error("Unable to update playbooks from github")
// 					continue
// 				}
// 				runCfg.SetLastRepoUpdate()
// 				err = runCfg.Save()
// 				if err != nil {
// 					logrus.WithError(err).Error("Unable to save last run date")
// 				}

// 				reloaded, err := loadTasks(repoPath)
// 				if err != nil {
// 					logrus.WithError(err).Error("Unable to parse tasks on reload")
// 					continue
// 				}
// 				d.tasksMu.Lock()
// 				d.tasks = reloaded
// 				d.sched.Stop()
// 				d.sched = startCron(d.cl, d.runCfg, d.tasks)
// 				d.tasksMu.Unlock()
// 			}
// 		}
// 	}()

// 	fmt.Println("🤖    startup complete")
// 	wg.Wait()
// 	return nil
// }

// func loadTasks(path string) ([]Task, error) {
// 	tasks := []Task{}
// 	files, err := ioutil.ReadDir(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, f := range files {
// 		taskPath := filepath.Join(path, f.Name())
// 		if f.IsDir() {
// 			t, err := loadTasks(taskPath)
// 			if err != nil {
// 				return tasks, err
// 			}
// 			tasks = append(tasks, t...)
// 			continue
// 		}

// 		if !strings.HasSuffix(f.Name(), TaskFileSuffix) {
// 			continue
// 		}

// 		rdr, err := os.OpenFile(taskPath, os.O_RDONLY, 0)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer rdr.Close()

// 		t, err := parseTask(rdr)
// 		if err != nil {
// 			return tasks, err
// 		}
// 		tasks = append(tasks, t)
// 	}

// 	return tasks, nil
// }

// func startCron(n *notifier, cfg *runConfig, tasks []Task) *cron.Cron {
// 	sched := cron.New()
// 	// daily open task shouter
// 	sched.AddFunc("0 0 15 * * 1-5", func() {
// 		ctx, cancel := context.WithTimeout(context.Background(), 120*tardis.Second)
// 		defer cancel()
// 		issues, _, err := n.Github.Issues.ListByRepo(ctx, "strongdm", "web", &github.IssueListByRepoOptions{Labels: Tags})
// 		if err != nil {
// 			logrus.WithError(err).Error("Error loading issues from Github for the morning shout")
// 		}

// 		if len(issues) == 0 {
// 			return
// 		}

// 		var buf bytes.Buffer
// 		fmt.Fprintln(&buf, "Good morning! I'm Mr. Meeseeks. Here are the open tasks for today.")
// 		for _, i := range issues {
// 			fmt.Fprintf(&buf, "- %s %s\n", *i.Title, *i.HTMLURL)
// 		}

// 		ctx, cancel = context.WithTimeout(context.Background(), 120*tardis.Second)
// 		defer cancel()
// 		_, err = n.Slack.Chat().PostMessage("dev").Text(buf.String()).Do(ctx)
// 		if err != nil {
// 			logrus.WithError(err).Error("Error posting morning shout to slack")
// 		}
// 	})

// 	// milestone shouter
// 	var (
// 		milestoneMu      sync.Mutex
// 		milestoneRunning bool
// 	)
// 	sched.AddFunc("0 * * * * *", func() {
// 		milestoneMu.Lock()
// 		if milestoneRunning {
// 			milestoneMu.Unlock()
// 			return
// 		}
// 		milestoneRunning = true
// 		milestoneMu.Unlock()
// 		defer func() {
// 			milestoneMu.Lock()
// 			milestoneRunning = false
// 			milestoneMu.Unlock()
// 		}()

// 		ctx, cancel := context.WithTimeout(context.Background(), 120*tardis.Second)
// 		defer cancel()

// 		mss, _, err := n.Github.Issues.ListMilestones(ctx, "strongdm", "web", &github.MilestoneListOptions{State: "closed"})
// 		if err != nil {
// 			logrus.WithError(err).Error("Error loading milestones from Github")
// 			return
// 		}

// 		lastClosedAt := cfg.LastMilestoneRefresh

// 		for _, ms := range mss {
// 			ctx, cancel = context.WithTimeout(context.Background(), 120*tardis.Second)
// 			defer cancel()
// 			issues, _, err := n.Github.Issues.ListByRepo(ctx, "strongdm", "web", &github.IssueListByRepoOptions{Milestone: strconv.Itoa(*ms.Number), State: "all"})
// 			if err != nil {
// 				logrus.WithError(err).Error("Error loading issues for Milestone from Github")
// 				return
// 			}

// 			closedIssues := []*github.Issue{}
// 			openIssues := []*github.Issue{}

// 			for _, i := range issues {
// 				if i.ClosedAt != nil {
// 					closedIssues = append(closedIssues, i)
// 					continue
// 				}

// 				openIssues = append(openIssues, i)
// 			}

// 			if !cfg.LastMilestoneRefresh.IsZero() && ms.ClosedAt.After(cfg.LastMilestoneRefresh) {
// 				logrus.Infoln("starting notification of milestone closure", *ms.ID, *ms.Title)
// 				var buf bytes.Buffer
// 				fmt.Fprintf(&buf, "%s milestone was released.\n\n%s\n\n", *ms.Title, *ms.Description)

// 				if len(closedIssues) > 0 {
// 					fmt.Fprintf(&buf, "Closes issues:\n\n")
// 					for _, i := range closedIssues {
// 						fmt.Fprintf(&buf, "- %s %s\n", *i.Title, *i.HTMLURL)
// 					}
// 				}

// 				if len(openIssues) > 0 {
// 					fmt.Fprintf(&buf, "Remaining open issues:\n\n")
// 					for _, i := range openIssues {
// 						fmt.Fprintf(&buf, "- %s %s\n", *i.Title, *i.HTMLURL)
// 					}
// 				}

// 				ctx, cancel = context.WithTimeout(context.Background(), 120*tardis.Second)
// 				defer cancel()
// 				_, err = n.Slack.Chat().PostMessage("releases").Text(buf.String()).Do(ctx)
// 				if err != nil {
// 					logrus.WithError(err).Error("Error posting morning shout to slack")
// 					return
// 				}

// 				ctx, cancel = context.WithTimeout(context.Background(), 120*tardis.Second)
// 				defer cancel()
// 				_, err = n.Slack.Chat().PostMessage("general").Text(buf.String()).Do(ctx)
// 				if err != nil {
// 					logrus.WithError(err).Error("Error posting morning shout to slack")
// 					return
// 				}
// 				logrus.Infoln("completed notification of milestone closure", *ms.ID, *ms.Title)
// 				if ms.ClosedAt.After(lastClosedAt) {
// 					lastClosedAt = *ms.ClosedAt
// 				}
// 			}
// 		}

// 		cfg.SetLastMilestoneRefresh(lastClosedAt)
// 		cfg.Save()
// 	})

// 	catchUp := false
// 	for _, t := range tasks {
// 		s, err := cron.Parse(t.Cron)
// 		if err != nil {
// 			logrus.WithError(err).WithField("task", t).Error("Invalid cron specification")
// 			continue
// 		}

// 		// See if we need to catch up
// 		nextAct := s.Next(cfg.LastIssueCreated)
// 		if !nextAct.IsZero() && nextAct.Before(tardis.Now()) {
// 			err := n.Notify(t)
// 			if err != nil {
// 				logrus.WithError(err).WithField("title", t.Title).Error("Could not catch task")
// 			}
// 			catchUp = true
// 		}

// 		sched.AddFunc(t.Cron, func() {
// 			err := n.Notify(t)
// 			if err != nil {
// 				logrus.WithError(err).WithField("title", t.Title).Error("Could not create new task")
// 			}
// 			cfg.SetLastIssueCreated()
// 			err = cfg.Save()
// 			if err != nil {
// 				logrus.WithError(err).Error("Error saving last run times")
// 			}
// 		})
// 	}

// 	if catchUp {
// 		cfg.SetLastIssueCreated()
// 		if err := cfg.Save(); err != nil {
// 			logrus.WithError(err).Error("Error saving last run times")
// 		}
// 	}

// 	sched.Start()

// 	return sched
// }
