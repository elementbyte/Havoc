package logr

import (
    "bufio"
    "log"
    "os"
    "regexp"

    "github.com/Cracked5pider/Havoc/teamserver/pkg/logger"
)


func strip(str []byte) []byte {
    var (
    	ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
        re   = regexp.MustCompile(ansi)
    )

    return []byte(re.ReplaceAllString(string(str), ""))
}

func (l Logr) ServerStdOutInit() {
    var (
    	PathStdOut  = l.Path + "/teamserver.log"
        OldStdout   = os.Stdout

        StdRead, StdWrite, _  = os.Pipe()
    )

    File, err := os.OpenFile(PathStdOut, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    os.Stdout = StdWrite

    logger.LoggerInstance = logger.NewLogger(StdWrite)

    go func() {
        var Reader = bufio.NewReader(StdRead)

        for {
            if Reader.Size() > 0 {
                line, _, _ := Reader.ReadLine()
                line = []byte(string(line) + "\n")

                _, err := File.Write(strip(line))
                if err != nil {
                    return
                }

                _, err = OldStdout.Write(line)
                if err != nil {
                    return
                }
            }
        }
    }()
}
