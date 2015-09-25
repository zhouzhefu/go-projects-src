package session

import (
	"fmt"
	"time"
	"net/http"
	"crypto/rand"
	"encoding/base64"
	"sync"
)

func init() {
	fmt.Println("init(): This func will be executed when this package was imported for side-effects. ")
}

func genSessionId() string {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}


type SessionManager struct {
	store map[string]Session
	rwlock sync.RWMutex
}

func (sm *SessionManager) Init() {
	sm.store = make(map[string]Session)
}

/**
* The author of the tutorial suggest add sync.Mutex.Lock() and defer sync.Mutex.Unlock(), for every 
* SessionManager operation on session, read/write. I doubt this a good practice, with which the whole App 
* has a global single Lock/syncronizer, and performance will be compromised. Since all the request handling 
* code are limited to its own unique sessionId, they will work fine together without Locking. 
* 
* To be considered is the goroutine for session GC, but we may redesign the GC logic to limit the Lock is on 
* a specific Session item (to be kept or to be GC). 
* 
* Now we have more thinking about this point, please read doc on IsExpired() method
*/
func (sm *SessionManager) CreateOrUpdateSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	sm.rwlock.Lock()
	defer sm.rwlock.Unlock()

	cookie, err := r.Cookie("gosessionid")
	if err != nil {
		fmt.Println(err)
		//return
	}
	var sessionId string
	if cookie == nil || cookie.Value == "" {
		sessionId = genSessionId()
		http.SetCookie(w, &http.Cookie{Name:"gosessionid", Value:sessionId, Path:"/"})
	} else {
		sessionId = cookie.Value
	}
fmt.Println("hit-1")
	//re-create session if last session expire or not exists, otherwise Retouch() to keep session alive
	session, exists := sm.store[sessionId]
	if !exists || session.IsExpired() {
		session = Session{}
		session.Init()
		sm.store[sessionId] = session
	} else {
		session.Retouch()
	}
fmt.Println("hit-2")
	return &session, nil
}

func (sm *SessionManager) GetSession(sessionId string) *Session {
	sm.rwlock.RLock()
	defer sm.rwlock.RUnlock()

	session, exists := sm.store[sessionId]
	if !exists {
		return nil
	} else {
		return &session
	}
}

func (sm *SessionManager) GC() {
	fmt.Println("SessionManager GC()!")
	discardedSessionIds := make([]string, 10)
	for sessionId, session := range sm.store {
		fmt.Println("try GC in sm:", session, sessionId, session.GC())
		if session.GC() {
			discardedSessionIds = append(discardedSessionIds, sessionId)
		}
	}

	fmt.Println("Session to be removed:", discardedSessionIds)
	for _, sessionId := range discardedSessionIds {
		delete(sm.store, sessionId)
	}

	time.AfterFunc(time.Duration(5) * time.Second, sm.GC)
}



type Session struct {
	Attributes map[string]string //how to store value as Object in Java?
	ExpireAt time.Time
	rwlock sync.RWMutex
}

func (ses *Session) Init() {
	ses.rwlock.Lock()
	defer ses.rwlock.Unlock()

	ses.Attributes = make(map[string]string)
	ses.ExpireAt = time.Now().Add(time.Duration(10) * time.Second)
}

func (ses *Session) Retouch() {
	ses.rwlock.Lock()
	defer ses.rwlock.Unlock()

	ses.ExpireAt = time.Now().Add(time.Duration(10) * time.Second)
}

/**
* Here we have to remove the RLock() & RUnlock() pair from the method, quite a few points worth to say: 
* 1. Unlike Java ReentrantReadWriteLock, RWMutex in Go is NOT reentrant. Therefore, you should never 
*    allow a Lock pair contain another one. Two consecutive Locker.Lock()/RLock() will deadlock the 
*    context goroutine. 
* 2. "No reentrant" applies between Lock() and RLock(). Java allows Read lock to be contained by a Write 
*    Lock, but in Go they are still exclusive to each other. 
* 3. Even if you can handle concerns above, you should still avoid using RWMutex down to small units, they 
*    are reusable and highly possible to be contained by other methods, which may also need to be protected 
*    by Lock, worsely some of them are in if..else.. block, everywhere is the trap of deadlock in runtime. 
* 4. Conclusion? Always try to avoid using session, if must, design it as small/lightweight as possible. As 
*    we have seen to handle session the Lock is needed, it is the true killer of concurrent performance. 
*/
func (ses *Session) IsExpired() bool {
	//ses.rwlock.RLock()
	//defer ses.rwlock.RUnlock()

	return time.Now().After(ses.ExpireAt)
}
/**
* Return "true" means the session satisfies the GC condition: 
* 1. IsExpired == true
* 2. ......
* And thus was successfully GC. 
*/
func (ses *Session) GC() bool {
	ses.rwlock.Lock()
	defer ses.rwlock.Unlock()

	gcStatus := false

	isExpired := ses.IsExpired()
	if isExpired {
		if ses.Attributes != nil {
			ses.Attributes = nil
		} else {
			gcStatus = true
		}
	}

	return gcStatus
}

