package main

import (
    "fmt"
    "log"
    "os/exec"
    "bytes"
    "sort"
    "strings"
    "runtime"
    "github.com/gosuri/uitable"
    "menteslibres.net/gosexy/to"
    // "io"
)

// windows support
// https://commandwindows.com/netstat.htm
// https://stackoverflow.com/questions/18059798/windows-batch-equivalent-of-fuser-k-folder

var version = "HEAD"

func sliceContains(a []string, x string) bool {
    for _, n := range a {
        if x == n {
            return true
        }
    }
    return false
}

func trimSuffix(s, suffix string) string {
    if strings.HasSuffix(s, suffix) {
        s = s[:len(s)-len(suffix)]
    }
    return s
}

func fuserInfo(proto string, port string) []string {

    proto = trimSuffix(proto, "6")    
    cmd := exec.Command("sudo", "fuser", "-a", "-v", fmt.Sprintf("%s/%s", port, proto))
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
    lines := strings.Split(out.String(), "\n")
    for _, item := range lines {
        if strings.HasPrefix(item, port) {
            fields := strings.Fields(item)            
            // fmt.Printf("%+q\n", fields)
            return fields
        }
    }
    return []string{}
}

func netstatInfo() [][]string {
    cmd := exec.Command("netstat", "-4", "-6", "--numeric", "--all")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
    lines := strings.Split(out.String(), "\n")
    results := make([][]string, 0)
    string_results := make([]string, 0)
    for _, item := range lines {

        if (strings.HasPrefix(item, "tcp") || strings.HasPrefix(item, "udp")) {
            fields := strings.Fields(item) 
            proto := trimSuffix(fields[0], "6")
            local_ip_port := fields[3]
            addr := strings.Split(local_ip_port, ":")
            port := addr[len(addr)-1]
            details := fuserInfo(proto, port)
            user := details[1]
            process := details[4]
            // pid := details[2]

            id := proto + port + process + user
            if !sliceContains(string_results, id){
                results = append(results, []string{proto, port, process, user})
                string_results = append(string_results, id)
            }
           
        }
        
    }
    return results
}

func main() {

    if runtime.GOOS != "linux" {
        log.Fatal("Only Linux is supported")
    }

    table := uitable.New()  
    table.MaxColWidth = 50
    table.AddRow("PORT", "PROTO", "PROCESS", "USER")

    items := netstatInfo()

    // sort items by port number ascending
    sort.SliceStable(items, func(i, j int) bool {
        porti := to.Int64(items[i][1])
        portj := to.Int64(items[j][1])
        return porti < portj
    })

    for _, v := range items {
        table.AddRow(v[1], v[0], v[2], v[3])
    }

    fmt.Println(table)

}

