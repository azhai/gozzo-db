package cache

import (
	"fmt"
	"strings"

	"github.com/azhai/gozzo-utils/rdspool"
)

const (
	SESS_ONLINE_KEY = "onlines" // 在线用户
	SESS_PREFIX     = "sess"    // 会话缓存前缀
	SESS_TIMEOUT    = 7200      // 会话缓存时间
	SESS_LIST_SEP   = ";"       // 角色名之间的分隔符
)

var (
	rds      rdspool.Redis
	onlines  = GetRedisHash(SESS_ONLINE_KEY, -1)
	sessions = make(map[string]*RedisBackend)
)

func SetRedisBackend(rdsConn rdspool.Redis) {
	rds = rdsConn
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

// 用户基本信息
type UserInfo struct {
	UID          string  `json:"uid" gorm:"unique_index;size:16;not null;comment:'唯一ID'"` // 唯一ID
	Realname     *string `json:"realname" gorm:"size:30;comment:'昵称/称呼'"`                 // 昵称/称呼
	Avatar       *string `json:"avatar" gorm:"size:100;comment:'头像'"`                     // 头像
	Introduction *string `json:"introduction" gorm:"size:500;comment:'介绍说明'"`             // 介绍说明
}

// 绑定用户信息
func BindUserInfo(sess *RedisBackend, user *UserInfo, roles []string) (oldSid string) {
	// 用于踢掉重复登录
	oldSid, _ = onlines.GetString(user.UID)
	onlines.Set(user.UID, sess.GetName())
	// 缓存用户基本信息
	sess.Set("uid", user.UID)
	sess.Set("roles", SessListJoin(roles))
	if user.Realname != nil {
		sess.Set("name", *user.Realname)
	}
	if user.Avatar != nil {
		sess.Set("avatar", *user.Avatar)
	}
	if user.Introduction != nil {
		sess.Set("introduction", *user.Introduction)
	}
	return
}
