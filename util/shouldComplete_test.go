package util

import (
    "reflect"
    "testing"
)

func TestShouldComplete(t *testing.T){

    tests := []struct {
        name               string
        previousRunStatus  string
        wantShouldComplete bool   
    }{
        {
            name: "Return true - last workflow status is 'completed'",
            previousRunStatus: "completed",
            wantShouldComplete: true,
        },
        {
            name: "Return false - last workflow status is 'in_progress'",
            previousRunStatus: "in_progress",
            wantShouldComplete: false,
        },
        {
            name: "Return true - last workflow status is 'queued'",
            previousRunStatus: "queued",
            wantShouldComplete: false,
        },

    }

    for _, tt := range tests {

        t.Run(tt.name, func(t *testing.T) {

            gotShouldComplete := ShouldComplete(tt.previousRunStatus)

            if !reflect.DeepEqual(gotShouldComplete, tt.wantShouldComplete){

                t.Errorf("ShouldComplete() failed - ID expects %t but received %t", tt.wantShouldComplete, gotShouldComplete)

            }

        })
    }

}
