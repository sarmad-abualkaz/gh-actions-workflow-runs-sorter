package main

import (
    "flag"
    "fmt"
	"time"

    util "gh-actions-workflow-runs-sorter/util"
	gh "gh-actions-workflow-runs-sorter/gh"

    log "github.com/sirupsen/logrus"
)


func main(){

	branch               := flag.String("branch", "main", "git branch name")
	mode                 := flag.String("run-mode", "shouldExecute", "which run mode to run - options available are 'shouldExecute' or 'shouldComplete'")
	owner                := flag.String("owner", "sarmad-abualkaz", "owner of github repo")
	repo                 := flag.String("repo", "test-repo", "github repoistory name")
	runNumber            := flag.Int("run_number", 0, "unique number for each run of a particular workflow in a repository")
	previousRunId        := flag.Int("prev_run_number", 0, "unique number for the previous run of a particular workflow in a repository")
	workflowFile         := flag.String("workflowFile", "cron_and_dispatch.yml", "workflow to link users to")
	workflowRunsToReturn := flag.Int("workflow_run_to_return", 20, "number of workflow runs to return")
	waitBetweenChecks    := flag.Int("wait_between_checks", 10, "how long, in seconds, to wait between checks on previous workflow run")
        waitBeforeComplete   := flag.Float64("wait_before_complete", 60, "how long, in seconds, to wait after a completed previous workflow run")

    flag.Parse()

    // initialize github client

    ctx, client := gh.CreateClient()

    // mode is to check should the workflow execute
    if *mode == "shouldExecute" {

		// get last x number of workflow runs to return (x = workflowRunsToReturn)
        runs, ghErr := gh.ReturnWorkflowRuns(*branch, ctx, client, *owner, *repo, *workflowFile, int(*workflowRunsToReturn))

        if ghErr != nil {
            log.WithFields(log.Fields{
                "repo":         *repo,
                "owner":        *owner,
                "workflowFile": *workflowFile,
                "workflowRunsToReturn": *workflowRunsToReturn,
            }).Error(ghErr.Error())
        }

        // get whether 3 parameters - to be used in the nex mode:
		shouldRunExecute, shouldWaitForPastRun, pastRunIdStr,  ShouldExecuteErr := util.ShouldExecute(runs, *runNumber)

		if ShouldExecuteErr != nil {
            log.WithFields(log.Fields{
                "repo":         *repo,
                "owner":        *owner,
                "workflowFile": *workflowFile,
                "workflowRunsToReturn": *workflowRunsToReturn,
            }).Error(ShouldExecuteErr.Error())
        }

        // export variables retrieved from ShouldExecute() to the environment:
		fmt.Printf("export SHOULD_RUN_EXECUTE=%s\n", shouldRunExecute)
		fmt.Printf("export SHOULD_WAIT_FOR_PAST_RUN=%s\n", shouldWaitForPastRun)
		fmt.Printf("export PAST_RUN_ID=%s\n", pastRunIdStr)

    // check if worklflow should complete:
	} else if *mode == "shouldComplete" {

        // continously loop if status on a previous workflow run is not "completed"
		for {
            
            // retrieve details on last workflow run - status and update_time:
			lastRunStatus, lastRunUpdateTime, ReturnWorkflowRunStatusErr := gh.ReturnWorkflowRunStatus(ctx, client, *owner, *repo, int(*previousRunId))

			if ReturnWorkflowRunStatusErr != nil {

                log.WithFields(log.Fields{
                    "repo":         *repo,
                    "owner":        *owner,
                    "workflowFile": *workflowFile,
                    "workflowRunsToReturn": *workflowRunsToReturn,
                }).Error(ReturnWorkflowRunStatusErr.Error())
            }

            // sleep if status on last workflow run is not "completed"
			if lastRunStatus != "completed" {

                log.WithFields(log.Fields{
                    "repo":             *repo,
                    "owner":            *owner,
                    "previousRunId":    *previousRunId,
                    "currentRunNumber": *runNumber,
                }).Info("must sleep - waiting on previous run to complete ...")
                
                // sleep for provided duration
            	time.Sleep(time.Duration(*waitBetweenChecks)*time.Second)

            // if status is "completed" check if the update_time on last workflow passed provided wait_before_complete
			} else {
                
                // loop until current_time - update_time is less than wait_before_complete
                for {

                    currTime := time.Now()

                    // sleep - if current_time - update_time (on last run) is less than wait_before_complete
                    // use the difference for sleep duration
                    if currTime.Sub(lastRunUpdateTime.Time).Seconds() < *waitBeforeComplete {

                        log.WithFields(log.Fields{
                            "repo":             *repo,
                            "owner":            *owner,
                            "previousRunId":    *previousRunId,
                            "currentRunNumber": *runNumber,
                        }).Info("must sleep - waiting post-completion of previous workflow run ...")

                        log.WithFields(log.Fields{
                            "repo":             *repo,
                            "owner":            *owner,
                            "previousRunId":    *previousRunId,
                            "currentRunNumber": *runNumber,
                        }).Info(fmt.Sprintf("sleeping for %f seconds ...", *waitBeforeComplete - currTime.Sub(lastRunUpdateTime.Time).Seconds()))

                        // sleep for the difference between current_time - update_time (on last workflow)
                        time.Sleep(time.Duration(*waitBeforeComplete - currTime.Sub(lastRunUpdateTime.Time).Seconds())*time.Second)
                    
                    // break if current_time - update_time (on last run) is greater or equal to the wait_before_complete
                    } else {

                        log.WithFields(log.Fields{
                            "repo":             *repo,
                            "owner":            *owner,
                            "currentRunNumber": *runNumber,
                        }).Info("Good to complete this workflow ...")

                        break
                    }

                }

                // break the main loop
                break
			}

		}
    
    // panic if run_mode is neither shouldExecute or shouldComplete
	} else {
		panic(fmt.Sprintf("mode passed is %s - allowed values are shouldExecute or shouldComplete", *mode))
	}

}
