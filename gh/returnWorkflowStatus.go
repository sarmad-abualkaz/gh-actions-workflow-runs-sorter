package gh

import (
    "context"
    "fmt"
    "time"

    "github.com/google/go-github/v47/github"

    log "github.com/sirupsen/logrus"
)

func ReturnWorkflowRunStatus(ctx context.Context, client *github.Client, owner string, repo string, workflowRunId int) (string, *github.Timestamp, error) {

    log.WithFields(log.Fields{
        "repo":         repo,
        "owner":        owner,
        "workflowRunId": workflowRunId,
    }).Info("Calling for a previous workflow RunId...")

    run, res, err := client.Actions.GetWorkflowRunByID(ctx, owner, repo, int64(workflowRunId))

    if res.StatusCode == 404 {

        log.WithFields(log.Fields{
            "repo":         repo,
            "owner":        owner,
            "workflowRunId": workflowRunId,
        }).Warn("Workflow not found ...")

        return "", &github.Timestamp{time.Time{}}, fmt.Errorf("Workflow run not found")

    }

    if res.StatusCode == 410 {

        log.WithFields(log.Fields{
            "repo":         repo,
            "owner":        owner,
            "workflowRunId": workflowRunId,
        }).Warn("received 410 code: API Method Gone...")

        return "", &github.Timestamp{time.Time{}}, fmt.Errorf("API Method Gone")
    }

    if res.StatusCode != 200 {

        log.WithFields(log.Fields{
            "Response Status":      res.StatusCode,
            "repo":         repo,
            "owner":        owner,
            "workflowRunId": workflowRunId,
        }).Warn("Request did not succeed: Response status received was not 200 ...")

        return "", &github.Timestamp{time.Time{}}, fmt.Errorf("Response status received was not 200")
    }

    if err != nil {

        return "", &github.Timestamp{time.Time{}}, err

    }

    log.WithFields(log.Fields{
        "repo":         repo,
        "owner":        owner,
        "workflowRunId": workflowRunId,
    }).Info("Workflow run was returned ...")

    return *run.Status, run.UpdatedAt, nil

}
