# gh-actions-workflow-runs-sorter

## What is this?
This is a basic command-line tool aiming to help organize Github Actions' workflow runs in their coronological order.

The tool's main purpose is to help serialize arifact release where Github Actions' runs of the same workflow can run into race-conditions (e.g. two commits pushed in a short timeframe, where the newer commit can finish it's workflow run before the previously pushed commit).

Using the tool can help organize releasing artifacts (e.g. Docker images, Helm Charts etc.) in their correct order.

The tool makes use of the `GITHUB_RUN_NUMBER` or `github.run_number` for aligning workflow runs (A unique number for each run of a particular workflow in a repository. This number begins at 1 for the workflow's first run, and increments with each new run. This number does not change if you re-run the workflow run.)

## How does this work

There are two modes this tool can run with:
1. shouldExecute - check if this workflow run should execute (or run) in the first place. If `SHOULD_RUN_EXECUTE` is returned as `true`, the command will also return `SHOULD_WAIT_FOR_PAST_RUN` (either - true/false) and `PAST_RUN_ID` (the workflow run ID with a run_number lower than currently running workflow run).
2. shouldComplete - this mode can check if a workflow run with `PAST_RUN_ID` is still running or is `completed`. If the former it will wait based on user-provided wait-time. If the run with `PAST_RUN_ID` is `completed` it will check if the completion time exceeds user-provided pos-completion wait time and complete the running workflow if it exceeds. If there's a lag required per user-requirement then it will sleep until that time has surparsed post-completion wait.

### 1. `shouldExecute` Mode
#### Usage

