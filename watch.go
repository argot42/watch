package main

import (
    "fmt"
    "io"
    "os"
    "os/exec"
    "os/signal"
    "flag"
    "syscall"
    "strings"
    "github.com/argot42/watcher"
)

var prog = flag.String("p", "", "program to be executed")

func main() {
    flag.Usage = usage
    flag.Parse()

    if flag.NArg() < 1 || *prog == "" {
        flag.Usage()
    }

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGTERM)

    subscription := watcher.Watch(flag.Args()[0])
    command := strings.Split(*prog, " ")
    first := true

    End:
    for {
        select {
        case <-subscription.Out:
            if !first {
                if err := Run(command, os.Stdout); err != nil {
                    fmt.Fprintf(os.Stderr, "%s error: %s\n", os.Args[0], err)
                    /* exec.ExitError are not treated as fatal errors since
                     * we don't want the program to end when a program exit status is 1
                     * on bad input */
                    if _, ok := err.(*exec.ExitError); !ok {
                        os.Exit(1)
                    }
                }
            }
            first = false
        case err := <-subscription.Err:
            fmt.Fprintf(os.Stderr, "%s error: %s\n", os.Args[0], err)
            os.Exit(1)
        case <- sigs:
            break End
        }
    }
}

func usage() {
    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr, "\n%s -p cmd file\n", os.Args[0])
    os.Exit(1)
}

func Run(command []string, out io.Writer) (err error) {
    cmd := exec.Command(command[0], command[1:]...)

    cmdout, err := cmd.StdoutPipe()
    if err != nil {
        return
    }

    if err = cmd.Start(); err != nil {
        return
    }

    if _, err = io.Copy(out, cmdout); err != nil {
        return
    }

    err = cmd.Wait()

    return
}
