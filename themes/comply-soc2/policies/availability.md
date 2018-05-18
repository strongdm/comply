name: Availability Policy
acronym: AP
satisfies:
  TSC:
    - A1.1
    - CC9.1
majorRevisions:
  - date: Jun 1 2018
    comment: Initial document
---

# Purpose and Scope

a. The purpose of this policy is to define requirements for proper controls to protect the availability of the organization’s information systems.

a. This policy applies to all users of information systems within the organization. This typically includes employees and contractors, as well as any external parties that come into contact with systems and information controlled by the organization (hereinafter referred to as “users”). This policy must be made readily available to all users.

# Background

a. The intent of this policy is to minimize the amount of unexpected or unplanned downtime (also known as outages) of information systems under the organization’s control. This policy prescribes specific measures for the organization that will increase system redundancy, introduce failover mechanisms, and implement monitoring such that outages are prevented as much as possible. Where they cannot be prevented, outages will be quickly detected and remediated.

a. Within this policy, an availability is defined as a characteristic of information or information systems in which such information or systems can be accessed by authorized entities whenever needed.

# References

a. Risk Assessment Policy

# Policy

a. Information systems must be consistently available to conduct and support business operations.

a. Information systems must have a defined availability classification, with appropriate controls enabled and incorporated into development and production processes based on this classification.

a. System and network failures must be reported promptly to the organization’s lead for Information Technology (IT) or designated IT operations manager.

a. Users must be notified of scheduled outages (e.g., system maintenance) that require periods of downtime. This notification must specify the date and time of the system maintenance, expected duration, and anticipated system or service resumption time.

a. Prior to production use, each new or significantly modified application must have a completed risk assessment that includes availability risks. Risk assessments must be completed in accordance with the Risk Assessment Policy (reference (a)).

a. Capacity management and load balancing techniques must be used, as deemed necessary, to help minimize the risk and impact of system failures.

a. Information systems  must have an appropriate data backup plan that ensures:

    i. All sensitive data can be restored within a reasonable time period.

    i. Full backups of critical resources are performed on at least a weekly basis.

    i. Incremental backups for critical resources are performed on at least a daily basis.

    i. Backups and associated media are maintained for a minimum of thirty (30) days and retained for at least one (1) year, or in accordance with legal and regulatory requirements.

    i. Backups are stored off-site with multiple points of redundancy and protected using encryption and key management.

    i. Tests of backup data must be conducted once per quarter. Tests of configurations must be conducted twice per year.

a. Information systems  must have an appropriate redundancy and failover plan that meets the following criteria:

    i. Network infrastructure that supports critical resources must have system-level redundancy (including but not limited to a secondary power supply, backup disk-array, and secondary computing system). Critical core components (including but not limited to routers, switches, and other devices linked to Service Level Agreements (SLAs)) must have an actively maintained spare. SLAs must require parts replacement within twenty-four (24) hours.

    i. Servers that support critical resources must have redundant power supplies and network interface cards. All servers must have an actively maintained spare. SLAs must require parts replacement within twenty-four (24) hours.

    i. Servers classified as high availability must use disk mirroring.

a. Information systems must have an appropriate business continuity plan that meets the following criteria:

    i. Recovery time and data loss limits are defined in Table 3. 

    i. Recovery time requirements and data loss limits must be adhered to with specific documentation in the plan.

    i. Company and/or external critical resources, personnel, and necessary corrective actions must be specifically identified.

    i. Specific responsibilities and tasks for responding to emergencies and resuming business operations must be included in the plan.

    i. All applicable legal and regulatory requirements must be satisfied.

+-------------------+------------------+---------------+-------------------+------------------+
|**Availability**   | **Availability** | **Scheduled** | **Recovery Time** | **Data Loss or** |
|**Classification** | **Requirements** | **Outage**    | **Requirements**  | **Impact Loss**  |
+===================+==================+===============+===================+==================+
| High              | High to          | 30 minutes    | 1 hour            | Minimal          |
|                   | Continuous       |               |                   |                  |
+-------------------+------------------+---------------+-------------------+------------------+
|                   |                  |               |                   |                  |
+-------------------+------------------+---------------+-------------------+------------------+
| Medium            | Standard         | 2 hours       | 4 hours           | Some data loss   |
|                   | Availability     |               |                   | is tolerated if  |
|                   |                  |               |                   | it results in    |
|                   |                  |               |                   | quicker          |
|                   |                  |               |                   | restoration      |
+-------------------+------------------+---------------+-------------------+------------------+
|                   |                  |               |                   |                  |
+-------------------+------------------+---------------+-------------------+------------------+
| Low               | Limited          | 4 hours       | Next              | Some data loss   |
|                   | Availability     |               | business day      | is tolerated if  |
|                   |                  |               |                   | it results in    |
|                   |                  |               |                   | quicker          |
|                   |                  |               |                   | restoration      |
+-------------------+------------------+---------------+-------------------+------------------+

Table 3: Recovery Time and Data Loss Limits 
