package controllers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var (
	ValidStatusCodes = map[int]bool{200: true, 301: true, 303: true, 307: true, 202: true, 201: true}

        commandInjectionFound = map[string]bool{"uid": true, "gid": true, "groups": true, "Permission denied": true, "root": true, "bin": true, "sys": true, "www-data": true, "dhcpcd": true, "syslog": true, "uuidd": true, "daemon": true, "LISTEN": true, "No such file or directory": true, "icmp_seq": true}

)

func Fuzz() {
	fmt.Println("Fuzzing target ...")
}

func CommandInjection(urlString string) {
	payload := Payload("/home/shayan/Desktop/projects/smartfuzz/pkg/controllers/command_injection_payload.txt")

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
	if err != nil{
		fmt.Println("Could not send request to web application: ", err)
	}

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

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "input" || n.Data == "textarea" || n.Data == "form"){
			fmt.Println("\nFound: ", n.Data)
			for _, attr := range n.Attr {
				fmt.Printf(" - %s = %s\n", attr.Key, attr.Val)

				var attribute_val string
				//command injection
				if attr.Key == "action"{
					attribute_val = attr.Val
				}
				if attr.Key == "method" && attr.Val == "post" {
					for _, command := range payload{

							v := url.Values{}
                                                        v.Set(attr.Key, command)


							response, err := http.PostForm(u.String(), v)
							if err != nil{
								fmt.Println("Error in sendig the GET request")
							}

							body, err := io.ReadAll(response.Body)
							if err != nil{
								fmt.Println("Error in reading response's body: ",err)
							}

							//fmt.Println(string(body))
							bodyStr := string(body)
							for key := range commandInjectionFound{
								if strings.Contains(bodyStr, key){
									fmt.Println("{")
									fmt.Println("Command Injection Possibility")
									fmt.Println(command)
									fmt.Println("}")
								}
							}

						}
				}

				if attr.Key == "method" && attr.Val == "get"{
					for _, command := range payload{

                                                payloadVal := url.QueryEscape(command)
                                                parseUrl, err := url.Parse(attribute_val)
                                                if err != nil{
                                                        fmt.Println("Error in parsing the URL: ",err)
                                                        return
                                                }
                                                params := parseUrl.Query()
                                                for key := range params{
                                                        fullUrl := fmt.Sprintf("%s?%s=%s", u.String(), key, payloadVal)

                                                        response, err := http.Get(fullUrl)
                                                        if err != nil{
                                                                fmt.Println("Error in sendig the GET request")
                                                        }

							body, err := io.ReadAll(response.Body)
                                                        if err != nil{
                                                                fmt.Println("Error in reading response's body: ",err)
                                                        }

                                                        //fmt.Println(string(body))
							bodyStr := string(body)
                                                       	for key := range commandInjectionFound{
                                                                if strings.Contains(bodyStr, key){
                                                                        fmt.Println("{")
                                                                        fmt.Println("Command Injection Possibility")
                                                                        fmt.Println(command)
                                                                        fmt.Println("}")
                                                                }
                                                        }

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

// Task: Add time base sql injection
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
			fmt.Println("\nFound: ", n.Data)
			for _, attr := range n.Attr {
				fmt.Printf(" - %s = %s\n", attr.Key, attr.Val)
				if attr.Key == "name" {
					for _, item := range payload {
						//if the method is utilizing Get request
						//if attr.Key == "method" && attr.Val == "GET" {
							//Starting time to send payloads
							start := time.Now()

							//for some url encodings
							payloadVal := url.QueryEscape(item)

							fullUrl := fmt.Sprintf("%s?%s=%s", u.String(), attr.Val, payloadVal)

							response, err := http.Get(fullUrl)
							if err != nil {
								fmt.Println("Error in sending request: ", err)
							}
							elapse := time.Since(start)
							fmt.Println("{\n")
							fmt.Println("Found !!!!")
							fmt.Println(item)
							if elapse > 5*time.Second {
								fmt.Println("Could be vulnerable to sql injection time-base attack!!!")
							}
							fmt.Println("Response Status code: ", response.StatusCode)
							fmt.Println("\n}\n")

						//} else if attr.Key == "method" && attr.Val == "POST" {
							//Send POST request
							v := url.Values{}
							v.Set(attr.Key, item)

							startPost := time.Now()

							PostRequest, err := http.PostForm(u.String(), v)
							if err != nil {
								fmt.Println("Error in sending Post form request: ", err)
							}
							if ValidStatusCodes[PostRequest.StatusCode] {
								fmt.Println("{\n")
								fmt.Println("Found !!!!")
								elapsedPost := time.Since(startPost)
								fmt.Println(item)
								if elapsedPost > 5*time.Second {
									fmt.Println("Could be vulnerable to sql injection time-base attack!!!")
								}
								fmt.Println("Response Status code: ", PostRequest.StatusCode)
								fmt.Println("\n}\n")
							}

							// if !ValidStatusCodes[PostRequest.StatusCode] {
							// 	fmt.Println("No potential SQL Injection may be found!")
							// 	fmt.Println(PostRequest.StatusCode)
							// }
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
