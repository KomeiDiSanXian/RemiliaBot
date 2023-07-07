// Package nativewife 本地老婆
package nativewife

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/process"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
)

func init() {
	engine := control.Register("本地老婆", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "本地抽老婆",
		Help:              "- 抽wife[@xxx]\n- 添加wife[名字][图片]\n- 删除wife[名字]\n- 删除所有wife\n- [让 | 不让]所有人均可添加wife\n- 查看所有wife",
		PrivateDataFolder: "nwife",
	})
	base := engine.DataFolder()
	baseuri := "file:///" + file.BOTPATH + "/" + base
	engine.OnPrefix("抽wife", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			grpf := strconv.FormatInt(ctx.Event.GroupID, 36)
			wifes, err := os.ReadDir(base + grpf)
			if err != nil {
				ctx.SendChain(message.Text("一个wife也没有哦~"))
				return
			}
			switch len(wifes) {
			case 0:
				ctx.SendChain(message.Text("一个wife也没有哦~"))
			case 1:
				wn := wifes[0].Name()
				ctx.SendChain(message.Text("大家的wife都是", wn, "\n"), message.Image(baseuri+grpf+"/"+wn), message.Text("\n哦~"))
			default:
				// 获取名字
				name := ctx.NickName()
				now := time.Now()
				s := md5.Sum(helper.StringToBytes(fmt.Sprintf("%s%d%d%d", name, now.Year(), now.Month(), now.Day())))
				r := rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(s[:]))))
				n := r.Intn(len(wifes))
				wn := wifes[n].Name()
				ctx.SendChain(message.Text(name, "的wife是", wn, "\n"), message.Image(baseuri+grpf+"/"+wn), message.Text("\n哦~"))
			}
		})
	// 上传一张图
	engine.OnPrefix("添加wife", zero.OnlyGroup, chkAddWifePermission, zero.MustProvidePicture).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := ""
			for _, elem := range ctx.Event.Message {
				if elem.Type == "text" {
					name = strings.ReplaceAll(elem.Data["text"], " ", "")
					name = name[strings.LastIndex(name, "添加wife")+10:]
					name = strings.ReplaceAll(name, "/", "")
					name = strings.ReplaceAll(name, "\\", "")
					break
				}
			}
			if name == "" {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("没有找到wife的名字！"))
				return
			}
			url := ctx.State["image_url"].([]string)[0]
			grpfolder := base + strconv.FormatInt(ctx.Event.GroupID, 36)
			if file.IsNotExist(grpfolder) {
				err := os.Mkdir(grpfolder, 0755)
				if err != nil {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("错误：", err.Error()))
					return
				}
			}
			err := file.DownloadTo(url, grpfolder+"/"+name)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("错误：", err))
				return
			}
			process.SleepAbout1sTo2s()
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("成功！"))
		})
	// 查看所有wife
	engine.OnFullMatch("查看所有wife", zero.AdminPermission, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			grpf := strconv.FormatInt(ctx.Event.GroupID, 36)
			wifes, err := os.ReadDir(base + grpf)
			if err != nil {
				ctx.SendChain(message.Text("一个wife也没有哦~"))
				return
			}
			wifesCnt := len(wifes)
			ctx.SendChain(message.Text("获取到", wifesCnt, "个wife，准备发送相关信息..."))
			process.SleepAbout1sTo2s()
			ctx.Send("将完整发送wife信息, 请耐心等待...")
			wifesSlice := divideWife(wifes, 100)
			// 发送图片消息
			var waitgroup sync.WaitGroup
			var mu sync.Mutex
			sendAllWifes(wifesSlice, baseuri, grpf, ctx, &waitgroup, &mu)
			waitgroup.Wait()
			process.SleepAbout1sTo2s()
			ctx.Send("wife列表发送完成")
		})
	engine.OnPrefix("删除wife", zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := ""
			for _, elem := range ctx.Event.Message {
				if elem.Type == "text" {
					name = strings.ReplaceAll(elem.Data["text"], " ", "")
					name = name[strings.LastIndex(name, "删除wife")+10:]
					name = strings.ReplaceAll(name, "/", "")
					name = strings.ReplaceAll(name, "\\", "")
					break
				}
			}
			if name != "" {
				grpfolder := base + strconv.FormatInt(ctx.Event.GroupID, 36)
				err := os.Remove(grpfolder + "/" + name)
				if err == nil {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("成功！"))
				} else {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("错误：", err.Error()))
				}
			} else {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("没有找到wife的名字！"))
			}
		})
	// 一次删除
	engine.OnFullMatch("删除所有wife", zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			grpfolder := base + strconv.FormatInt(ctx.Event.GroupID, 36)
			if err := os.RemoveAll(grpfolder); err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("删除时发生错误：", err))
				return
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("成功！"))
		})
	engine.OnSuffix("所有人均可添加wife", zero.AdminPermission, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			text := ""
			for _, elem := range ctx.Event.Message {
				if elem.Type == "text" {
					text = strings.ReplaceAll(elem.Data["text"], " ", "")
					text = text[:strings.LastIndex(text, "所有人均可添加wife")]
					break
				}
			}
			var err error
			switch text {
			case "设置", "授予", "让":
				err = setEveryoneCanAddWife(ctx, true)
			case "取消", "撤销", "不让":
				err = setEveryoneCanAddWife(ctx, false)
			}
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("错误：", err))
				return
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("成功！"))
		})
}

func chkAddWifePermission(ctx *zero.Ctx) bool {
	gid := ctx.Event.GroupID
	if gid > 0 {
		m, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
		if ok {
			data := m.GetData(gid)
			if data&1 == 1 {
				return true
			}
			return zero.AdminPermission(ctx)
		}
	}
	return false
}

func setEveryoneCanAddWife(ctx *zero.Ctx, canadd bool) error {
	m, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
	if ok {
		if canadd {
			return m.SetData(ctx.Event.GroupID, 1)
		}
		return m.SetData(ctx.Event.GroupID, 0)
	}
	return errors.New("no such plugin")
}

func divideWife(wifes []fs.DirEntry, arrlen int) (wifesSlice [][]fs.DirEntry) {
	wifesCnt := len(wifes)
	for i := 0; i < wifesCnt; i += arrlen {
		end := i + arrlen
		if end > wifesCnt {
			end = wifesCnt
		}
		wifesSlice = append(wifesSlice, wifes[i:end])
	}
	return
}

func sendAllWifes(wifesSlice [][]fs.DirEntry, baseuri, grpf string, ctx *zero.Ctx, waitgroup *sync.WaitGroup, mu *sync.Mutex) {
	for _, wfs := range wifesSlice {
		waitgroup.Add(1)
		go func(wifes []fs.DirEntry) {
			msg := make(message.Message, 0, len(wifes))
			var wg sync.WaitGroup
			for _, wf := range wifes {
				wg.Add(1)
				go func(wife fs.DirEntry) {
					wifeName := wife.Name()
					msgNode := message.Message{
						ctxext.FakeSenderForwardNode(ctx, message.Text("wife名: ", wifeName), message.Image(baseuri+grpf+"/"+wifeName)),
					}
					mu.Lock()
					msg = append(msg, msgNode...)
					mu.Unlock()
					wg.Done()
				}(wf)
			}
			wg.Wait()
			ctx.Send(msg)
			waitgroup.Done()
		}(wfs)
	}
}
