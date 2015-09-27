package main 

import (
	"fmt"
	"encoding/xml"
	"io/ioutil"
	"os"

	"html/template"

	"strings"
	"strconv"
)

type RecurlyServers struct {
	XMLName xml.Name `xml:"servers"`
	Version string `xml:"version,attr"`
	Servers []Server `xml:"server"`
	Description string `xml:",innerxml"`

	Desc string `xml:"anDesc>desc"`
	AnyInfo string `xml:",any"`
}

type Server struct {
	XMLName xml.Name `xml:"server"`
	ServerName string `xml:"serverName"`
	ServerIP string `xml:"serverIP"`
}

func readXML() {
	file, err := os.Open("xmlJson/Servers.xml")
	if err != nil {
		fmt.Println("Error when Opening Servers.xml")
		return
	} 
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error when reading Servers.xml")
		return
	}

	serversXml := RecurlyServers{}
	err = xml.Unmarshal(bytes, &serversXml)
	if err != nil {
		fmt.Println("Error when Unmarshal")
		return
	}

	fmt.Println(serversXml)
	fmt.Println(serversXml.Version, serversXml.Description, serversXml.Desc)
	fmt.Println(serversXml.Servers)
	for _, node := range serversXml.Servers {
		fmt.Println("ServerIP: ", node.ServerIP)
	}

	fmt.Println(serversXml.AnyInfo)
}

type Friend struct {
	FriendName string
	Age int
}

type Person struct {
	Name string
	Emails []string
	Friends []*Friend
}

func dealWithEmail(s string) string {
	return "<" + s + ">"
}

type MyHtml struct {
	Header string
	Body string
	Footer string
}

func readTemplate() {
	t := template.New("hello") //imagine there is a map[string]*Template, can use Lookup(name) to find it out
	t, _ = t.Parse("Hello {{.Name}}\n")
	t.Execute(os.Stdout, Person{Name:"Jeremy"})

	//if-else, with, range, self-defined func...
	t1 := template.New("tree")
	t1 = t1.Funcs(template.FuncMap{"dealEmail": dealWithEmail})
	t1, _ = t1.Parse(`
			Hello {{.Name}}!
			Your emails: 
			{{range .Emails}}
				an email: {{.}} {{. | dealEmail}}
			{{end}}

			Your friends: 
			{{with .Friends}}
				{{range .}}
					{{if .Age}}{{.FriendName}}{{end}}
				{{end}}
			{{end}}
		`)

	f1 := Friend{FriendName: "Ah Mao", Age: 13}
	f2 := Friend{FriendName: "Ah Gau", Age: 19}
	f3 := Friend{FriendName: "Ah Jiu"}
	p1 := Person{
		Name: "Jeremy1", 
		Emails: []string{"jeremy@yahoo.com", "jeremy@gmail.com", "jeremy@163.com"},
		Friends: []*Friend{&f1, &f2, &f3},
	}

	t1.Execute(os.Stdout, p1)


	// reusable template
	myHtml := MyHtml{Header: "MyHeader", Body: "MyBody", Footer: "MyFooter"}
	s1, _ := template.ParseFiles("templates/header.tmpl", "templates/body.tmpl", "templates/footer.tmpl")

	// why embeded templates cannot bind the data fields?
	s1.ExecuteTemplate(os.Stdout, "header", myHtml)
    fmt.Println("header done.")
    s1.ExecuteTemplate(os.Stdout, "body", myHtml)
    fmt.Println("body done.")
    s1.ExecuteTemplate(os.Stdout, "footer", myHtml)
    fmt.Println("footer done. ")
    s1.Execute(os.Stdout, myHtml)
    fmt.Println("full done. ")
}

func processString() {
	str := make([]byte, 0, 100)
	// str := ""
	fmt.Println(str)
	str = strconv.AppendInt(str, 4567, 10)
	str = strconv.AppendBool(str, false)
	str = strconv.AppendQuote(str, "abcde")
	str = strconv.AppendQuoteRune(str, '周')

	fmt.Println(str)
	fmt.Println(string(str))

	str1 := strconv.FormatBool(false)
	fmt.Println(str1)

	str1 = strings.Repeat(str1, 2)
	fmt.Println(str1)
	fmt.Println(strings.Contains(str1, "al"))
	fmt.Println(strings.Index(str1, "al"))
	fmt.Println(strings.Trim("!a     james  May       !a", "!a"))
}

func main() {
	readXML()
	
	fmt.Println()
	fmt.Println()

	readTemplate()

	fmt.Println()
	fmt.Println()

	processString()
}