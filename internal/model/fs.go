package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/path"
	"gopkg.in/yaml.v2"
)

// ReadData loads all records from both the filesystem and ticket cache.
func ReadData() (*Data, error) {
	tickets, err := ReadTickets()
	if err != nil {
		return nil, err
	}
	narratives, err := ReadNarratives()
	if err != nil {
		return nil, err
	}
	policies, err := ReadPolicies()
	if err != nil {
		return nil, err
	}
	procedures, err := ReadProcedures()
	if err != nil {
		return nil, err
	}
	standards, err := ReadStandards()
	if err != nil {
		return nil, err
	}

	return &Data{
		Tickets:    tickets,
		Narratives: narratives,
		Policies:   policies,
		Procedures: procedures,
		Standards:  standards,
	}, nil
}

// ReadTickets returns all known tickets, or an empty list in the event the ticket cache is empty or unavailable.
func ReadTickets() ([]*Ticket, error) {
	rt, err := DB().ReadAll("tickets")
	if err != nil {
		// empty list
		return []*Ticket{}, nil
	}
	return tickets(rt)
}

func tickets(rawTickets []string) ([]*Ticket, error) {
	var tickets []*Ticket
	for _, rt := range rawTickets {
		t := &Ticket{}
		err := json.Unmarshal([]byte(rt), t)
		if err != nil {
			return nil, errors.Wrap(err, "malformed ticket JSON")
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

// ReadStandards loads standard definitions from the filesystem.
func ReadStandards() ([]*Standard, error) {
	var standards []*Standard

	files, err := path.Standards()
	if err != nil {
		return nil, errors.Wrap(err, "unable to enumerate paths")
	}

	for _, f := range files {
		s := &Standard{}
		sBytes, err := ioutil.ReadFile(f.FullPath)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read "+f.FullPath)
		}

		yaml.Unmarshal(sBytes, &s)
		standards = append(standards, s)
	}

	return standards, nil
}

// ReadNarratives loads narrative descriptions from the filesystem.
func ReadNarratives() ([]*Document, error) {
	var narratives []*Document

	files, err := path.Narratives()
	if err != nil {
		return nil, errors.Wrap(err, "unable to enumerate paths")
	}

	for _, f := range files {
		n := &Document{}
		mdmd := loadMDMD(f.FullPath)
		err = yaml.Unmarshal([]byte(mdmd.yaml), &n)
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse "+f.FullPath)
		}
		n.Body = mdmd.body
		n.FullPath = f.FullPath
		n.ModifiedAt = f.Info.ModTime()
		n.OutputFilename = fmt.Sprintf("%s-%s.pdf", config.Config().FilePrefix, n.Acronym)
		narratives = append(narratives, n)
	}

	return narratives, nil
}

// ReadProcedures loads procedure descriptions from the filesystem.
func ReadProcedures() ([]*Procedure, error) {
	var procedures []*Procedure
	files, err := path.Procedures()
	if err != nil {
		return nil, errors.Wrap(err, "unable to enumerate paths")
	}

	for _, f := range files {
		p := &Procedure{}
		mdmd := loadMDMD(f.FullPath)
		err = yaml.Unmarshal([]byte(mdmd.yaml), &p)
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse "+f.FullPath)
		}
		p.Body = mdmd.body
		p.FullPath = f.FullPath
		p.ModifiedAt = f.Info.ModTime()
		procedures = append(procedures, p)
	}

	return procedures, nil
}

// ReadPolicies loads policy documents from the filesystem.
func ReadPolicies() ([]*Document, error) {
	var policies []*Document

	files, err := path.Policies()
	if err != nil {
		return nil, errors.Wrap(err, "unable to enumerate paths")
	}

	for _, f := range files {
		p := &Document{}
		mdmd := loadMDMD(f.FullPath)
		err = yaml.Unmarshal([]byte(mdmd.yaml), &p)
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse "+f.FullPath)
		}
		p.Body = mdmd.body
		p.FullPath = f.FullPath
		p.ModifiedAt = f.Info.ModTime()
		p.OutputFilename = fmt.Sprintf("%s-%s.pdf", config.Config().FilePrefix, p.Acronym)
		policies = append(policies, p)
	}

	return policies, nil
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
