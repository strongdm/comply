package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/aktau/github-release/github"
)

func infocmd(opt Options) error {
	user := nvls(opt.Info.User, EnvUser)
	repo := nvls(opt.Info.Repo, EnvRepo)
	token := nvls(opt.Info.Token, EnvToken)
	tag := opt.Info.Tag

	if user == "" || repo == "" {
		return fmt.Errorf("user and repo need to be passed as arguments")
	}

	// Find regular git tags.
	foundTags, err := Tags(user, repo, token)
	if err != nil {
		return fmt.Errorf("could not fetch tags, %v", err)
	}
	if len(foundTags) == 0 {
		return fmt.Errorf("no tags available for %v/%v", user, repo)
	}

	tags := foundTags[:0]
	for _, t := range foundTags {
		// If the user only requested one tag, filter out the rest.
		if tag == "" || t.Name == tag {
			tags = append(tags, t)
		}
	}

	renderer := renderInfoText

	if opt.Info.JSON {
		renderer = renderInfoJSON
	}

	// List releases + assets.
	var releases []Release
	if tag == "" {
		// Get all releases.
		vprintf("%v/%v: getting information for all releases\n", user, repo)
		releases, err = Releases(user, repo, token)
		if err != nil {
			return err
		}
	} else {
		// Get only one release.
		vprintf("%v/%v/%v: getting information for the release\n", user, repo, tag)
		release, err := ReleaseOfTag(user, repo, tag, token)
		if err != nil {
			return err
		}
		releases = []Release{*release}
	}

	return renderer(tags, releases)
}

func renderInfoText(tags []Tag, releases []Release) error {
	fmt.Println("tags:")
	for _, tag := range tags {
		fmt.Println("-", &tag)
	}

	fmt.Println("releases:")
	for _, release := range releases {
		fmt.Println("-", &release)
	}

	return nil
}

func renderInfoJSON(tags []Tag, releases []Release) error {
	out := struct {
		Tags     []Tag
		Releases []Release
	}{
		Tags:     tags,
		Releases: releases,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	return enc.Encode(&out)
}

func uploadcmd(opt Options) error {
	user := nvls(opt.Upload.User, EnvUser)
	repo := nvls(opt.Upload.Repo, EnvRepo)
	token := nvls(opt.Upload.Token, EnvToken)
	tag := opt.Upload.Tag
	name := opt.Upload.Name
	label := opt.Upload.Label
	file := opt.Upload.File

	vprintln("uploading...")

	if file == nil {
		return fmt.Errorf("provided file was not valid")
	}
	defer file.Close()

	if err := ValidateCredentials(user, repo, token, tag); err != nil {
		return err
	}

	// Find the release corresponding to the entered tag, if any.
	rel, err := ReleaseOfTag(user, repo, tag, token)
	if err != nil {
		return err
	}

	// If asked to replace, first delete the existing asset, if any.
	if assetID := findAssetID(rel.Assets, name); opt.Upload.Replace && assetID != -1 {
		URL := nvls(EnvApiEndpoint, github.DefaultBaseURL) +
			fmt.Sprintf(ASSET_DOWNLOAD_URI, user, repo, assetID)
		resp, err := github.DoAuthRequest("DELETE", URL, "application/json", token, nil, nil)
		if err != nil || resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("could not replace asset %s (ID: %d), deletion failed (error: %v, status: %s)",
				name, assetID, err, resp.Status)
		}
	}

	v := url.Values{}
	v.Set("name", name)
	if label != "" {
		v.Set("label", label)
	}

	url := rel.CleanUploadUrl() + "?" + v.Encode()

	resp, err := github.DoAuthRequest("POST", url, "application/octet-stream",
		token, nil, file)
	if err != nil {
		return fmt.Errorf("can't create upload request to %v, %v", url, err)
	}
	defer resp.Body.Close()

	vprintln("RESPONSE:", resp)
	if resp.StatusCode != http.StatusCreated {
		if msg, err := ToMessage(resp.Body); err == nil {
			return fmt.Errorf("could not upload, status code (%v), %v",
				resp.Status, msg)
		}
		return fmt.Errorf("could not upload, status code (%v)", resp.Status)
	}

	if VERBOSITY != 0 {
		vprintf("BODY: ")
		if _, err := io.Copy(os.Stderr, resp.Body); err != nil {
			return fmt.Errorf("while reading response, %v", err)
		}
	}

	return nil
}

