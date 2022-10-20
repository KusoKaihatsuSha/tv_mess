package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	gtranslate            = "http://translate.google.com/translate_a/single?client=gtx&dt=t&dj=1&ie=UTF-8&sl=%s&tl=%s&q=%s"
	resource              = "https://www.googleapis.com/youtube/v3/"
	maxResults            = "50"
	tryingDownload        = 2
	startDelimeter        = "?"
	delimeter             = "&"
	pathReplace           = `●.=^;-~{}<>/\*+:&?|%#$'"@!`
	mp4                   = ".mp4"
	mp3                   = ".mp3"
	jpg                   = ".jpg"
	telegramUrl           = "https://api.telegram.org/bot"
	DisableNotification   = "/mute"
	Start                 = "start"
	Subscribe             = "subscribe"
	Autoload              = "/auto"
	Private               = "/priv"
	limitFileTelegram     = int64(45000000)
	partSize              = int64(2000000)
	paramParam            = "parameters"
	userParam             = "user"
	sendParam             = "/sendMessage"
	updateParam           = "/getUpdates"
	diceParam             = "/sendDice"
	editParam             = "/editMessageText"
	delParam              = "/deleteMessage"
	readinfoParamComplete = "total_list_done"
	sortinfoParamComplete = "sort_list_done"
	saveinfoParamComplete = "save_list_done"
	paramTypeVideo        = "atype"
	paramLink             = "linkonly"
	paramGood             = "+++"
	infoLabel             = "⚠"
	timerLabel            = "⏳"
	htmlMode              = "HTML"

	delimeterString = "--!--"

	commandCancel        = "!!!cancel!!!"
	commandStart         = "start"
	commandStartConfirm  = "!!!start_confirm_good!!!"
	commandType          = "!!!type!!!"
	commandFind          = "!!!find!!!"
	commandDeleteCurrent = "!!!delthis!!!"
	commandSettingsJpg   = "!!!front_picture!!!"
	commandSettingsLog   = "!!!logs!!!"

	defaultWebHook = "get"

	databasename = "database"
)

var (
	printf  = log.Printf
	sprintf = fmt.Sprintf
	p       = fmt.Println

	yt3key     = ""
	api        = ""
	debug      = false
	usewebhook = false
)

// defHandler(http.ResponseWriter, *http.Request, interface{})
// default handler
func defHandler(w http.ResponseWriter, req *http.Request, ext interface{}) {
	if req.Method == "GET" {
		fmt.Fprintf(w, "fail")
	}
}

// debugHandler(http.ResponseWriter, *http.Request, interface{})
// handler for debugging
func debugHandler(w http.ResponseWriter, req *http.Request, ext interface{}) {
	if req.Method == "GET" {
		logUrl := ext.(string)
		if logUrl != "" && debug {
			text, _ := ioutil.ReadFile(logUrl)
			fmt.Fprintf(w, string(text))
		}
	}
}

// extHandler(func(http.ResponseWriter, *http.Request, interface{}), interface{}) http.HandlerFunc
// handlers func wrapper
func extHandler(fn func(http.ResponseWriter, *http.Request, interface{}), p interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, p)
	}
}

// extHandlerFunc(http.Handler) http.Handler
// handlers wrapper
func extHandlerFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

// setWebHook(string, bool) string
// set or remove webhook on telegram side
func setWebHook(host string, flagDel bool) string {
	ext := "set"
	if flagDel {
		ext = "delete"
	}
	type tmpUrl struct {
		Url string `json:"url"`
	}
	tmpUrlBytes, err := json.Marshal(tmpUrl{Url: "https" + "://" + host + "/" + defaultWebHook})
	if err != nil {
		toLog(err)
	}
	body := bytes.NewReader(tmpUrlBytes)
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+os.Getenv("TAPI")+"/"+ext+"Webhook", body)
	if err != nil {
		toLog(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		toLog(err)
	}
	bodyReturn, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	returnStr := string(tmpUrlBytes) + "\n" + string(bodyReturn)
	toLog(returnStr)
	return returnStr
}

// getHandler(http.ResponseWriter, *http.Request, interface{})
// handler for getting data via webhook
func getHandler(w http.ResponseWriter, req *http.Request, ext interface{}) {
	if req.Method == "POST" {
		extt := ext.([]any)
		MainTasker := extt[0].(*Tasker)
		cmds := extt[1].(*Commands)
		obj := extt[3].(*Action)
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			toLog(err)
		}
		r := new(Result)
		json.Unmarshal(data, &r)
		obj.getHandlerManual(r, MainTasker, cmds, nil)
		fmt.Fprintf(w, "/")
	}
}

