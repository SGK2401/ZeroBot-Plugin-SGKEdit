package xp

import (
	"encoding/json"
	"math/rand"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type dict = map[string]*[]string

func init() {
	engine := control.Register("xp", &ctrl.Options[*zero.Ctx]{
		Help: "随机丢给你一个XP, 代码参考(照搬)Thesaurus插件\n" +
			"- XP[xxx]\n- XPLIST",
		DisableOnDefault: true,
		PublicDataFolder: "XP",
	})
	go func() {
		data, err := engine.GetLazyData("XP.json", false)
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

		logrus.Infoln("[XP]Loaded ", len(keylist), " XP(s)")

		engine.OnPrefix("XP", zero.SuffixRule(keylist...)).SetBlock(true).Handle(func(ctx *zero.Ctx) { // 中间有什么无所谓, 可以修. 懒.
			rawkey := ctx.MessageString()
			key := strings.TrimPrefix(rawkey, "XP") // 包括空格将无法匹配 以后再修
			val := *dictmap[key]
			text := val[rand.Intn(len(val))]
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(text))
		})
		engine.OnFullMatch("XPLIST").SetBlock(true).Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(keylist))
		})
	}()
}
