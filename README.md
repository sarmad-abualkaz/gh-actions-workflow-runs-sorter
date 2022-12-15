# gh-actions-workflow-runs-sorter

## What is this?
This is a basic command-line tool aiming to help organize Github Actions' workflow runs in their chronological order.

The tool's main purpose is to help serialize artifact release where Github Actions' runs of the same workflow can run into race-conditions (e.g. two commits pushed in a short timeframe, where the newer commit can finish its workflow run before the previously pushed commit).

Using the tool can help organize releasing artifacts (e.g. Docker images, Helm Charts etc.) in their correct order.

The tool makes use of the `GITHUB_RUN_NUMBER` or `github.run_number` for aligning workflow runs (A unique number for each run of a particular workflow in a repository. This number begins at 1 for the workflow's first run, and increments with each new run. This number does not change if you re-run the workflow run.)

## How does this work?

### Usage:

There are two modes this tool can run with:
1. `shouldExecute` - check if this workflow run should execute (or run) in the first place. If `SHOULD_RUN_EXECUTE` is returned as `true`, the command will also return `SHOULD_WAIT_FOR_PAST_RUN` (either - true/false) and `PAST_RUN_ID` (the workflow run ID with a run_number lower than currently running workflow run).
2. `shouldComplete` - this mode can check if a workflow run with `PAST_RUN_ID` is still running or is `completed`. If the former it will wait based on user-provided wait-time. If the run with `PAST_RUN_ID` is `completed` it will check if the completion time exceeds user-provided pos-completion wait time and complete the running workflow based on pos-completion wait time. If there's a lag required per user-requirement then it will sleep until that time has surpassed post-completion wait.

#### 1. `shouldExecute` Mode

```
gh-actions-workflow-runs-sorter \
  --mode=shouldExecute \
  --run_number=${{ github.run_number }} \
  --branch=<git-branch> --owner=<git-repo-owner> --repo=<git-repo> \
  --workflowFile=<workflow-file-name>
```

#### 2. `shouldComplete` Mode

```
gh-actions-workflow-runs-sorter \
  --mode=shouldComplete \
  --run_number=${{ github.run_number }} \
  --branch=<git-branch> --owner=<git-repo-owner> --repo=<git-repo> \
  --previousRunId=${PAST_RUN_ID} \
  --waitBeforeComplete=<how long to wait after PAST_RUN_ID workflow run completes>
```

#### Flags to note:

| flag | purpose | default |
| --- | --- | --- | 
| `--branch` | which branch to point to for workflow file name | `main` |
|`--mode` | which mode to run this cli with. |`shouldExecute`|
| `--owner` | owner of the git repo where this workflow is running | |
| `--repo` | the git repo where this workflow is running | |
| `--run_number`| the `GITHUB_RUN_NUMBER` or `github.run_number` of currently running workflow run | |
| `--prev_run_number` | used in `shouldComplete` mode the workflow run id `GITHUB_RUN_ID` or `github.run_id` of previous workflow run | | 
| `--workflowFile` | the workflow file name running triggering the workflow | | 
| `--workflow_run_to_return` | how many workflow runs do you want to visit per check | `20` |
| `--wait_between_checks` | used in `shouldComplete` mode when `SHOULD_WAIT_FOR_PAST_RUN` is true - how long to wait before checking the status of workflow run with `previousRunId` again | `10s` |
| `--wait_before_complete` | used in `shouldComplete` mode - how long to wait post-completion of workflow run with `previousRunId` | `300s` |

## Explanation:
Running this cli using the `shouldExecute` mode will return three variables `SHOULD_RUN_EXECUTE`, `SHOULD_WAIT_FOR_PAST_RUN`, and `PAST_RUN_ID`. All three variables are exportable using the cli output - note the command execution below. 

```
$(gh-actions-workflow-runs-sorter \
  --mode=shouldExecute \
  --run_number=${{ github.run_number }} \
  --branch=<git-branch> \
  --owner=<git-repo-owner> \
  --repo=<git-repo> 
  --workflowFile=<workflow-file-name>)

echo ${SHOULD_RUN_EXECUTE}` #should output true or false
echo ${SHOULD_WAIT_FOR_PAST_RUN}` #should output true or false
echo ${PAST_RUN_ID}` #should output an integer
```

### How are variables calculated in `shouldExecute` mode?

### `SHOULD_RUN_EXECUTE`:
The `SHOULD_RUN_EXECUTE` variable is calculated by looking over x number of previous runs from a workflow (x is provided by `--workflow_run_to_return` defaulting to 20). 

If a run with a HIGHER `github.run_number` than what was set in `--run_number` is found to have `completed`, then this variable is set to `false` - since a new commit has already ran and completed - the CURRENT run (with run_number=`--run_number`) has lost its order and should not be executed. 

Setting this to false will gate against manually re-running a previously failed (or succeeded) workflow runs. 

If the last `completed` run is found to have a LOWER `github.run_number` than what was set in `--run_number`, then this variable is set to `true`.


### `SHOULD_WAIT_FOR_PAST_RUN`:
The `SHOULD_WAIT_FOR_PAST_RUN` variable is calculated by looking over x number of previous runs from a workflow (x is provided by `--workflow_run_to_return` defaulting to 20). 

If the last run with a `github.run_number` LOWER than what is set in `--run_number` is found to not be in a `completed` state, then this flag will be set to `true`. Otherwise the assumption is there's no run to wait on.

### `PAST_RUN_ID`:
The `PAST_RUN_ID` variable is calculated by looking over x number of previous runs from a workflow (x is provided by `--workflow_run_to_return` defaulting to 20). 

This will provide the `github.run_id` of the last run found with a `github.run_number` LOWER than what is set in `--run_number`.

### How is wait time calculated in `shouldComplete` mode?
Based on what is provided in `--prev_run_number`, `--waitBetweenChecks` and `--waitBeforeComplete` the following logic will take place:
1. if `prev_run_number` is still not in `completed` state, the tool will wait `--waitBetweenChecks` seconds.
2. retry 1 until `prev_run_number` is in `completed` state.
3. if `prev_run_number` is `completed` check the `LastUpdateTime` on `prev_run_number` workflow run.
4. if `current_time` - (`LastUpdateTime` on `prev_run_number` workflow run) is less than `--waitBeforeComplete` seconds, then sleep for (`--waitBeforeComplete`) - (the diff of current_time - last_update_time on `prev_run_number`).
5. repeat 5 until `current_time` - (`LastUpdateTime` on `prev_run_number` workflow run) is greater than `--waitBeforeComplete` seconds.
6. Exit successfully.


