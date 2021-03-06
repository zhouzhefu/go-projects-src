package main 

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"strings"
	"os"
	"io"
	//"crypto/rand"
	//"encoding/base64"
	. "github.com/zhouzhefu/util/session" // please note Go has 3 ways of import
	_ "github.com/zhouzhefu/util/session" //Not a real import, just to execute the init() func of that package

	"time"
	"net"
	"strconv"

	"code.google.com/p/go.net/websocket"

	"net/rpc"

	"crypto/md5"
	"crypto/sha256"

	"errors"
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

		//output to page should be escaped in case of injection/CRSF attack
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


var glbSess *SessionManager

func main() {
	// startSessionServer()

	// startTcpServer()

	// startWebSocketServer()

	// startRpcServer()

	// startEncrptLoginServer()

	// startErrorServer()
	startDesignedErrorServer()
}

func startDesignedErrorServer() {
	http.Handle("/errorCode", AppHandler(fineErrorView))

	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type AppHandler func(w http.ResponseWriter, r *http.Request) *MyError
func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// default rule to handle errors
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
func fineErrorView(w http.ResponseWriter, r *http.Request) *MyError {
	return &MyError{time.Now(), "This is from fineErrorView() error"}
}

func startErrorServer() {
	http.HandleFunc("/error1", func(w http.ResponseWriter, r *http.Request) {
		result, err := generateError1()
		fmt.Println(result, err)

		//throw 500 error to browser
		http.Error(w, err.Error(), 500)
	})
	http.HandleFunc("/error2", func(w http.ResponseWriter, r *http.Request) {
		_, err := generateError2()

		// This "Type Assertion" equivalent to (as in Java): 
		// myErr = (MyError) err;
		// correct = err instanceof MyError;
		if myErr, correct := err.(MyError); correct {
			fmt.Println("It is really a MyError!", myErr)
		} else {
			fmt.Println("Sorry it is not that MyError. ", err)
		}
		
	})
	// Please carefully identify the subtle difference between /error2 & /error3. 
	// BOTH pointer & value receiver of Error(), makes Go believe "error" interface has been implemented. 
	http.HandleFunc("/error3", func(w http.ResponseWriter, r *http.Request) {
		_, err := generateError3()

		// This "Type Assertion" equivalent to (as in Java): 
		// myErr = (*MyErrorP) err;
		// correct = err instanceof *MyErrorP;
		if myErr, correct := err.(*MyErrorP); correct {
			fmt.Println("It is really a MyErrorP!", myErr)
		} else {
			fmt.Println("Sorry it is not that MyErrorP. ", err)
		}
		
	})

	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func generateError1() (string, error) {
	return "error1 tested", errors.New("error1 here")
}

type MyError struct {
	When time.Time
	What string
}
func (me MyError) Error() string {
	return me.What + "@" + me.When.String()
}
func generateError2() (string, error) {
	return "", MyError{When: time.Now(), What: "MyError2 here"}
}

type MyErrorP struct {
	When time.Time
	What string
}
func (me *MyErrorP) Error() string {
	return me.What + "@" + me.When.String()
}
func generateError3() (string, error) {
	return "", &MyErrorP{When: time.Now(), What: "MyError3P here"}
}

func startEncrptLoginServer() {
	http.HandleFunc("/login", encryptLogin)

	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// submitted password will be hashed(one-way) to avoid storing plain-text password
func encryptLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { //GET means user just reach login panel
		t, _ := template.ParseFiles("login.gtpl")
		//t.Execute(w, nil)
		t.Execute(w, "dummyToken")
	} else { //POST means user try to login
		r.ParseForm()
		username := r.FormValue("username")
		plainPassword := r.FormValue("password")
		fmt.Println("username: ", username)
		fmt.Println("(plain)password: ", plainPassword)

		hashMd5 := md5.New()
		passwordMd5 := hashMd5.Sum([]byte(plainPassword))

		hashSha256 := sha256.New()
		io.WriteString(hashSha256, plainPassword)
		passwordSha256 := hashSha256.Sum(nil)

		fineEncryptedPassSHA256 := fmt.Sprintf("%x", passwordSha256)
		fineEncryptedPassMD5 := fmt.Sprintf("%x", passwordMd5)
		fmt.Println("MD5 Password:", fineEncryptedPassMD5, "SHA256 Password:", fineEncryptedPassSHA256)
	}
}

func startRpcServer() {
	rpc.Register(new(Arith1)) //registered as "<ReceiverTypeName>.<MethodName>" for client Call()
	rpc.HandleHTTP()

	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type Args struct {
	A, B int
}
type Arith1 int
func (a *Arith1) Multiply(args Args, ret *int) error {
	*ret = args.A * args.B
	return nil
}

func startWebSocketServer() {
	http.Handle("/websocket", websocket.Handler(EchoWebSocket))

	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func EchoWebSocket(wsConn *websocket.Conn) {
	for {
		var reply string
		if err := websocket.Message.Receive(wsConn, &reply); err != nil {
			fmt.Println("Cannot Receive!")
			break
		}

		fmt.Println("Receive from client:", reply)

		msg := "Received: " + reply
		fmt.Println("Sending client message:", msg)

		if err := websocket.Message.Send(wsConn, msg); err != nil {
			fmt.Println("Cannot Send!")
			break
		}
	}
}

// UDP is very similar so won't create extra examples
func startTcpServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	checkError(err)
	fmt.Println("tcpAddr:", tcpAddr)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for { // event loop while(true)
		conn, err := listener.Accept()
		if err != nil { 
			checkError(err)
			continue; //of course, you don't want to stop the whole server just because of one conn error
		}

		//processing logic, here is just a simple timestamp as response
		// handleTcpClient(conn)
		handleLongConnTcpClient(conn)
	}
}

func handleLongConnTcpClient(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(2 * time.Second)) //timeout in 2 min
	defer conn.Close()

	connId := time.Now().String()

	for {
		fmt.Println("Ready to read request:", connId)
		request := make([]byte, 128) // larger than 128 bytes request will be identified as flood attack
		readLen, err := conn.Read(request)
		if err != nil {
			fmt.Println("Error found on server:", err)
			// this error is booked by conn.SetReadDeadLine(), without this if block, 
			// the connection will never get Close() even if timeout already hit. 
			if strings.Contains(err.Error(), "i/o timeout") {break} 
			// this error is possible when client Close() the connection, without this if block, 
			// the connection will never get Close() even if EOF already hit
			if strings.Contains(err.Error(), "EOF") {break}

			continue
		}
		fmt.Println("Processing conn:", connId, "\nreadLen:", readLen)

		if readLen == 0 {
			continue
		} else if string(request) == "timestamp" {
			daytime := strconv.FormatInt(time.Now().Unix(), 10) + " @" + connId
			fmt.Println("Writing daytime:", daytime)
			conn.Write([]byte(daytime))
		} else {
			daytime := time.Now().String() + " @" + connId
			fmt.Println("Writing daytime:", daytime)
			conn.Write([]byte(daytime))

			
		}
	}

	conn.Close()
	fmt.Println("End of the connection. ")
}

func handleTcpClient(conn net.Conn) {
	time.AfterFunc(time.Duration(2) * time.Second, func() {
			time := time.Now().String()
			conn.Write([]byte(time))
			conn.Close()		
		})
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(0)
	}
}

func startSessionServer() {
	glbSess = new(SessionManager)
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