# {{.Name}} Compliance Program

This repository consolidates all documents related to the {{.Name}} Compliance Program.

# Structure

Compliance documents are organized as follows:

```
narratives/     Narratives provide an overview of the organization and the compliance environment.
policies/       Policies govern the behavior of employees and contractors.
procedures/     Procedures prescribe specific steps that are taken in response to key events.
standards/      Standards specify the controls satisfied by the compliance program.
templates/      Templates control the output format of the HTML Dashboard and PDF assets.
```

# Building

Assets are built using [`comply`](https://comply.strongdm.com), which can be installed via `brew install comply` (macOS) or `go get github.com/strongdm/comply`

# Publishing

The `output/` directory contains all generated assets. Links in the HTML dashboard a relative, and all dependencies are included via direct CDN references. The entire `output/` directory therefore may be uploaded to an S3 bucket or other static asset host without further modification.

# Dashboard Status

Procedure tracking is updated whenever `comply sync` is invoked. Invoke a sync prior to `comply build` to include the most current ticket status.

# Procedure Scheduler

Any `procedures/` that include a `cron` schedule will automatically created in your configured ticketing system whenever `comply scheduler` is executed. The scheduler will backfill any overdue tickets.

# Deployment Recommendation

Invoke a script similar to the following at least once per day:

```
#!/bin/bash
#
# prerequisites:
#   git access
#   ticketing configuration in comply.yml
#   upload.sh to publish static site
#

# get latest policies and procedures
git pull

# update ticketing status
comply sync

# trigger creation of scheduled tickets
comply scheduler

# build latest
comply build

# publish static site from output/ directory
upload.sh output/
```