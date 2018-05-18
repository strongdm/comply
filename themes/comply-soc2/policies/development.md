name: Software Development Lifecycle Policy
acronym: SDLCP
satisfies:
  TSC:
    - CC8.1
majorRevisions:
  - date: Jun 1 2018
    comment: Initial document
---

# Purpose and Scope

a. The purpose of this policy is to define requirements for establishing and maintaining baseline protection standards for company software, network devices, servers, and desktops.

a. This policy applies to all users performing software development, system administration, and management of these activities within the organization. This typically includes employees and contractors, as well as any relevant external parties involved in these activities (hereinafter referred to as “users”). This policy must be made readily available to all users.

a. This policy also applies to enterprise-wide systems and applications developed by the organization or on behalf of the organization for production implementation.

# Background

a. The intent of this policy is to ensure a well-defined, secure and consistent process for managing the entire lifecycle of software and information systems, from initial requirements analysis until system decommission. The policy defines the procedure, roles, and responsibilities, for each stage of the software development lifecycle.

a. Within this policy, the software development lifecycle consists of requirements analysis, architecture and design, development, testing, deployment/implementation, operations/maintenance, and decommission. These processes may be followed in any form; in a waterfall model, it may be appropriate to follow the process linearly, while in an agile development model, the process can be repeated in an iterative fashion.

# References

a. Risk Assessment Policy

# Policy

a. The organization’s Software Development Life Cycle (SDLC) includes the following phases:

    i. Requirements Analysis

    i. Architecture and Design

    i. Testing

    i. Deployment/Implementation

    i. Operations/Maintenance

    i. Decommission

a. During all phases of the SDLC where a system is not in production, the system must not have live data sets that contain information identifying actual people or corporate entities, actual financial data such as account numbers, security codes, routing information, or any other financially identifying data. Information that would be considered sensitive must never be used outside of production environments.

a. The following activities must be completed and/or considered during the requirements analysis phase:

    i. Analyze business requirements.

    i. Perform a risk assessment. More information on risk assessments is discussed in the Risk Assessment Policy (reference (a)).

    i. Discuss aspects of security (e.g., confidentiality, integrity, availability) and how they might apply to this requirement.

    i. Review regulatory requirements and the organization’s policies, standards, procedures and guidelines.

    i. Review future business goals.

    i. Review current business and information technology operations.

    i. Incorporate program management items, including:

        1. Analysis of current system users/customers.

        1. Understand customer-partner interface requirements (e.g., business-level, network).

        1. Discuss project timeframe.

    i. Develop and prioritize security solution requirements.

    i. Assess cost and budget constraints for security solutions, including development and operations.

    i. Approve security requirements and budget.

    i. Make “buy vs. build” decisions for security services based on the information above.

a. The following must be completed/considered during the architecture and design phase:

    i. Educate development teams on how to create a secure system.

    i. Develop and/or refine infrastructure security architecture.

    i. List technical and non-technical security controls.

    i. Perform architecture walkthrough.

    i. Create a system-level security design.

    i. Create high-level non-technical and integrated technical security designs.

    i. Perform a cost/benefit analysis for design components.

    i. Document the detailed technical security design.

    i. Perform a design review, which must include, at a minimum, technical reviews of application and infrastructure, as well as a review of high-level processes.

    i. Describe detailed security processes and procedures, including: segregation of duties and segregation of development, testing and production environments. 

    i. Design initial end-user training and awareness programs.

    i. Design a general security test plan.

    i. Update the organization’s policies, standards, and procedures, if appropriate.

    i. Assess and document how to mitigate residual application and infrastructure vulnerabilities.

    i. Design and establish separate development and test environments.

a. The following must be completed and/or considered during the development phase:

    i. Set up a secure development environment (e.g., servers, storage).

    i. Train infrastructure teams on installation and configuration of applicable software, if required.

    i. Develop code for application-level security components.

    i. Install, configure and integrate the test infrastructure.

    i. Set up security-related vulnerability tracking processes.

    i. Develop a detailed security test plan for current and future versions (i.e., regression testing).

    i. Conduct unit testing and integration testing.

a. The following must be completed and/or considered during the testing phase:

    i. Perform a code and configuration review through both static and dynamic analysis of code to identify vulnerabilities. 

    i. Test configuration procedures.

    i. Perform system tests.

    i. Conduct performance and load tests with security controls enabled.

    i. Perform usability testing of application security controls.


    i. Conduct independent vulnerability assessments of the system, including the infrastructure and application.

a. The following must be completed and/or considered during the deployment phase:

    i. Conduct pilot deployment of the infrastructure, application and other relevant components.

    i. Conduct transition between pilot and full-scale deployment.

    i. Perform integrity checking on system files to ensure authenticity.

    i. Deploy training and awareness programs to train administrative personnel and users in the system’s security functions.

    i. Require participation of at least two developers in order to conduct full-scale deployment to the production environment.

a. The following must be completed and/or considered during the operations/maintenance phase:

    i. Several security tasks and activities must be routinely performed to operate and administer the system, including but not limited to:

        1. Administering users and access.

        1. Tuning performance.

        1. Performing backups according to requirements defined in the System Availability Policy 

        1. Performing system maintenance (i.e., testing and applying security updates and patches).

        1. Conducting training and awareness.

        1. Conducting periodic system vulnerability assessments.

        1. Conducting annual risk assessments.

    i. Operational systems must:

        1. Be reviewed to ensure that the security controls, both automated and manual, are functioning correctly and effectively.

        1. Have logs that are periodically reviewed to evaluate the security of the system and validate audit controls.

        1. Implement ongoing monitoring of systems and users to ensure detection of security violations and unauthorized changes.

        1. Validate the effectiveness of the implemented security controls through security training as required by the Procedure For Executing Incident Response.

        1. Have a software application and/or hardware patching process that is performed regularly in order to eliminate software bug and security problems being introduced into the organization’s technology environment. Patches and updates must be applied within ninety (90) days of release to provide for adequate testing and propagation of software updates. Emergency, critical, break-fix, and zero-day vulnerability patch releases must be applied as quickly as possible.

a. The following must be completed and/or considered during the decommission phase:

    i. Conduct unit testing and integration testing on the system after component removal.

    i. Conduct operational transition for component removal/replacement.

    i. Determine data retention requirements for application software and systems data.

    i. Document the detailed technical security design.

    i. Update the organization’s policies, standards and procedures, if appropriate.

    i. Assess and document how to mitigate residual application and infrastructure vulnerabilities.

