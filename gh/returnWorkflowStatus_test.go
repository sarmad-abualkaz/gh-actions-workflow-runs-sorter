package gh

import (
    "context"
    "fmt"
    "net/http"
    "reflect"
    "testing"
    "time"

    "github.com/google/go-github/v47/github"
)

func TestReturnWorkflowRunStatustatus(t *testing.T){

    type endpoint struct{
        branch       string
        owner        string
        repo         string
        runId        int
        run          string
    }
    
    type args struct{
        branch       string
        owner        string
        repo         string        
        runId        int
    }

    tests := []struct {
        args           args
        endpoint       endpoint
        name           string
        wantErr        error
        wantStatus     string
        wantUpdateTime *github.Timestamp
    }{
        {
            name: "should succefully return completed run",
            args: args{
                branch:          "ft/test-branch",
                owner:           "testowner",
                repo:            "testrepo",
                runId:           1111111111,
            },
            endpoint: endpoint{
                branch:          "ft/test-branch",
                owner:           "testowner",
                repo:            "testrepo",
                runId:           1111111111,
                run: `{
                        "id": 1111111111,
                        "name": "Test Workflow",
                        "node_id": "fakenode03",
                        "run_number": 3,
                        "event": "push",
                        "status": "completed",
                        "conclusion": "success",
                        "created_at": "2022-12-12T23:34:57Z",
                        "updated_at": "2022-12-12T23:47:06Z"
                    }`,

            },

            wantStatus: "completed",

            wantErr:  nil,

            wantUpdateTime: &github.Timestamp{time.Date(2022, time.December, 12, 23, 47, 06, 0, time.UTC)},
        },
        {
            name: "should succefully return in_progress run",
            args: args{
                branch:          "ft/test-branch",
                owner:           "testowner",
                repo:            "testrepo",
                runId:           1111111111,
            },
            endpoint: endpoint{
                branch:          "ft/test-branch",
                owner:           "testowner",
                repo:            "testrepo",
                runId:           1111111111,
                run: `{
                        "id": 1111111111,
                        "name": "Test Workflow",
                        "node_id": "fakenode03",
                        "run_number": 3,
                        "event": "push",
                        "status": "in_progress",
                        "conclusion": "success",
                        "created_at": "2022-12-12T22:34:57Z",
                        "updated_at": "2022-12-12T22:47:06Z"
                    }`,

            },

            wantStatus: "in_progress",

            wantUpdateTime: &github.Timestamp{time.Date(2022, time.December, 12, 22, 47, 06, 0, time.UTC)},

            wantErr:  nil,
        },

    }

    for _, tt := range tests {

        // var apiurl string
        
        t.Run(tt.name, func(t *testing.T) {

            client, mux, _, teardown := Setup()
            defer teardown()

            ctx := context.Background()
                
            apiurl := fmt.Sprintf("/repos/%s/%s/actions/runs/%d", tt.endpoint.owner, tt.endpoint.repo, tt.endpoint.runId)

            mux.HandleFunc(apiurl, func(w http.ResponseWriter, r *http.Request) {
                TestingMethod(t, r, "GET")
                fmt.Fprint(w, tt.endpoint.run)
            })
            
            gotStatus, gotUpdateTime, gotErr := ReturnWorkflowRunStatus(ctx, client, tt.args.owner, tt.args.repo, tt.args.runId)

            if tt.wantErr == nil {
                
                if gotErr != nil {
                    t.Errorf("ReturnWorkflowRunStatus() returned error: %v expect %v", gotErr, tt.wantErr)
                }

            } else if gotErr.Error() != tt.wantErr.Error() {
                
                t.Errorf("ReturnWorkflowRunStatus() returned error: %v expect %v", gotErr, tt.wantErr)
            }

            if !reflect.DeepEqual(gotStatus, tt.wantStatus){

                t.Errorf("ReturnWorkflowRunStatus() failed - expects %s but received %s", tt.wantStatus, gotStatus)

            }

            if !reflect.DeepEqual(gotUpdateTime, tt.wantUpdateTime){

                t.Errorf("ReturnWorkflowRunStatus() failed - expects %v but received %v", tt.wantUpdateTime, gotUpdateTime)

            }

        })
    }

}