// getHandlerManual(*Result, *Tasker, *Commands, error)
// handle telegram incoming data
func (obj *Action) getHandlerManual(val *Result, MainTasker *Tasker, cmds *Commands, err error) {
	tid := fmax(&val.Message.Chat.ID, &val.CallbackQuery.Message.Chat.ID)
	command := val.Message.Text + val.CallbackQuery.Data
	user := new(botUser)
	user.New(obj.Db, tid)
	user.Name = val.Message.Chat.Username
	tempMessage := Message{}
	tempMessage.MessageId = fmax(&val.Message.MessageID, &val.CallbackQuery.Message.MessageID)
	tempMessage.ReturnMessageId = tempMessage.MessageId
	tempMessage.MessageIdStr = sprintf("%d", tempMessage.MessageId)
	tempMessage.ChatID = tid
	tempMessage.ChatIDStr = sprintf("%d", tempMessage.ChatID)
	tempMessage.LanguageCode = fmaxStr(&val.Message.From.LanguageCode, &val.CallbackQuery.From.LanguageCode)
	tempMessage.UUID = replaceSpecialSymbols(uuid.New().String())
	re := regexp.MustCompile(`.+(watch\?v=|youtu.be/)`)
	command = re.ReplaceAllString(command, commandFind)
	re = regexp.MustCompile(`.+playlist\?list=`)
	command = re.ReplaceAllString(command, commandFind)
	if strings.HasPrefix(command, commandFind) {
		re = regexp.MustCompile(`[a-zA-Z0-9_-]{11,41}`)
		command = commandFind + re.FindString(command)
	}
	if strings.HasPrefix(command, commandType) {
		re = regexp.MustCompile(`[0-9]{1,4}`)
		type_ := re.FindString(command)
		typetext_ := ""
		ext := commandType
		user.setParameter(paramParam, mp4, true)
		user.setParameter(paramParam, mp3, true)
		switch type_ {
		case "251":
			typetext_ = "high (audio)"
		case "140":
			typetext_ = "medium (audio)"
		case "249":
			typetext_ = "low (audio)"
		case "22":
			typetext_ = "medium (720p)"
			user.setParameter(paramParam, mp3, false)
		case "18":
			typetext_ = "low (360p)"
			user.setParameter(paramParam, mp3, false)
		default:
			typetext_ = "/settingsQuality"
			if !sBool(user.getParameter(paramParam, paramLink)) {
				user.setParameter(paramParam, paramLink, true)
			} else {
				user.setParameter(paramParam, paramLink, false)
			}
			user.setParameter(paramParam, mp4, false)
			user.setParameter(paramParam, mp3, false)
			ext = ""
		}
		if sBool(user.getParameter(paramParam, paramLink)) {
			user.setParameter(paramParam, mp4, false)
			user.setParameter(paramParam, mp3, false)
		}
		if type_ != "" {
			user.setParameter(paramParam, paramTypeVideo, type_)
		}
		command = ext + typetext_
	}
	tempMessage.Command = command
	tempMessage.User = val.Message.Chat.Username
	tempMessage.ParseMode = htmlMode
	tempMessage.DelAfterDelay = 5 * time.Second
	ok := false
	switch command {
	case "/" + commandStartConfirm:
		tempMessage.ReplyMarkup = cmds.B.Return()
		ok = true
	case "/" + commandStart:
		tempMessage.ReplyMarkup = Buttons{}
		ok = true
	default:
		tempMessage.ReplyMarkup = Buttons{}
		ok = sBool(user.getParameter(paramParam, Subscribe))
	}
	cmdss := cmds.Find(command)
	if cmdss.IsCommand {
		if ok {
			go cmdss.F(tempMessage, user)
			delEmpty()
		} else {
			tempMessage.Text = tempMessage.tr("Try 'start' again, please. → ") + " /start"
			tempMessage.ReplyMarkup = new(Buttons).NewLine().Add("/" + commandStart).Return()
			MainTasker.Add(nil, tempMessage.SendMessageWrapperTask, &tempMessage)
		}
	}
}