func downloadcmd(opt Options) error {
	user := nvls(opt.Download.User, EnvUser)
	repo := nvls(opt.Download.Repo, EnvRepo)
	token := nvls(opt.Download.Token, EnvToken)
	tag := opt.Download.Tag
	name := opt.Download.Name
	latest := opt.Download.Latest

	vprintln("downloading...")

	if err := ValidateTarget(user, repo, tag, latest); err != nil {
		return err
	}

	// Find the release corresponding to the entered tag, if any.
	var rel *Release
	var err error
	if latest {
		rel, err = LatestRelease(user, repo, token)
	} else {
		rel, err = ReleaseOfTag(user, repo, tag, token)
	}
	if err != nil {
		return err
	}

	assetID := findAssetID(rel.Assets, name)
	if assetID == -1 {
		return fmt.Errorf("coud not find asset named %s", name)
	}

	var resp *http.Response
	if token == "" {
		// Use the regular github.com site it we don't have a token.
		resp, err = http.Get("https://github.com" + fmt.Sprintf("/%s/%s/releases/download/%s/%s", user, repo, tag, name))
	} else {
		url := nvls(EnvApiEndpoint, github.DefaultBaseURL) + fmt.Sprintf(ASSET_DOWNLOAD_URI, user, repo, assetID)
		resp, err = github.DoAuthRequest("GET", url, "", token, map[string]string{
			"Accept": "application/octet-stream",
		}, nil)
	}
	if err != nil {
		return fmt.Errorf("could not fetch releases, %v", err)
	}
	defer resp.Body.Close()

	vprintln("GET", resp.Request.URL, "->", resp)

	contentLength, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github did not respond with 200 OK but with %v", resp.Status)
	}

	out := os.Stdout // Pipe the asset to stdout by default.
	if isCharDevice(out) {
		// If stdout is a char device, assume it's a TTY (terminal). In this
		// case, don't pipe th easset to stdout, but create it as a file in
		// the current working folder.
		if out, err = os.Create(name); err != nil {
			return fmt.Errorf("could not create file %s", name)
		}
		defer out.Close()
	}

	return mustCopyN(out, resp.Body, contentLength)
}

// mustCopyN attempts to copy exactly N bytes, if this fails, an error is
// returned.
func mustCopyN(w io.Writer, r io.Reader, n int64) error {
	an, err := io.Copy(w, r)
	if an != n {
		return fmt.Errorf("data did not match content length %d != %d", an, n)
	}
	return err
}

func ValidateTarget(user, repo, tag string, latest bool) error {
	if user == "" {
		return fmt.Errorf("empty user")
	}
	if repo == "" {
		return fmt.Errorf("empty repo")
	}
	if tag == "" && !latest {
		return fmt.Errorf("empty tag")
	}
	return nil
}

func ValidateCredentials(user, repo, token, tag string) error {
	if err := ValidateTarget(user, repo, tag, false); err != nil {
		return err
	}
	if token == "" {
		return fmt.Errorf("empty token")
	}
	return nil
}

