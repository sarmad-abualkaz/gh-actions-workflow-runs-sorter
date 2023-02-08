package util

import (
    "fmt"
    "reflect"
    "testing"
    "time"

    "github.com/google/go-github/v47/github"
)

func TestShouldExecute(t *testing.T){

    tests := []struct {
        name                     string
        runs                     []*github.WorkflowRun
        runNumber                int
        wantShouldExecute        string
        wantShouldWaitForPastRun string
        wantPastRunId            string
        wantError                error
        workflowRunsToReturn     int
    }{
        {
            name: "should not execute",
            runs: []*github.WorkflowRun{
                {ID: github.Int64(3333333333), Name: github.String("Test Workflow"), NodeID: github.String("fakenode03"), RunNumber: github.Int(30), Event: github.String("push") , Status: github.String("completed"), Conclusion: github.String("success"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 23, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 23, 47, 06, 0, time.UTC)}},
                {ID: github.Int64(2222222222), Name: github.String("Test Workflow"), NodeID: github.String("fakenode02"), RunNumber: github.Int(29), Event: github.String("push") , Status: github.String("completed"), Conclusion: github.String("success"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 22, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 22, 47, 06, 0, time.UTC)}},
                {ID: github.Int64(1111111111), Name: github.String("Test Workflow"), NodeID: github.String("fakenode01"), RunNumber: github.Int(28), Event: github.String("push") , Status: github.String("completed"), Conclusion: github.String("success"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 21, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 21, 47, 06, 0, time.UTC)}},
            },
            runNumber: 20,
            wantShouldExecute: "false",
            wantShouldWaitForPastRun: "false",
            wantPastRunId: "0",
            wantError: nil,
        },
        {
            name: "should execute but not wait",
            runs: []*github.WorkflowRun{
                {ID: github.Int64(4444444444), Name: github.String("Test Workflow"), NodeID: github.String("fakenode04"), RunNumber: github.Int(31), Event: github.String("push") , Status: github.String("in_progress"), Conclusion: github.String("tbc"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 24, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 24, 47, 06, 0, time.UTC)}},
                {ID: github.Int64(3333333333), Name: github.String("Test Workflow"), NodeID: github.String("fakenode03"), RunNumber: github.Int(30), Event: github.String("push") , Status: github.String("completed"), Conclusion: github.String("success"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 23, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 23, 47, 06, 0, time.UTC)}},
                {ID: github.Int64(2222222222), Name: github.String("Test Workflow"), NodeID: github.String("fakenode02"), RunNumber: github.Int(29), Event: github.String("push") , Status: github.String("completed"), Conclusion: github.String("success"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 22, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 22, 47, 06, 0, time.UTC)}},
                {ID: github.Int64(1111111111), Name: github.String("Test Workflow"), NodeID: github.String("fakenode01"), RunNumber: github.Int(28), Event: github.String("push") , Status: github.String("completed"), Conclusion: github.String("success"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 21, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 21, 47, 06, 0, time.UTC)}},
            },
            runNumber: 31,
            wantShouldExecute: "true",
            wantShouldWaitForPastRun: "false",
            wantPastRunId: "3333333333",
            wantError: nil,
            workflowRunsToReturn: 20,
        },        
        {
            name: "should execute and wait",
            runs: []*github.WorkflowRun{
                {ID: github.Int64(4444444444), Name: github.String("Test Workflow"), NodeID: github.String("fakenode04"), RunNumber: github.Int(31), Event: github.String("push") , Status: github.String("in_progress"), Conclusion: github.String("tbc"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 24, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 24, 47, 06, 0, time.UTC)}},
                {ID: github.Int64(3333333333), Name: github.String("Test Workflow"), NodeID: github.String("fakenode03"), RunNumber: github.Int(30), Event: github.String("push") , Status: github.String("in_progress"), Conclusion: github.String("tbc"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 23, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 23, 47, 06, 0, time.UTC)}},
                {ID: github.Int64(2222222222), Name: github.String("Test Workflow"), NodeID: github.String("fakenode02"), RunNumber: github.Int(29), Event: github.String("push") , Status: github.String("completed"), Conclusion: github.String("success"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 22, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 22, 47, 06, 0, time.UTC)}},
                {ID: github.Int64(1111111111), Name: github.String("Test Workflow"), NodeID: github.String("fakenode01"), RunNumber: github.Int(28), Event: github.String("push") , Status: github.String("completed"), Conclusion: github.String("success"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 21, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 21, 47, 06, 0, time.UTC)}},
            },
            runNumber: 31,
            wantShouldExecute: "true",
            wantShouldWaitForPastRun: "true",
            wantPastRunId: "3333333333",
            wantError: nil,
            workflowRunsToReturn: 30,       
        },
        {
            name: "should return error - zero runs",
            runs: []*github.WorkflowRun{},
            runNumber: 31,
            wantShouldExecute: "false",
            wantShouldWaitForPastRun: "false",
            wantPastRunId: "0",
            wantError: fmt.Errorf("No previous runs were returned from Github Actions API"),
            workflowRunsToReturn: 20,
        },

    }

    for _, tt := range tests {

        gotShouldExecute, gotShouldWaitForPastRun, gotPastRunId, gotError := ShouldExecute(tt.runs, tt.runNumber, tt.workflowRunsToReturn)

        if tt.wantError == nil {
                
            if gotError != nil {
                t.Errorf("ShouldExecute() returned error: '%v' expect '%v'", gotError, tt.wantError)
            }

        } else if gotError.Error() != tt.wantError.Error() {
            
            t.Errorf("ShouldExecute() returned error: '%v' expect '%v'", gotError, tt.wantError)
        }

        if !reflect.DeepEqual(tt.wantShouldExecute, gotShouldExecute){
            
            t.Errorf("ShouldExecute() failed - shouldExecute? expects '%s' but received '%s'", tt.wantShouldExecute, gotShouldExecute)

        }

        if !reflect.DeepEqual(tt.wantShouldWaitForPastRun, gotShouldWaitForPastRun){
            
            t.Errorf("ShouldExecute() failed - shouldWaitForPastRun? expects '%s' but received '%s'", tt.wantShouldWaitForPastRun, gotShouldWaitForPastRun)

        }

        if !reflect.DeepEqual(tt.wantPastRunId, gotPastRunId){
            
            t.Errorf("ShouldExecute() failed - pastRunId expects '%s' but received '%s'", tt.wantPastRunId, gotPastRunId)

        }
    }

}
