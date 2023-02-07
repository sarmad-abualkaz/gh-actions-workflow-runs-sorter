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

func TestReturnWorkflowRuns(t *testing.T){

    type endpoint struct{
        branch       string
        owner        string
        repo         string
        workflowFile string
        runs         string
    }
    
    type args struct{
        branch       string
        httpstatus   int
        owner        string
        repo         string        
        workflowFile string
    }

    tests := []struct {
        name       string
        args       args
        endpoint   endpoint
        wantRuns   []*github.WorkflowRun
        wantErr  error
    }{
        {
            name: "should succefully return 3 runs",
            args: args{
                branch:          "ft/test-branch",
                httpstatus:      200,
                owner:           "testowner",
                repo:            "testrepo",
                workflowFile:    "testfile.yaml",
            },
            endpoint: endpoint{
                branch:          "ft/test-branch",
                owner:           "testowner",
                repo:            "testrepo",
                workflowFile:    "testfile.yaml",
                runs: `{"total_count":3,"workflow_runs":[
                    {
                        "id": 3333333333,
                        "name": "Test Workflow",
                        "node_id": "fakenode03",
                        "run_number": 3,
                        "event": "push",
                        "status": "completed",
                        "conclusion": "success",
                        "created_at": "2022-12-12T23:34:57Z",
                        "updated_at": "2022-12-12T23:47:06Z"
                    },
                    {
                        "id": 2222222222,
                        "name": "Test Workflow",
                        "node_id": "fakenode02",
                        "run_number": 2,
                        "event": "push",
                        "status": "completed",
                        "conclusion": "success",
                        "created_at": "2022-12-12T22:34:57Z",
                        "updated_at": "2022-12-12T22:47:06Z"
                    },
                    {
                        "id": 1111111111,
                        "name": "Test Workflow",
                        "node_id": "fakenode01",
                        "run_number": 1,
                        "event": "push",
                        "status": "completed",
                        "conclusion": "success",
                        "created_at": "2022-12-12T21:34:57Z",
                        "updated_at": "2022-12-12T21:47:06Z"
					}
				]}`,

            },
            wantRuns: []*github.WorkflowRun{
                {ID: github.Int64(3333333333), Name: github.String("Test Workflow"), NodeID: github.String("fakenode03"), RunNumber: github.Int(3), Event: github.String("push") , Status: github.String("completed"), Conclusion: github.String("success"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 23, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 23, 47, 06, 0, time.UTC)}},
                {ID: github.Int64(2222222222), Name: github.String("Test Workflow"), NodeID: github.String("fakenode02"), RunNumber: github.Int(2), Event: github.String("push") , Status: github.String("completed"), Conclusion: github.String("success"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 22, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 22, 47, 06, 0, time.UTC)}},
                {ID: github.Int64(1111111111), Name: github.String("Test Workflow"), NodeID: github.String("fakenode01"), RunNumber: github.Int(1), Event: github.String("push") , Status: github.String("completed"), Conclusion: github.String("success"), CreatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 21, 34, 57, 0, time.UTC)}, UpdatedAt: &github.Timestamp{time.Date(2022, time.December, 12, 21, 47, 06, 0, time.UTC)}},
            },

            wantErr:  nil,
        },
        {
            name: "should fail with 404",
            args: args{
                branch:          "ft/test-branch",
                httpstatus:      404,
                owner:           "testowner",
                repo:            "testrepo",
                workflowFile:    "nonexistenttestfile.yaml",
            },
            endpoint: endpoint{
                branch:          "ft/test-branch",
                owner:           "testowner",
                repo:            "testrepo",
                workflowFile:    "nonexistenttestfile.yaml",
                runs: `{"total_count":0,"workflow_runs":[]}`,

            },
            wantRuns: []*github.WorkflowRun{},

            wantErr:  fmt.Errorf("Workflow not found"),
        },

    }

    for _, tt := range tests {

        // var apiurl string
        
        t.Run(tt.name, func(t *testing.T) {

            client, mux, _, teardown := Setup()
            defer teardown()

            ctx := context.Background()
                
            apiurl := fmt.Sprintf("/repos/%s/%s/actions/workflows/%s/runs", tt.endpoint.owner, tt.endpoint.repo, tt.endpoint.workflowFile)

            mux.HandleFunc(apiurl, func(w http.ResponseWriter, r *http.Request) {
                TestingMethod(t, r, "GET")
                if tt.args.httpstatus == 200 {
                    w.WriteHeader(http.StatusOK)
                } else if tt.args.httpstatus == 404 {
                    w.WriteHeader(http.StatusNotFound)
                }
                fmt.Fprint(w, tt.endpoint.runs)
            })
            
            gotRuns, gotErr := ReturnWorkflowRuns(tt.args.branch, ctx, client, tt.args.owner, tt.args.repo, tt.args.workflowFile, 20)

            fmt.Println(gotRuns)
            fmt.Println(gotErr)

            if tt.wantErr == nil {
                
                if gotErr != nil {
                    t.Errorf("ReturnWorkflowRuns() returned error: %v expect %v", gotErr, tt.wantErr)
                }

            } else if gotErr.Error() != tt.wantErr.Error() {
                
                t.Errorf("ReturnWorkflowRuns() returned error: %v expect %v", gotErr, tt.wantErr)
            }

            if len(gotRuns) != len(tt.wantRuns) {

                t.Errorf("expected %d elements but received %d instead", len(tt.wantRuns), len(gotRuns))
            }

            if !reflect.DeepEqual(gotRuns, tt.wantRuns){
                
				for i, run := range gotRuns {
                    
					if run != tt.wantRuns[i] {

                        if *run.ID != *tt.wantRuns[i].ID {
                            t.Errorf("ReturnWorkflowRuns() failed to retrieve same run on element %d - ID received %d but expects %d", i, *run.ID, *tt.wantRuns[i].ID)
                        }

                        if *run.Name != *tt.wantRuns[i].Name {
                            t.Errorf("ReturnWorkflowRuns() failed to retrieve same run on element %d - Name received %s but expects %s", i, *run.Name, *tt.wantRuns[i].Name)
                        }

                        if *run.NodeID != *tt.wantRuns[i].NodeID {
                            t.Errorf("ReturnWorkflowRuns() failed to retrieve same run on element %d - NodeID received %s but expects %s", i, *run.NodeID, *tt.wantRuns[i].NodeID)
                        }

                        if *run.RunNumber != *tt.wantRuns[i].RunNumber {
                            t.Errorf("ReturnWorkflowRuns() failed to retrieve same run on element %d - RunNumber received %d but expects %d", i, *run.RunNumber, *tt.wantRuns[i].RunNumber)
                        }

                        if *run.Event != *tt.wantRuns[i].Event {
                            t.Errorf("ReturnWorkflowRuns() failed to retrieve same run on element %d - Event received %s but expects %s", i, *run.Event, *tt.wantRuns[i].Event)
                        }

                        if *run.Status != *tt.wantRuns[i].Status {
                            t.Errorf("ReturnWorkflowRuns() failed to retrieve same run on element %d - Status received %s but expects %s", i, *run.Status, *tt.wantRuns[i].Status)
                        }

                        if *run.Conclusion != *tt.wantRuns[i].Conclusion {
                            t.Errorf("ReturnWorkflowRuns() failed to retrieve same run on element %d - Conclusion received %s but expects %s", i, *run.Conclusion, *tt.wantRuns[i].Conclusion)
                        }

                        if *run.CreatedAt != *tt.wantRuns[i].CreatedAt {
                            t.Errorf("ReturnWorkflowRuns() failed to retrieve same run on element %d - CreatedAt expects %v but received %v", i, *run.CreatedAt, *tt.wantRuns[i].CreatedAt)
                        }

                        if *run.UpdatedAt != *tt.wantRuns[i].UpdatedAt {
                            t.Errorf("ReturnWorkflowRuns() failed to retrieve same run on element %d - UpdatedAt received %s but expects %s", i, *run.UpdatedAt, *tt.wantRuns[i].UpdatedAt)
                        }

                    }
                }
            }

        })
    }

}
