package cache

import (
	"fmt"
	"strings"

	"github.com/azhai/gozzo-utils/rdspool"
)

const (
	MAX_TIMEOUT     = 86400 * 30 // 接近无限时间
	SESS_ONLINE_KEY = "onlines"  // 在线用户
	SESS_PREFIX     = "sess"     // 会话缓存前缀
	SESS_TIMEOUT    = 7200       // 会话缓存时间
	SESS_LIST_SEP   = ";"        // 角色名之间的分隔符
)

var (
	onlines  *RedisBackend
	sessions = make(map[string]*RedisBackend)
)

func GetRedisHash(r rdspool.Redis, key string, timeout int) *RedisBackend {
	if sess, ok := sessions[key]; ok {
		return sess
	}
	sess := NewRedisBackend(key, timeout)
	_ = sess.SetRedisInst(r)
	sessions[key] = sess
	return sess
}

func DelRedisHash(key string) bool {
	if sess, ok := sessions[key]; ok {
		sess.ClearAll()
		delete(sessions, key)
		return true
	}
	return false
}

func InitOnlines() {
	if onlines == nil {
		onlines = GetRedisHash(GetRedisPool(), SESS_ONLINE_KEY, MAX_TIMEOUT)
	}
}

func GetSession(token string) *RedisBackend {
	key := fmt.Sprintf("%s:%s", SESS_PREFIX, token)
	return GetRedisHash(GetRedisPool(), key, SESS_TIMEOUT)
}

func DelSession(token string) bool {
	key := fmt.Sprintf("%s:%s", SESS_PREFIX, token)
	return DelRedisHash(key)
}

func SessListJoin(data []string) string {
	return strings.Join(data, SESS_LIST_SEP)
}

func SessListSplit(data string) []string {
	return strings.Split(data, SESS_LIST_SEP)
}

// 绑定用户角色，返回旧的sid
func (sess *RedisBackend) BindRoles(uid string, roles []string) (string, error) {
	InitOnlines()
	newSid := sess.GetName()
	oldSid, _ := onlines.GetString(uid) // 用于踢掉重复登录
	if oldSid == newSid {               // 同一个token
		oldSid = ""
	}
	_, err := onlines.SetVal(uid, newSid)
	_, err = sess.SetVal("uid", uid)
	_, err = sess.SetVal("roles", SessListJoin(roles))
	return oldSid, err
}
