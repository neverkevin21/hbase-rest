# hbase-rest

Usege:

```
import (
    "fmt"
    "github.com/neverkevin21/hbase-rest"
)

func main() {
    headers := make(map[string]string)
    headers["auth"] = "authentication"
    timeout := 10
    client := rest.NewRest("http://127.0.0.1:8088", timeout, headers)

    // get
    row, err := client.Get("<table>", "<rowkey>", "")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(row)

    // put
    err = client.Put("<table>", "<rowkey>", "<family>:<column>", "<value>")
    if err != nil {
        fmt.Println(err)
    }
}
```
