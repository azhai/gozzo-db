package cache

import (
	"fmt"
	"strings"

	"github.com/azhai/gozzo-utils/common"
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
	rds      rdspool.Redis
	onlines  *RedisBackend
	sessions = make(map[string]*RedisBackend)
)

func SetRedisBackend(rdsConn rdspool.Redis) {
	rds = rdsConn
	onlines = GetRedisHash(SESS_ONLINE_KEY, MAX_TIMEOUT)
}

func GetRedisHash(key string, timeout int) *RedisBackend {
	if sess, ok := sessions[key]; ok {
		return sess
	}
	sess := NewRedisBackend(key, timeout)
	_ = sess.SetRedisInst(rds)
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

func GetSession(token string) *RedisBackend {
	key := fmt.Sprintf("%s:%s", SESS_PREFIX, token)
	return GetRedisHash(key, SESS_TIMEOUT)
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
func BindUserRoles(sess *RedisBackend, uid string, roles []string) (string, error) {
	newSid := sess.GetName()
	oldSid, _ := onlines.GetString(uid) // 用于踢掉重复登录
	if oldSid == newSid {               // 同一个token
		oldSid = ""
	}
	_, err := onlines.Set(uid, newSid)
	_, err = sess.Set("uid", uid)
	_, err = sess.Set("roles", SessListJoin(roles))
	return oldSid, err
}

// 绑定用户信息
func BindUserInfo(sess *RedisBackend, info map[string]string) error {
	var args []string
	for key, value := range info {
		args = append(args, key, value)
	}
	_, err := sess.DoWith("HMSET", common.StrToList(args)...)
	return err
}
