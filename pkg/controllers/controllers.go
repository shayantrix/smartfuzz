package controllers

import (
	"bufio"
	"fmt"
	"time"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var (
	ValidStatusCodes = map[int]bool{200: true, 301: true, 303: true, 307: true, 202: true, 201: true}
)

func Fuzz() {
	fmt.Println("Fuzzing target ...")
}
//Task: Add time base sql injection 
func URLHandler(urlString string) {
	payload := Payload("/home/shayan/Desktop/projects/smartfuzz/pkg/controllers/time_payload.txt")
	 //Sending a Get request; Having the response page HTML code
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	u, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Error in URL string: ", err)
		return
	}

	//resp, err := client.Get(u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		fmt.Println("Error in handling the URL: ", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:137.0) Gecko/20100101 Firefox/137.0")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in parsing the response body: ", err)
		return
	}
	// store html content
	htmlContent := string(body)

	//fmt.Println(htmlContent)
	//html parsing to have html.Node
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal("Could not parse the html content: ", err)
	}

	//finding html element nodes
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "input" || n.Data == "textarea" || n.Data == "form") {
			fmt.Println("Found: ", n.Data)
			for _, attr := range n.Attr {
				fmt.Printf(" - %s = %s\n", attr.Key, attr.Val)
				if attr.Key == "name" {

					//Here we kinda won the match, Cause its not controlled by JS
					//So we try to send requsts to this <input> tag
					//https://ililearn.com/?s=sda
					//here the name is s and the value is "sda"
					//func Post(url, contentType string, body io.Reader) (resp *Response, err error)
					//func Get(url string) (resp *Response, err error)
					//func PostForm(url string, data url.Values) (resp *Response, err error)

					for _, item := range payload {
						//Starting time to send payloads
						//will be used for time-based injection
						start := time.Now()

						//for some url encodings
						payloadVal := url.QueryEscape(item)
						
						fullUrl := fmt.Sprintf("%s?%s=%s", u.String(), attr.Val, payloadVal)

						response, err := http.Get(fullUrl)
						if err != nil {
							fmt.Println("Error in sending request: ", err)
						}
						if !ValidStatusCodes[response.StatusCode] {
							//Send POST request
							v := url.Values{}
							v.Set(attr.Key, item)

							startPost := time.Now()
							
							PostRequest, err := http.PostForm(u.String(), v)
							if err != nil {
								fmt.Println("Error in sending Post form request: ", err)
							}
							if  ValidStatusCodes[PostRequest.StatusCode]{
								fmt.Println("Found !!!!")
								elapsedPost := time.Since(startPost)
								if elapsedPost > 5 * time.Second {
									fmt.Println("Could be vulnerable to sql injection time-base attack!!!")
								}
								fmt.Println(item)
								//body_req, err := io.ReadAll(PostRequest.Body)
								//if err != nil{
								//	log.Fatal("Error in parsing response body: ", err)
								//}
								//fmt.Println(string(body_req))
								fmt.Println(PostRequest.StatusCode)
							}else{
								fmt.Println("No potential SQL Injection may be found!")
							}
						} else {
							elapse := time.Since(start)

							if elapse > 5 * time.Second{
								fmt.Println("Could be vulnerable to sql injection time-base attack!!!")
							}
							fmt.Println("Found !!!!")
							fmt.Println(item)
							//body_request, err := io.ReadAll(response.Body)
							//if err != nil{
							//	log.Fatal("Error in parsing response body: ", err)
							//}

							//fmt.Println(string(body_request))
							fmt.Println(response.StatusCode)
						}
					}

				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
                                f(c)
                }
	}

	f(doc)

}

func Payload(path string) []string {
	//os.ReadFile(name string)
	//Open(file string) (io.Reader, error)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error in reading the payload file: ", err)
	}

	defer file.Close()

	//NewScanner(io.Reader)
	scanner := bufio.NewScanner(file)

	var payload []string

	for scanner.Scan() {
		//scanner.Text()
		payload = append(payload, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return payload
}