func releasecmd(opt Options) error {
	cmdopt := opt.Release
	user := nvls(cmdopt.User, EnvUser)
	repo := nvls(cmdopt.Repo, EnvRepo)
	token := nvls(cmdopt.Token, EnvToken)
	tag := cmdopt.Tag
	name := nvls(cmdopt.Name, tag)
	desc := nvls(cmdopt.Desc, tag)
	target := nvls(cmdopt.Target)
	draft := cmdopt.Draft
	prerelease := cmdopt.Prerelease

	vprintln("releasing...")

	if err := ValidateCredentials(user, repo, token, tag); err != nil {
		return err
	}

	// Check if we need to read the description from stdin.
	if desc == "-" {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("could not read description from stdin: %v", err)
		}
		desc = string(b)
	}

	params := ReleaseCreate{
		TagName:         tag,
		TargetCommitish: target,
		Name:            name,
		Body:            desc,
		Draft:           draft,
		Prerelease:      prerelease,
	}

	/* encode params as json */
	payload, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("can't encode release creation params, %v", err)
	}
	reader := bytes.NewReader(payload)

	URL := nvls(EnvApiEndpoint, github.DefaultBaseURL) + fmt.Sprintf("/repos/%s/%s/releases", user, repo)
	resp, err := github.DoAuthRequest("POST", URL, "application/json", token, nil, reader)
	if err != nil {
		return fmt.Errorf("while submitting %v, %v", string(payload), err)
	}
	defer resp.Body.Close()

	vprintln("RESPONSE:", resp)
	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == 422 {
			return fmt.Errorf("github returned %v (this is probably because the release already exists)",
				resp.Status)
		}
		return fmt.Errorf("github returned %v", resp.Status)
	}

	if VERBOSITY != 0 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error while reading response, %v", err)
		}
		vprintln("BODY:", string(body))
	}

	return nil
}

func editcmd(opt Options) error {
	cmdopt := opt.Edit
	user := nvls(cmdopt.User, EnvUser)
	repo := nvls(cmdopt.Repo, EnvRepo)
	token := nvls(cmdopt.Token, EnvToken)
	tag := cmdopt.Tag
	name := nvls(cmdopt.Name, tag)
	desc := nvls(cmdopt.Desc, tag)
	draft := cmdopt.Draft
	prerelease := cmdopt.Prerelease

	vprintln("editing...")

	if err := ValidateCredentials(user, repo, token, tag); err != nil {
		return err
	}

	id, err := IdOfTag(user, repo, tag, token)
	if err != nil {
		return err
	}

	vprintf("release %v has id %v\n", tag, id)

	// Check if we need to read the description from stdin.
	if desc == "-" {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("could not read description from stdin: %v", err)
		}
		desc = string(b)
	}

	/* the release create struct works for editing releases as well */
	params := ReleaseCreate{
		TagName:    tag,
		Name:       name,
		Body:       desc,
		Draft:      draft,
		Prerelease: prerelease,
	}

	/* encode the parameters as JSON, as required by the github API */
	payload, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("can't encode release creation params, %v", err)
	}

	URL := nvls(EnvApiEndpoint, github.DefaultBaseURL) + fmt.Sprintf("/repos/%s/%s/releases/%d", user, repo, id)
	resp, err := github.DoAuthRequest("PATCH", URL, "application/json", token, nil, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("while submitting %v, %v", string(payload), err)
	}
	defer resp.Body.Close()

	vprintln("RESPONSE:", resp)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 422 {
			return fmt.Errorf("github returned %v (this is probably because the release already exists)",
				resp.Status)
		}
		return fmt.Errorf("github returned unexpected status code %v", resp.Status)
	}

	if VERBOSITY != 0 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error while reading response, %v", err)
		}
		vprintln("BODY:", string(body))
	}

	return nil
}

func deletecmd(opt Options) error {
	user, repo, token, tag := nvls(opt.Delete.User, EnvUser),
		nvls(opt.Delete.Repo, EnvRepo),
		nvls(opt.Delete.Token, EnvToken),
		opt.Delete.Tag
	vprintln("deleting...")

	id, err := IdOfTag(user, repo, tag, token)
	if err != nil {
		return err
	}

	vprintf("release %v has id %v\n", tag, id)

	baseURL := nvls(EnvApiEndpoint, github.DefaultBaseURL)
	resp, err := github.DoAuthRequest("DELETE", baseURL+fmt.Sprintf("/repos/%s/%s/releases/%d",
		user, repo, id), "application/json", token, nil, nil)
	if err != nil {
		return fmt.Errorf("release deletion failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete the release corresponding to tag %s on repo %s/%s",
			tag, user, repo)
	}

	return nil
}
