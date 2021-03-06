## Chapter 6: Monitoring distributed systems

#### Why monitor?
* Analyze long-term trends
* Compare over time or do experiments
* Alerting
* Building dashboards
* Debugging

*Our problem isn’t that we move too slowly, it’s that we build the wrong thing. I wonder how we could get from where we are today to having enough instrumentation to be able to make informed decisions when building new systems.*

#### Setting reasonable expectations
* Monitoring is non-trivial
* 10-12 person SRE team typically has 1-2 people building and maintaining monitoring
* Number has decreased over time due to improvements in tooling/libs/centralized monitoring infra
* General trend towards simpler/faster monitoring systems, with better tools for post hoc analysis
* Avoid “magic” systems (i.e. automatically detect causality)
* Limited success with complex dependency hierarchies (e.g., seldom do “if DB slow, alert for DB, otherwise alert for website”).
    * Used mostly for very stable parts of system
* Rules that generate alerts for humans should be simple to understand and represent a clear failure
* Lots of white-box monitoring
* Some black-box monitoring for critical stuff
* Four golden signals
    * Latency
    * Traffic
    * Errors
    * Saturation

*Interesting examples from Bigtable and Gmail are skipped in notes*

#### Monitoring/alerting philosophy endorsed by Google SRE team
* When creating alerting/monitoring:
    * Rule detect an otherwise undetected condition that is urgent, actionable, and actively user-visible?
    * Am I able to ignore this alert?
    * Are there cases where users not negatively impacted? (i.e. test deployment that could be filtered out)
    * Can I take action on this alert? Urgent or later? Short-term/long-term action?
    * Are other people paged for the same thing?

* When creating page:
    * Every pager, react with sense of urgency. Can only react few times a day.
    * Every page should be actionable
    * Every page response requires intelligence, if robotic response only no need.
    * Pages should be problem that hasn't been seen before.

* The long run
    * There’s often a tension between long-run and short-run availability (hack vs proper fix)
    * Can sometimes fix unreliable systems through heroic effort, but that’s a burnout risk and also a failure risk
    * Taking a controlled hit in short-term reliability is usually the better trade

## Chapter 7: Evolution of automation at Google
* “Automation is a force multiplier, not a panacea”
* Value of automation
    * Consistency
    * Extensibility
    * Centralize mistakes
    * MTTR
    * Faster actions
    * Time savings

*Multiple interesting case studies and explanations skipped in notes.*

## Chapter 8: Release engineering

A specific job function at Google

#### Release engineer role
* Release engineers work with SWEs and SREs to define how software is released
    * Allows dev teams to focus on dev work
* Define best practices
    * Compiler flags, formats for build ID tags, etc.
* Releases automated
* Models vary between teams
    * Could be “push on green” and deploy every build
    * Could be hourly builds and deploys
* Hermetic builds
    * Building same rev number should always give identical results
    * Self-contained -- this includes versioning everything down the compiler used
    * Can cherry-pick fixes against an old rev to fix production software
* Virtually all changes require code review
* Branching
    * All code in main branch
    * Releases are branched off
    * Fixes can go from master to branch
    * Branches never merged back
* Testing
    * CI
    * Release process creates an audit trail that runs tests and shows that tests passed
    
#### Config management
* Although look simple, can cause instability
* Many possible schemes (all involve storing config in source control and having strict config review)
* Use mainline for config -- config maintained at head and applied immediately
    * Originally used for Borg (and pre-Borg systems)
    * Binary releases and config changes decoupled!
* Include config files and binaries in same package
    * Simple
    * Tightly couples binary and config -- ok for projects with few config files or where few configs change
* Package config into “configuration packages”
    * Same hermetic principle as for code
* Release engineering shouldn’t be an afterthought!
* Budget resources at beginning of dev cycle