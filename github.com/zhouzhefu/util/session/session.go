package session

import (
	"fmt"
	"time"
	"net/http"
	"crypto/rand"
	"encoding/base64"
	"sync"
)

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
	discardedSessionIds := make([]string, 10)
	for sessionId, session := range sm.store {
		if session.GC() {
			discardedSessionIds = append(discardedSessionIds, sessionId)
		}
	}

	for _, sessionId := range discardedSessionIds {
		delete(sm.store, sessionId)
	}

	time.AfterFunc(time.Duration(60) * time.Second, sm.GC)
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
	ses.ExpireAt = time.Now().Add(time.Duration(30) * time.Second)
}

func (ses *Session) Retouch() {
	ses.rwlock.Lock()
	defer ses.rwlock.Unlock()

	ses.ExpireAt = time.Now().Add(time.Duration(30) * time.Second)
}

func (ses *Session) IsExpired() bool {
	ses.rwlock.RLock()
	defer ses.rwlock.RUnlock()

	return time.Now().After(ses.ExpireAt)
}
/**
* Return "true" means the session satisfies the GC condition: 
* 1. IsExpired == true
* 2. ......
* And thus was successfully GC. 
*/
func (ses *Session) GC() bool {
	ses.rwlock.RLock()
	defer ses.rwlock.RUnlock()

	gcStatus := false

	isExpired := ses.IsExpired()
	if isExpired {
		ses.rwlock.Lock()

		if ses.Attributes != nil {
			ses.Attributes = nil
		} else {
			gcStatus = true
		}

		ses.rwlock.Unlock()
	}

	return gcStatus
}

