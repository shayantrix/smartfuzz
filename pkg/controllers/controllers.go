package controllers

import (
        "fmt"
        "net/http"
        "io"
        "strings"
        "log"
        "golang.org/x/net/html"
)

func Fuzz(){
        fmt.Println("Fuzzing target ...")
}

func URLHandler (url string){
        // Sending a Get request; Having the response page HTML code
        resp, err := http.Get(url)
        if err != nil {
                fmt.Println("Error in handling the URL: ", err)
                return
        }

        defer resp.Body.Close()

        body, err := io.ReadAll(resp.Body)
        if err != nil{
                fmt.Println("Error in parsing the response body: ", err)
                return
        }
        // store html content 
        htmlContent := string(body)

        //html parsing to have html.Node
        doc, err := html.Parse(strings.NewReader(htmlContent))
        if err != nil{
                log.Fatal("Could not parse the html content: ", err)
        }

        //finding html element nodes
        var f func(*html.Node)
        f = func(n *html.Node){
                if n.Type == html.ElementNode && (n.Data == "input" || n.Data == "textarea" || n.Data == "form"){
                        fmt.Println("Found: ", n.Data)
                        for _, attr := range n.Attr{
                                fmt.Printf(" - %s = %s\n", attr.Key, attr.Val)
                        }
                }
                for c := n.FirstChild ; c != nil ; c = c.NextSibling {
                        f(c)
                }
        }

        f(doc)
}
