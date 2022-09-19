package keyword

import (
	"encoding/json"
	"math/rand"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type dict = map[string]*[]string

func init() {
	engine := control.Register("keyword", &ctrl.Options[*zero.Ctx]{
		Help: "关键词匹配回复随机内容, 代码参考Thesaurus插件" +
			"- KEYWORDS",
		DisableOnDefault: true,
		PublicDataFolder: "Keyword",
	})
	go func() {
		data, err := engine.GetLazyData("kw.json", false)
		if err != nil {
			panic(err)
		}
		dictmap := make(dict, 256)
		err = json.Unmarshal(data, &dictmap)
		if err != nil {
			panic(err)
		}
		keylist := make([]string, 0, 256)
		for k := range dictmap {
			keylist = append(keylist, k)
		}

		logrus.Infoln("[Keyword]Loaded ", len(keylist), " keyword(s)")

		engine.OnFullMatchGroup(keylist, zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
			key := ctx.MessageString()
			val := *dictmap[key]
			text := val[rand.Intn(len(val))]
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(text))
		})
		engine.OnFullMatch("KEY").SetBlock(true).Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(""))
		})
	}()
}
