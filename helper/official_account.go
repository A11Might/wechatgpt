// 公众号
package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"github.com/silenceper/wechat/v2/officialaccount/server"
)

type OfficialAccount struct {
	wc              *wechat.Wechat
	officialAccount *officialaccount.OfficialAccount
}

func NewOfficialAccount() *OfficialAccount {
	wc := wechat.NewWechat()
	//这里本地内存保存access_token，也可选择redis，memcache或者自定cache
	memory := cache.NewMemory()
	cfg := &config.Config{
		AppID:          DefaultConfig.AppID,
		AppSecret:      DefaultConfig.AppSecret,
		Token:          DefaultConfig.Token,
		EncodingAESKey: DefaultConfig.EncodingAESKey,
		Cache:          memory,
	}
	oa := wc.GetOfficialAccount(cfg)
	return &OfficialAccount{
		wc:              wc,
		officialAccount: oa,
	}
}

func (oa *OfficialAccount) GetServer(c *gin.Context) *server.Server {
	return oa.officialAccount.GetServer(c.Request, c.Writer)
}

func (oa *OfficialAccount) GetMaterial() *material.Material {
	return oa.officialAccount.GetMaterial()
}
