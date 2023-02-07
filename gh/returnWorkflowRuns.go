package gh

import (
    "context"
    "fmt"

    "github.com/google/go-github/v47/github"

    log "github.com/sirupsen/logrus"
)

func ReturnWorkflowRuns(branchName string, ctx context.Context, client *github.Client, owner string, repo string, workflowFile string, workflowRunsToReturn int) ([]*github.WorkflowRun, error) {

    log.WithFields(log.Fields{
        "repo":         repo,
        "owner":        owner,
        "workflowFile": workflowFile,
        "workflowRunsToReturn": workflowRunsToReturn,
    }).Info("Calling for last few runs from workflow...")

    opts := &github.ListWorkflowRunsOptions{
        Branch: branchName,
        ListOptions: github.ListOptions{
            Page: 1,
            PerPage: workflowRunsToReturn,
        },
    }

    runs, res, err := client.Actions.ListWorkflowRunsByFileName(ctx, owner, repo, workflowFile, opts)

    if res.StatusCode == 404 {

        log.WithFields(log.Fields{
            "Response Status":      res.StatusCode,
            "repo":                 repo,
            "owner":                owner,
            "workflowFile":         workflowFile,
            "workflowRunsToReturn": workflowRunsToReturn,
        }).Warn("Workflow not found ...")

        return nil, fmt.Errorf("Workflow not found")

    }

    if res.StatusCode == 410 {

        log.WithFields(log.Fields{
            "Response Status":      res.StatusCode,
            "repo":                 repo,
            "owner":                owner,
            "workflowFile":         workflowFile,
            "workflowRunsToReturn": workflowRunsToReturn,
        }).Warn("received 410 code: API Method Gone...")

        return nil, fmt.Errorf("API Method Gone")
    }


    if err != nil {

        return nil, err
    }

    log.WithFields(log.Fields{
        "repo":         repo,
        "owner":        owner,
        "workflowFile": workflowFile,
        "workflowRunsToReturn": workflowRunsToReturn,
    }).Info("Runs were returned ...")

    return runs.WorkflowRuns, nil

}
