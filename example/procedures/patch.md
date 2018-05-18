id: "patch"
name: "Apply OS patches"
cron: "0 0 0 15 * *"
---

# OS Patch Procedure

Resolve this ticket by executing the following steps:

- [ ] Pull the latest scripts from the Ops repository
- [ ] Execute `ENV=staging patch-all.sh`
- [ ] Inspect output
    - [ ] Errors? Investigate and resolve
- [ ] Execute `ENV=production patch-all.sh`
- [ ] Attach log output to this ticket