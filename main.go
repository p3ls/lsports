package main

import (
    "fmt"
    "log"
    "os/exec"
    "bytes"
    "strings"
    "github.com/gosuri/uitable"
    // "io"
)



func netstat_info() [][]string {
    cmd := exec.Command("netstat", "-4", "--numeric", "--all")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
    lines := strings.Split(out.String(), "\n")
    results := make([][]string, 0)
    for _, item := range lines {
        if (strings.HasPrefix(item, "tcp") || strings.HasPrefix(item, "udp")) {
            fields := strings.Fields(item) 
            proto := fields[0]
            local_ip_port := fields[3]
            port := strings.Split(local_ip_port, ":")[1]
            details := fuser_info(proto, port)
            user := details[1]
            // pid := details[2]
            process := details[4]
            results = append(results, []string{proto, port, process, user})
            // fmt.Printf("%s %s %s\n", proto, port, process)
           
        }
        
    }
    return results
}

func fuser_info(proto string, port string) []string {
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

func main() {

    table := uitable.New()  
    table.MaxColWidth = 50
    table.AddRow("PROTO", "PORT", "PROCESS", "USER")

    list_of_list_of_string := netstat_info()
    for _, v := range list_of_list_of_string {
        // fmt.Printf("%+q\n", v)
        table.AddRow(v[0], v[1], v[2], v[3])
    }
    fmt.Println(table)

}

