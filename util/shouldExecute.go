package util

import (
    "fmt"
    "strconv"

    "github.com/google/go-github/v47/github"
    log "github.com/sirupsen/logrus"
)

func ShouldExecute(runs []*github.WorkflowRun, runNumber int)(string, string, string, error){

    /*

    sudo-code logic:

    > should it pubish?
	
    - if no (latest completed/successful run has a higher run_number)

    - if yes
        > should it wait - prev run_number is in progress?
    */

    // currentRunNumber, err := strconv.Atoi(runNumber)

    // if err != nil {
	
    // 	return "", "", "", err

    // }

    // set the execution run requirement to false by default
    shouldRunExecute := "false"
    shouldWaitForPastRun := "false"
    pastRunId := int64(0)
	
    // loop through each run from this workflow:
    for _, run := range runs {

        // latest completed/successful run has a higher run_number:
        if (*run.RunNumber > runNumber) && (*run.Status == "completed") {

            log.WithFields(log.Fields{
                "runNumber": runNumber,
            }).Warn(fmt.Sprintf("There's no need to re-run this workflow run; latest 'future' workflow run has completed with id %d\n", *run.RunNumber))

            // do not update the shouldRunExecute, shouldWaitForPastRun or pastRunId
            // break loop - the rest of the logic is not required
            break

        // found the first previous run with a complete status:
        } else if (*run.RunNumber < runNumber) && (*run.Status == "completed") {

            // do not update shouldWaitForPastRun - fine as false
            // update shouldRunExecute and pastRunId
            shouldRunExecute = "true"
            pastRunId = *run.ID
            break

        } else if (*run.RunNumber < runNumber) && (*run.Status != "completed") {

            // update shouldRunExecute, shouldWaitForPastRun and pastRunId
            shouldRunExecute = "true"
            shouldWaitForPastRun = "true"
            pastRunId = *run.ID
            break
        }
    }

    // convert to string from pastRunId
    pastRunIdStr := strconv.Itoa(int(pastRunId))

    // -- Post-Logic check --
    // Check for the size of the list of runs - expected to be minimum of 20
    // Q: when can it not be 20?
    // A: if this cli is used on a new repository that had fewer workflow run executions than 20
    //    or
    //    if gh-action api has failed to return any runs for the workflow

    // log warning if less than 20 runs were returned:
    switch runs_length := len(runs); {

    // return error if less than 1 (or zero) runs were retunred:
    case runs_length < 1:
        return shouldRunExecute, shouldWaitForPastRun, pastRunIdStr, fmt.Errorf("No previous runs were returned from Github Actions API")

    case runs_length < 20:
        log.WithFields(log.Fields{
            "runNumber":               runNumber,
            "number of previous runs": len(runs),
        }).Warn(fmt.Sprintf("Number of workflow runs recieved from API is less than 20."))
    }

    log.WithFields(log.Fields{
        "runNumber": runNumber,
    }).Info(fmt.Sprintf("updating data with SHOULD_RUN_EXECUTE = %s; SHOULD_WAIT_FOR_PAST_RUN = %s; PAST_RUN_ID = %s\n", shouldRunExecute, shouldWaitForPastRun, pastRunIdStr))

    return shouldRunExecute, shouldWaitForPastRun, pastRunIdStr, nil
}

/*
sudo-code logic:

> should it pubish?
	
    - if no (latest completed/successful run has a higher run_number)

    - if yes
        > should it wait - prev run_number is in progress?

            - if yes
                > At the end shoud it complete?

                    - if yes (prev run_number is completed - success/failed):
                        > sleep x seconds for post success-complete of prev
                        or
                        > don't sleep if prev run_number failed

                    - if no
                        > sleep x seconds - waiting for previous run to complete
*/
