package main

import (
	"bytes"
	"os/exec"
	"testing"
)

type RunTC struct {
    Cmd []string
    Out string
    Fail bool
}

func TestRun(t *testing.T) {
    tcs := []RunTC {
        {
            []string{"echo", "hello"},
            "hello\n",
            false,
        },
        {
            []string{"false"},
            "",
            true,
        },
        {
            []string{"(exit 1)"},
            "",
            false,
        },
    }

    for _, tc := range tcs {
        buf := new(bytes.Buffer)
        err := Run(tc.Cmd, buf)

        if err != nil {
            t.Logf("%s -> %s -> %T", tc.Cmd, err, err)
            // count ExitErrors as non fatal since some commands will fail with
            // the wrong input
            if _, ok := err.(*exec.ExitError); !ok {
                if tc.Fail {
                    t.Errorf("%s ExitError is not treated as fatal error", tc.Cmd)
                    continue
                }
            } else {
                if !tc.Fail {
                    t.Errorf("Command %s shouldn't have failed", tc.Cmd)
                }
                continue
            }
        }

        if tc.Fail {
            t.Errorf("Command %s should have failed", tc.Cmd)
            continue
        }

        if tc.Out != buf.String() {
            t.Errorf("%s got [%s] should have been [%s]", tc.Cmd, buf.String(), tc.Out)
            continue
        }
    }
}
