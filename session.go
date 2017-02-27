package ggweb

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	ID         string
	LastAccess time.Time
	Values     map[string]interface{}
}
type SessionMgr struct {
	CookieName  string
	Lock        sync.RWMutex
	MaxLifeTime int64
	sessions    map[string]*Session
}

func NewSessionMgr(cookieName string, maxLifeTime int64) *SessionMgr {
	return &SessionMgr{CookieName: cookieName, MaxLifeTime: maxLifeTime, sessions: make(map[string]*Session)}
}

func (mgr *SessionMgr) StartSession(rw http.ResponseWriter, re *http.Request) {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()

}

func (mgr *SessionMgr) NewSessionID() string {
	hash := md5.New()
	hash.Write([]byte(string(time.Now().UnixNano())))
	sessionID := hex.EncodeToString(hash.Sum(nil))
	return sessionID
}

func (mgr *SessionMgr) GetSessionIDList() []string {
	var sessionIDList []string
	for _, v := range mgr.sessions {
		sessionIDList = append(sessionIDList, v.ID)
	}
	return sessionIDList
}

func (mgr *SessionMgr) SetSessionVal(sessionID string, key string, value interface{}) {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	if v, ok := mgr.sessions[sessionID]; ok {
		v.Values[key] = value
	}
}

func (mgr *SessionMgr) GetSessionVal(sessionID string, key string) (interface{}, bool) {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	if v, ok := mgr.sessions[sessionID]; ok {
		return v.Values[key], true
	}
	return nil, false
}

func (mgr *SessionMgr) GetLastAccess(sessionID string) time.Time {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	if v, ok := mgr.sessions[sessionID]; ok {
		return v.LastAccess
	}
	return time.Now()
}

func (mgr *SessionMgr) EndSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(mgr.CookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	delete(mgr.sessions, cookie.Value)
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}

func (mgr *SessionMgr) EndSessionByID(sessionID string) {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	delete(mgr.sessions, sessionID)
}

func (mgr *SessionMgr) CheckValid(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie(mgr.CookieName)
	if err != nil {
		log.Println(err)
		return false
	}
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	if session, ok := mgr.sessions[cookie.Value]; ok {
		session.LastAccess = time.Now()
		return true
	}
	return false
}

func (mgr *SessionMgr) GC() {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()

	for sessionID, session := range mgr.sessions {
		if time.Now().Unix()-session.LastAccess.Unix() > mgr.MaxLifeTime {
			delete(mgr.sessions, sessionID)
		}
	}
	time.AfterFunc(time.Duration(mgr.MaxLifeTime)*time.Second, func() { mgr.GC() })
}
