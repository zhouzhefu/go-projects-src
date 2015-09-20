package main 

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	//"strings"
	//"md5"
	"os"
	"io"
	//"crypto/rand"
	//"encoding/base64"
	//"time"
	"github.com/zhouzhefu/util/session"
)

/**
* As per form submitted with token, we can validate whether uer submit same thing multiple times
*/
func withToken() string {
	//h := md5.New()
	//token := fmt.Sprintf("%x", h.Sum(nil))
	token := "dummyToken"
	return token
}

func sayHallo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Web")
}


/**
* this pattern of request process, isn't it similar to an abused JSP code? if (method == "GET") {..} else {..}
*/
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("client Method: ", r.Method)
	if r.Method == "GET" { //GET means user just reach login panel
		session, _ := glbSess.CreateOrUpdateSession(w, r)
		fmt.Println("GET to retouch session:", session)

		t, _ := template.ParseFiles("login.gtpl")
		//t.Execute(w, nil)
		t.Execute(w, withToken())
	} else { //POST means user try to login
		r.ParseForm() //by default form will not be parsed until call out, 
		fmt.Println("username: ", r.Form["username"][0]) //only after ParseForm() was called, 
		fmt.Println("password: ", r.Form["password"]) //these fields can read value

		//validate token (usually we use session store & compare)
		//token := r.Form["token"] //Form[field] result is []string
		token := r.FormValue("token") //or r.Form["token"][0]
		if token != "" {
			fmt.Println("token: ", token, "submitted")
		} else {
			fmt.Println("Aiyo no token!")
		}

		//check session
		session, _ := glbSess.CreateOrUpdateSession(w, r)
		currUsrName, exists := session.Attributes["username"]
		if !exists || session.IsExpired() {
			currUsrName = r.Form.Get("username")
			session.Attributes["username"] = currUsrName
		} else {
			fmt.Println("Current you have been login as:", currUsrName)
		}
fmt.Println("hit-3")

		//output to page should be escaped in case of injection attack
		template.HTMLEscape(w, []byte("Welcome " + currUsrName))
	}

	gosessionId, _ := r.Cookie("gosessionid")
	fmt.Println("Your gosessionid is:", gosessionId.Value)
	fmt.Println("Current session object is:", glbSess.GetSession(gosessionId.Value))
	fmt.Println("global session:", glbSess)
}

/**
* handle file upload
*/
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("client Method: ", r.Method)

	if r.Method == "GET" {
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, withToken())
	} else {
		//ParseMultipartForm calls ParseForm if necessary, including parsing those non-multipart fields
		err := r.ParseMultipartForm(32 << 20) //maxMemory = 32 * 2^20 bytes
		if err != nil {fmt.Println("Aiyo error lah:", err)}

		file, handler, err := r.FormFile("uploadFile")
		if err != nil {
			fmt.Println("What happen?", err)
			return
		}
		defer file.Close() //finally {in.close()} in Java

		fmt.Fprintf(w, "%v", handler.Header)

		//prepare the space to copy the uploaded file to
		targetFile, err := os.OpenFile(
			"./testUpload/" + r.FormValue("renameTo"), 
			os.O_WRONLY|os.O_CREATE, 
			os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer targetFile.Close()

		io.Copy(targetFile, file)

	}
}


var glbSess session.SessionManager
func main() {
	glbSess.Init()
	go glbSess.GC()
	//add route/handler
	http.HandleFunc("/", sayHallo) //in this way, sayHallo() will be called every time when /login is accessed
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)

	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}