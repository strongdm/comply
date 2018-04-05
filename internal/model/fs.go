package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/path"
	"gopkg.in/yaml.v2"
)

// ReadData loads all records
func ReadData() *Data {
	d := &Data{}
	d.Tickets = ReadTickets()
	d.Narratives = ReadNarratives()
	d.Policies = ReadPolicies()
	d.Procedures = ReadProcedures()
	d.Standards = ReadStandards()
	return d
}

// ReadTickets returns all known tickets, or an empty list in the event the ticket cache is empty or unavailable
func ReadTickets() []*Ticket {
	rt, err := DB().ReadAll("tickets")
	if err != nil {
		// empty list
		return make([]*Ticket, 0)
	}
	return tickets(rt)
}

func tickets(rawTickets []string) []*Ticket {
	var tickets []*Ticket
	for _, rt := range rawTickets {
		t := &Ticket{}
		err := json.Unmarshal([]byte(rt), t)
		if err != nil {
			panic("Malformed ticket JSON")
		}
		tickets = append(tickets, t)
	}
	return tickets
}

// ReadStandards loads standard definitions from the filesystem
func ReadStandards() []*Standard {
	var standards []*Standard

	for _, f := range path.Standards() {
		s := &Standard{}
		sBytes, err := ioutil.ReadFile(f.FullPath)
		if err != nil {
			panic(err)
		}

		yaml.Unmarshal(sBytes, &s)
		standards = append(standards, s)
	}

	return standards
}

// ReadNarratives loads narrative descriptions from the filesystem
func ReadNarratives() []*Narrative {
	var narratives []*Narrative

	for _, f := range path.Narratives() {
		n := &Narrative{}
		mdmd := loadMDMD(f.FullPath)
		yaml.Unmarshal([]byte(mdmd.yaml), &n)
		n.Body = mdmd.body
		n.FullPath = f.FullPath
		n.ModifiedAt = f.Info.ModTime()
		rewriteFilename := strings.Replace(n.Name, " ", "-", -1)
		n.OutputFilename = fmt.Sprintf("%s-%s.pdf", config.Config().FilePrefix, rewriteFilename)
		narratives = append(narratives, n)
	}

	return narratives
}

// ReadProcedures loads procedure descriptions from the filesystem
func ReadProcedures() []*Procedure {
	var procedures []*Procedure

	for _, f := range path.Procedures() {
		p := &Procedure{}
		mdmd := loadMDMD(f.FullPath)
		yaml.Unmarshal([]byte(mdmd.yaml), &p)
		p.Body = mdmd.body
		p.FullPath = f.FullPath
		p.ModifiedAt = f.Info.ModTime()
		procedures = append(procedures, p)
	}

	return procedures
}

// ReadPolicies loads policy documents from the filesystem
func ReadPolicies() []*Policy {
	var policies []*Policy

	for _, f := range path.Policies() {
		p := &Policy{}
		mdmd := loadMDMD(f.FullPath)
		yaml.Unmarshal([]byte(mdmd.yaml), &p)
		p.Body = mdmd.body
		p.FullPath = f.FullPath
		p.ModifiedAt = f.Info.ModTime()
		p.OutputFilename = fmt.Sprintf("%s-%s.pdf", config.Config().FilePrefix, p.Acronym)

		policies = append(policies, p)
	}

	return policies
}

type metadataMarkdown struct {
	yaml string
	body string
}

func loadMDMD(path string) metadataMarkdown {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	content := string(bytes)
	components := strings.Split(content, "---")
	if len(components) == 1 {
		panic(fmt.Sprintf("Malformed metadata markdown in %s, must be of the form: YAML\\n---\\nmarkdown content", path))
	}
	yaml := components[0]
	body := strings.Join(components[1:], "---")
	return metadataMarkdown{yaml, body}
}