// UpdateMsg(*Tasker, *Commands, error)
// handle telegram incoming data manually
func (obj *Action) UpdateMsg(MainTasker *Tasker, cmds *Commands, err error) {
	for {
		r := new(UpdateReturn)
		telegramQuery(updateParam, obj.Update, r, false, "")
		if r != nil && err == nil && len(r.Result) > 0 {
			obj.Update.SetLast(r.Result[0].UpdateID)
			obj.Db.FindCreate("last").Put("id", sprintf("%d", obj.Update.GetLast())) //may be antivirus block if error
			if r != nil || len(r.Result) != 0 {
				for _, val := range r.Result {
					obj.getHandlerManual(&val, MainTasker, cmds, err)
				}
			}
		}
	}
}

// -----
func main() {
	// For deploying need set envs(file .env for example):
	// os.Setenv("GAPI", "xxxXXXxxx")
	// os.Setenv("TAPI", "xxxXXXxxx")
	// os.Setenv("PORT", "8910")
	// os.Setenv("COUNTTASK", "512")
	// os.Setenv("DEBUG", "0")
	// os.Setenv("WEBHOOK", "0")
	// os.Setenv("HOST", "xxxXXXxxx")
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread()
	runtime.Gosched()
	filePtr, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		toLog(err)
	}
	defer filePtr.Close()
	log.SetOutput(filePtr)
	taskscount, err := strconv.Atoi(os.Getenv("COUNTTASK"))
	if err != nil {
		toLog(err)
		taskscount = 512
	}
	yt3key = os.Getenv("GAPI")
	api = os.Getenv("TAPI")
	debug = os.Getenv("DEBUG") == "1"
	usewebhook = os.Getenv("WEBHOOK") == "1"
	MainTasker := new(Tasker).Init(runtime.NumCPU(), taskscount)
	playlistQ :=
		resource +
			"playlistItems" + startDelimeter +
			"key=" + yt3key + delimeter +
			"playlistId=" + "%s" + delimeter +
			"part=contentDetails" + delimeter +
			"maxResults=" + maxResults + delimeter +
			"pageToken="
	videoQ :=
		resource +
			"videos" + startDelimeter +
			"key=" + yt3key + delimeter +
			"id=" + "%s" + delimeter +
			"part=snippet,contentDetails"
	obj := new(Action)
	obj.Db = new(DataBase)
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		toLog(err)
	}
	obj.Db.Open(filepath.Join(path, databasename))
	defer obj.Db.Close()
	markCurrent := func(user *botUser, key_ string, list map[string]string) map[string]string {
		newlist := make(map[string]string)
		for k, v := range list {
			if strings.Contains(v, user.getParameter(paramParam, key_)) {
				newlist["👉 "+k] = v
			} else {
				newlist[k] = v
			}
		}
		return newlist
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", extHandler(defHandler, nil))
	mux.HandleFunc("/debug", extHandler(debugHandler, "logs.log"))
	obj.Update = new(Updates).New()
	cmds := new(Commands).DbConnect(obj.Db)
	cmds.B = new(Buttons)
	cmds.AddNewLine()
	cmds.Add(commandStart, "🆗start", false, false, func(message Message, usr *botUser) {
		message.DelAfter = true
		message.DelAfterDelay = 10 * time.Second
		MainTasker.Add(nil, message.SendRandomWrapperTask, &message)
		val := GetCtx[int](MainTasker, "random", &message)
		message.Text = message.tr("This is 'start' confirmation. Choose right number ↓")
		versions := make(map[string]string)
		pics := map[string]string{"1": "1️⃣", "2": "2️⃣", "3": "3️⃣", "4": "4️⃣", "5": "5️⃣", "6": "6️⃣"}
		for i := 1; i <= 6; i++ {
			if i == val {
				versions[pics[sprintf("%d", i)]] = "/" + commandStartConfirm
			} else {
				versions[pics[sprintf("%d", i)]] = "/" + commandStart
			}
		}
		message.ReplyMarkup = Button{InlineKeyboard: buttomsMap(versions)}
		message.DelAfter = true
		message.DelAfterDelay = 10 * time.Second
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	})
	cmds.Add(commandStartConfirm, "🆗start confirm good", false, false, func(message Message, usr *botUser) {
		usr.setParameter(paramParam, Subscribe, true)
		usr.setParameter(paramParam, mp4, true)
		usr.setParameter(paramParam, mp3, true)
		usr.setParameter(paramParam, commandSettingsLog, false)
		usr.setParameter(paramParam, paramTypeVideo, "140")
		usr.setParameter(paramParam, jpg, false)
		usr.setParameter(paramParam, "add_log", false)
		message.Text = message.tr("Hello! You are subscribe. Paste link to playlist/song. Or you can use the buttons below ↓. Do not delete this message, otherwise you may delete buttons below.")
		message.MessageId = -1
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	})
	cmds.Add("settingsQuality", "⚙Quality", true, true, func(message Message, usr *botUser) {
		message.Text = message.tr("🎵 Choose video/audio quality type")
		message.DelBefore = true
		versions := make(map[string]string)
		versions[message.tr("high (audio)")] = commandType + "251"
		versions[message.tr("medium (audio)")] = commandType + "140"
		versions[message.tr("low (audio)")] = commandType + "249"
		versions[message.tr("medium (720p)")] = commandType + "22"
		versions[message.tr("low (360p)")] = commandType + "18"
		if !sBool(usr.getParameter(paramParam, paramLink)) {
			versions[message.tr("🔴 link mode")] = commandType + paramLink
		} else {
			versions[message.tr("🟢 link mode")] = commandType + paramLink
		}
		versions[message.tr("❌ close")] = "/" + commandDeleteCurrent
		message.ReplyMarkup = Button{InlineKeyboard: buttomsMap(markCurrent(usr, paramTypeVideo, versions))}
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	})
	cmds.Add(commandType, commandType, false, false, func(message Message, usr *botUser) {
		MainTasker.Add(message.Cbmid, message.DeleteMessageWrapperTask, &message)
		text := strings.TrimPrefix(message.Command, commandType)
		message.Text = message.tr("You are choosing - " + text)
		message.DelBefore = true
		message.DelAfter = true
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	})
	cmds.Add("settings2", "⚙JPG", true, true, func(message Message, usr *botUser) {
		message.Text = message.tr("🎴 Do you want JPG of front side?")
		message.DelBefore = true
		versions := make(map[string]string)
		versions[message.tr("➕ JPG")] = "/+++" + commandSettingsJpg
		versions[message.tr("➖ JPG")] = "/---" + commandSettingsJpg
		versions[message.tr("❌ close")] = "/" + commandDeleteCurrent
		message.ReplyMarkup = Button{InlineKeyboard: buttomsMap(markCurrent(usr, jpg, versions))}
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	})
	cmds.Add("settings3", "⚙LOGS", true, true, func(message Message, usr *botUser) {
		message.Text = message.tr("Do you want get logs?")
		message.DelBefore = true
		versions := make(map[string]string)
		versions[message.tr("➕ logs")] = "/+++" + commandSettingsLog
		versions[message.tr("➖ logs")] = "/---" + commandSettingsLog
		versions[message.tr("❌ close")] = "/" + commandDeleteCurrent
		message.ReplyMarkup = Button{InlineKeyboard: buttomsMap(markCurrent(usr, commandSettingsLog, versions))}
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	}).AddNewLine()
	cmds.Add(commandCancel, "⚙cancel", true, false, func(message Message, usr *botUser) {
		message.UUID = strings.TrimPrefix(message.Command, "/"+commandCancel)
		tttt := GetCtx[*BranchContext](MainTasker, "context", &message)
		(*tttt).Cancel()
		for os.RemoveAll(message.UUID) != nil {
			<-time.After(1 * time.Second)
			select {
			case <-tttt.Context.Done():
				break
			default:
			}
		}
	}).AddNewLine()
	cmds.Add("+++"+commandSettingsJpg, "⚙add jpg", false, false, func(message Message, usr *botUser) {
		MainTasker.Add(message.Cbmid, message.DeleteMessageWrapperTask, &message)
		usr.setParameter(paramParam, jpg, true)
		message.Text = message.tr("You are choosed add JPG")
		message.DelBefore = true
		message.DelAfter = true
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	})
	cmds.Add("---"+commandSettingsJpg, "⚙no add jpg", false, false, func(message Message, usr *botUser) {
		MainTasker.Add(message.Cbmid, message.DeleteMessageWrapperTask, &message)
		usr.setParameter(paramParam, jpg, false)
		message.Text = message.tr("You are choosed load without JPG")
		message.DelBefore = true
		message.DelAfter = true
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	})
	cmds.Add("+++"+commandSettingsLog, "⚙add logs", false, false, func(message Message, usr *botUser) {
		MainTasker.Add(message.Cbmid, message.DeleteMessageWrapperTask, &message)
		usr.setParameter(paramParam, commandSettingsLog, true)
		message.Text = message.tr("You are choosed add logs")
		message.DelBefore = true
		message.DelAfter = true
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	})
	cmds.Add("---"+commandSettingsLog, "⚙no add logs", false, false, func(message Message, usr *botUser) {
		MainTasker.Add(message.Cbmid, message.DeleteMessageWrapperTask, &message)
		usr.setParameter(paramParam, commandSettingsLog, false)
		message.Text = message.tr("You are choosed load without logs")
		message.DelBefore = true
		message.DelAfter = true
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	})
	cmds.Add("cthulu", "🐙cthulu", true, false, func(message Message, usr *botUser) {
		message.Text = "🐙Ph'nglui mglw'nafh Cthulhu R'lyeh wgah'nagl fhtagn"
		message.DelBefore = true
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	}).AddNewLine()
	cmds.Add("unsubscibe", "⚙unsubscibe", false, true, func(message Message, usr *botUser) {
		message.Text = message.tr("Will you want unsubscribe? Do you sure?")
		message.DelBefore = true
		versions := make(map[string]string)
		versions[message.tr("Yes")] = "/!!!stop_is_ok"
		versions[message.tr("No")] = "/" + commandDeleteCurrent
		message.ReplyMarkup = Button{InlineKeyboard: buttomsMap(versions)}
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	}).AddNewLine()
	cmds.Add(commandDeleteCurrent, "🔵delete", false, false, func(message Message, usr *botUser) {
		MainTasker.Add(message.Cbmid, message.DeleteMessageWrapperTask, &message)
	})
	cmds.Add("!!!stop_is_ok", "🔵full stop", false, false, func(message Message, usr *botUser) {
		MainTasker.Add(message.Cbmid, message.DeleteMessageWrapperTask, &message)
		message.Text = message.tr("Bye! See ya!")
		message.DelBefore = true
		message.ReplyMarkup = new(Buttons).NewLine().Add("/" + commandStart).Return()
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
		usr.setParameter(paramParam, Subscribe, false)
	})
	cmds.Add("help", "🚫help", false, false, func(message Message, usr *botUser) {
		message.Text = "<i>" + message.Command + "</i> - " + message.tr("Fail operation. Try again, please.")
		message.DelBefore = true
		message.DelAfter = true
		MainTasker.Add(nil, message.SendMessageWrapperTask, &message)
	})
	cmds.Add(commandFind, commandFind, false, false, func(message Message, usr *botUser) {
		playlist := strings.TrimPrefix(message.Command, commandFind)
		MainTaskerT := new(Tasker).Init(runtime.NumCPU(), taskscount)
		message.AddCtx(MainTaskerT, "user", usr)
		time_ := time.Now().Format("2006_01_02_15_04_05")
		message.AddCtx(MainTaskerT, "uuid", message.UUID)
		message.AddCtx(MainTaskerT, mp3, sBool(usr.getParameter(paramParam, mp3)))
		message.AddCtx(MainTaskerT, "log_path", message.UUID+`\`+usr.Name+"_"+time_+".json")
		usr.setParameter(paramParam, "uuid", message.UUID)
		go func() {
			GetCtx[bool](MainTaskerT, saveinfoParamComplete, &message)
			param := DocumentMessage{}
			param.Src = message.UUID + `\` + usr.Name + "_" + time_ + ".json"
			param.Check = sBool(usr.getParameter(paramParam, commandSettingsLog))
			param.Title = "LOGS"
			MainTasker.Add(param, message.SendDocumentWrapperTask, &message)
		}()
		tmp := new(Query)
		tmp.M = new(sync.RWMutex)
		tmp.playlistQ = playlistQ
		tmp.videoQ = videoQ
		tmp.Playlists = strings.Split(playlist, ";")
		message.AddCtx(MainTasker, "context", &MainTaskerT.Branch)
		if len(playlist) > 12 {
			tmp.GetInformationPlaylist(MainTaskerT, &message)
		} else {
			tmp.GetInformationVideo(MainTaskerT, playlist, &message)
		}
		MainTaskerT.Wg.Wait()
		MainTaskerT.Branch.Cancel()
	})
	// Will activate webhook or delete, if not using.
	setWebHook(os.Getenv("HOST"), !usewebhook)
	mux.HandleFunc("/"+defaultWebHook, extHandler(getHandler, []any{MainTasker, cmds, err, obj}))
	// If not using webhook will activate manual getting update data
	if !usewebhook {
		go obj.UpdateMsg(MainTasker, cmds, err)
	}
	toLog("server run on port: " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), extHandlerFunc(mux))
}
