name: Security Architecture Narrative
acronym: SEN
satisfies:
  TSC:
    - CC6.6
    - CC6.7
    - CC7.1
    - CC7.2
majorRevisions:
  - date: Jun 1 2018
    comment: Initial document
---

# Security Architecture Narrative

Here we narrate why our org satisfies the control keys listed in the YML block

# {{.Name}} Product Architecture

Describe product architecture here, emphasizing security implications

# {{.Name}} Infrastructure

## Product Infrastructure

Describe product infrastructure, emphasizing security measures

### Authorized Personnel

- **AWS root account** access is granted only to the CTO and CEO
- **AWS IAM** access is granted to to a limited group of **Operators**
- **{{.Name}} SSH** access is granted to a limited group of **Operators**
- **{{.Name}} DB** access is granted to a limited group of **Data Operators**

## IT Infrastructure

{{.Name}} uses the following cloud services for its internal infrastructure:

- List cloud services

Access to these cloud services is limited according to the role of the {{.Name}} employee and is reviewed quarterly as well as via regular onboarding/offboarding tasks for new and departing employees.

# {{.Name}} Workstations

{{.Name}} workstations are hardened against logical and physical attack by the following measures:

- operating system must be within one generation of current
- full-disk encryption
- onboard antivirus/antimalware software
- OS and AV automatically updated

Workstation compliance with these measures is evaluated on a quarterly basis.

## Remote Access

Many {{.Name}} employees work remotely on a regular basis and connect to production and internal IT systems via the same methods as those employees connecting from the {{.Name}} physical office, i.e., direct encrypted access to cloud services. It is the employee's responsibility to ensure that only authorized personnel use {{.Name}} resources and access {{.Name}} systems.

# Access Review

Access to {{.Name}} infrastructure, both internal and product, is reviewed quarterly and inactive users are removed. Any anomalies are reported to the security team for further investigation. When employees start or depart, an onboarding/offboarding procedure is followed to provision or deprovision appropriate account access.

# Penetration Testing

{{.Name}} commissions an external penetration test on an annual basis. All findings are immediately reviewed and addressed to the satisfaction of the CTO/CEO.

# {{.Name}} Physical Security

{{.Name}} has one physical location, in San Francisco, CA. Key issuance is tracked by the Office Physical Security Policy Ledger. Office keys are additionally held by the lessor, property management, and custodial staff. These keys are not tracked by the Office Physical Security Policy Ledger. {{.Name}} managers regularly review physical access privileges.

{{.Name}} infrastructure is located within AWS. {{.Name}} does not have physical access to AWS infrastructure.

# Risk Assessment

{{.Name}} updates its Cyber Risk Assessment on an annual basis in order to keep pace with the evolving threat landscape. The following is an inventory of adversarial and non-adversarial threats assessed to be of importance to {{.Name}}.

## Adversarial Threats

The following represents the inventory of adversarial threats:

|Threat|Source|Vector|Target|Likelihood|Severity|
|----------------------------+--------------+------------+-----------------+----------+------|
| | | | | | |

## Non-Adversarial Threats

The following represents the inventory of non-adversarial threats:

|Threat|Vector|Target|Likelihood|Severity|
|----------------------------+--------------+-------------+----------+------|
| | | | | |

# References

## Narratives

Products and Services Narrative
System Architecture Narrative

## Policies

Encryption Policy
Log Management Policy
Office Security Policy
Remote Access Policy
Security Incident Response Policy
Workstation Policy

## Procedures

Apply OS Patches
Review & Clear Low-Priority Alerts
Review Access
Review Devices & Workstations
