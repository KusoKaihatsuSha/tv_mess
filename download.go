package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kkdai/youtube/v2"
)

// Counter type
type WriteCounter struct {
	Current  int64
	Total    int64
	PartSize int64
	Percent  string
	Done     chan bool
	WriteCounterExt
	Tasker  *Tasker
	File    *scrFile
	Message Message
}

// Extention for Counter
type WriteCounterExt struct {
	Title string
	Type  string
}

// type helping wok with file download
type scrFile struct {
	Name     string
	Type     string
	SizeFrom int64
	SizeTo   int64
	From     string
	To       string
	File     *os.File
}

// JsonPls - element of Resulting struct in Query
type JsonPlsMinimal struct {
	Num    int
	ID     string
	Artist string
	Song   string
}

type JsonPls struct {
	JsonPlsMinimal
	Title       string
	URL         string
	URLDl       string
	URLSaved    string
	PicturePath string
	M           sync.RWMutex
	UUID        string
	Status      string
}

// Query does work as central content which handle all program
type Query struct {
	KeyApi       string
	Playlists    []string
	Host         string
	Type         string
	Parameters   map[string]string
	Data         any
	Result       []*JsonPls
	Tasker       *Tasker
	Error        error
	Np           string
	M            *sync.RWMutex
	videoQ       string
	playlistQ    string
	JsonFilename string
}

// Query() []byte
// function wrap around web query
func (o *Query) Query() []byte {
	if o.Data == nil {
		o.Data = &bytes.Reader{}
	}
	r, err := http.NewRequest(o.Type, o.Host, o.Data.(io.Reader))
	if err != nil {
		toLog(sprintf("!!! %s\n", err.Error()))
		o.Error = err
	}
	for k, v := range o.Parameters {
		r.Header.Add(k, v)
	}
	tr := &http.Transport{
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   20 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		ForceAttemptHTTP2:     true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(r)
	if err != nil {
		toLog(sprintf("!!! %s\n", err.Error()))
		o.Error = err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			toLog(sprintf("!!! %s\n", err.Error()))
			o.Error = err
		}
		return body
	}
	return nil
}

// SortWrapperTask(*Tasker, Thing, *Message)
// sort resulting information of playlist elements to file
func (o *Query) SortWrapperTask(T *Tasker, task Thing, message *Message) {
	toLog("SORTING")
	sort.Slice(o.Result, func(i, j int) bool {
		return o.Result[i].Artist+o.Result[i].Song < o.Result[j].Artist+o.Result[j].Song
	})
	for k, v := range o.Result {
		v.Num = k + 1
	}
	message.AddCtx(T, sortinfoParamComplete, true)
}

// PrintWrapperTask(*Tasker, Thing, *Message)
// print resulting information of playlist elements to file
func (o *Query) PrintWrapperTask(T *Tasker, task Thing, message *Message) {
	toLog("PRINTING")
	for _, v := range o.Result {
		printf("%d) %s - %s  [%s]\n", v.Num, v.Artist, v.Song, v.URL)
	}
}

// SaveWrapperTask(*Tasker, Thing, *Message)
// save resulting information of playlist elements to file
func (o *Query) SaveWrapperTask(T *Tasker, task Thing, message *Message) {
	GetCtx[bool](T, sortinfoParamComplete, message)
	toLog("SAVE LIST TO JSON")
	var network bytes.Buffer
	enc := json.NewEncoder(&network)
	enc.SetIndent("", "    ")
	var tmp []*JsonPlsMinimal

	for _, tmpv := range o.Result {
		tmp = append(tmp, &tmpv.JsonPlsMinimal)
	}
	err := enc.Encode(&tmp)
	if err != nil {
		toLog(sprintf("!!! %s\n", err.Error()))
	}
	err = ioutil.WriteFile(GetCtx[string](T, "log_path", message), network.Bytes(), 0755)
	if err != nil {
		toLog(sprintf("!!! %s\n", err.Error()))
	}
	defer message.AddCtx(T, saveinfoParamComplete, true)
}

// Exist(*JsonPls) bool
// check doubles in playlist
func (o *Query) Exist(val *JsonPls) bool {
	for _, v := range o.Result {
		if v.Artist == val.Artist && v.Song == val.Song && v.URL == val.URL {
			v.toLog(sprintf("%s", errors.New("double")), true)
			return true
		}
	}
	return false
}

// GetInformationPlaylist(*Tasker, Thing, *Message)
// get information through YT api v3. Playlist element
func (o *Query) GetInformationPlaylist(T *Tasker, message *Message) {
	for _, vpls := range o.Playlists {
		toLog("GET PLAYLIST", "(", vpls, ")")
		o.GetInformationFromPlaylist(T, vpls, "", message)
	}
	timeout := false
	var notfound = make(chan struct{})
	go func() {
		for {
			if GetCtx[int](T, readinfoParamComplete, message) == len(o.Result) || timeout {
				notfound <- struct{}{}
				close(notfound)
				return
			}
		}
	}()
	select {
	case <-time.After(10 * time.Minute):
		timeout = true
	case <-notfound:
	}
	T.Add(nil, o.SortWrapperTask, message)
	T.Add(nil, o.SaveWrapperTask, message)
}

// GetInformationVideo(*Tasker, Thing, *Message)
// get information through YT api v3. One element
func (o *Query) GetInformationVideo(T *Tasker, id string, message *Message) {
	var vq Query
	vq.Host = fmt.Sprintf(o.videoQ, id)
	vq.Parameters = make(map[string]string)
	vq.Type = http.MethodGet
	list := vq.Query()
	var vJson ItemInformation
	json.Unmarshal(list, &vJson)
	go T.Add(&vJson, o.GetVideoWrapperTask, message)
	err := os.MkdirAll(message.UUID, 0775)
	if err != nil {
		toLog(err.Error())
	}
	timeout := false
	var notfound = make(chan struct{})
	go func() {
		for {
			if len(o.Result) != 0 || timeout {
				notfound <- struct{}{}
				close(notfound)
				return
			}
		}
	}()
	select {
	case <-time.After(10 * time.Minute):
		timeout = true
	case <-notfound:
	}
	T.Add(nil, o.SortWrapperTask, message)
	T.Add(nil, o.SaveWrapperTask, message)
}

// GetInformationFromPlaylist(*Tasker, Thing, *Message)
// get information through YT api v3
func (o *Query) GetInformationFromPlaylist(T *Tasker, vpls, next string, message *Message) {
	toLog("GET PART", "(", next, ")")
	o.Host = fmt.Sprintf(o.playlistQ, vpls) + next
	o.Parameters = make(map[string]string)
	o.Type = http.MethodGet
	o.Data = &bytes.Reader{}
	var ret PlaylistItem
	json.Unmarshal(o.Query(), &ret)
	message.AddCtx(T, readinfoParamComplete, ret.PageInfo.TotalResults)
	T.Add(&ret, o.GetWrapperTask, message)
	if ret.NextPageToken != "" {
		go o.GetInformationFromPlaylist(T, vpls, ret.NextPageToken, message)
	}
}

// GetWrapperTask(*Tasker, Thing, *Message)
// get information through YT api v3
func (o *Query) GetWrapperTask(T *Tasker, task Thing, message *Message) {
	ret := task.Input.(*PlaylistItem)
	for _, vv := range ret.Items {
		var vq Query
		vq.Host = fmt.Sprintf(o.videoQ, vv.ContentDetails.VideoID)
		vq.Parameters = make(map[string]string)
		vq.Type = http.MethodGet
		list := vq.Query()
		var vJson ItemInformation
		json.Unmarshal(list, &vJson)
		T.Add(&vJson, o.GetVideoWrapperTask, message)
	}
}

// GetVideoWrapperTask(*Tasker, Thing, *Message)
// get information through YT api v3
func (o *Query) GetVideoWrapperTask(T *Tasker, task Thing, message *Message) {
	vJson := task.Input.(*ItemInformation)
	for _, vv := range vJson.Items {
		next := new(JsonPls)
		next.Num = 0
		next.M = sync.RWMutex{}
		next.ID = vv.ID
		next.Artist = replaceSpecialSymbols(strings.Replace(vv.Snippet.ChannelTitle, " - Topic", "", -1))
		next.Song = replaceSpecialSymbols(vv.Snippet.Title)
		next.Title = next.Artist + "[" + next.Song + "]"
		next.UUID = message.UUID
		err := os.MkdirAll(next.UUID, 0777)
		if err != nil {
			next.toLog(err.Error(), true)
		}
		next.URLSaved = filepath.Join(next.UUID, next.Artist+"__"+next.Song+"__"+next.ID)
		if vv.Snippet.Thumbnails.Maxres.URL != "" {
			next.PicturePath = vv.Snippet.Thumbnails.Maxres.URL
		} else {
			next.PicturePath = vv.Snippet.Thumbnails.High.URL
		}
		if !o.Exist(next) {
			next.M.RLock()
			o.Result = append(o.Result, next)
			next.M.RUnlock()
			T.Add(next, o.DownloadWrapperTask, message)
		}
	}
}

// DownloadWrapperTask(*Tasker, Thing, *Message)
// download all elements (mp4, jpg, convert to mp3)
func (o *Query) DownloadWrapperTask(T *Tasker, task Thing, message *Message) {
	infoLabel := "⚠"
	select {
	case <-T.Branch.Context.Done():
		return
	default:
	}
	v := task.Input.(*JsonPls)
	if existFile(v.URLSaved+mp4) == "" && existFile(v.URLSaved+mp3) == "" {
		client := youtube.Client{}
		video, err := client.GetVideo(v.ID)
		if err != nil {
			v.toLog(err.Error(), true)
		}
		formats := video.Formats.WithAudioChannels()
		usr := GetCtx[*botUser](T, userParam, message)
		atype, _ := strconv.Atoi(usr.getParameter(paramParam, paramTypeVideo))
		audio := formats.Itag(atype)
		if audio != nil {
			v.URLDl, err = client.GetStreamURL(video, &audio[0])
			if err != nil {
				v.toLog(err.Error(), true)
			}
			if !sBool(usr.getParameter(paramParam, mp3)) && !sBool(usr.getParameter(paramParam, mp4)) {
				message.ReplyMarkup = Buttons{}
				message.DelAfter = false
				message.Text = `<a href="` + v.URLDl + `">` + v.Artist + " [" + v.Song + `]</a>`
				message.extensionMessaging(T, sendParam, false, message.sendMessage)
				go func() {
					if !debug {
						os.RemoveAll(v.UUID)
					}
				}()
				return
			}
			T.Add(v, o.DownloadJpgWrapperTask, message)
			T.Add(v, o.DownloadMp4WrapperTask, message)
			T.Add(v, o.DownloadMp3WrapperTask, message)
		} else {
			v.toLog("Error LINK", true)
			message.ReplyMarkup = Buttons{}
			message.DelAfter = false
			message.Text = message.tr(infoLabel+"[choose other quality] ") + v.Artist + "_" + v.Song
			message.extensionMessaging(T, sendParam, false, message.sendMessage)
		}

	}
}

// AppendResultWrapperTask(*Tasker, Thing, *Message)
// add to list of elements
func (o *Query) AppendResultWrapperTask(T *Tasker, task Thing, message *Message) {
	v := task.Input.(*JsonPls)
	v.M.RLock()
	o.Result = append(o.Result, v)
	v.M.RUnlock()
}

// DownloadMp4WrapperTask(*Tasker, Thing, *Message)
// wrapper for mp4 downloading
func (o *Query) DownloadMp4WrapperTask(T *Tasker, task Thing, message *Message) {
	select {
	case <-T.Branch.Context.Done():
		return
	default:
	}
	v := task.Input.(*JsonPls)
	DownloadFile(T, v.URLDl, v.URLSaved+mp4, 0, task, message)
}

// DownloadMp3WrapperTask(*Tasker, Thing, *Message)
// wrapper for mp3 converting
func (*Query) DownloadMp3WrapperTask(T *Tasker, task Thing, message *Message) {
	select {
	case <-T.Branch.Context.Done():
		return
	default:
	}
	v := task.Input.(*JsonPls)
	usr := GetCtx[*botUser](T, userParam, message)
	if sBool(usr.getParameter(paramParam, mp3)) {
		ConvertToMp3(T, task, 0, message)
		if !debug {
			sendFiles(T, task, mp3, message)
		}
	} else {
		message.AddCtx(T, paramGood+v.URLSaved+mp3, true)
	}
}

// DownloadJpgWrapperTask(*Tasker, Thing, *Message)
// wrapper for picture downloading
func (*Query) DownloadJpgWrapperTask(T *Tasker, task Thing, message *Message) {
	v := task.Input.(*JsonPls)
	DownloadFile(T, v.PicturePath, v.URLSaved+jpg, 0, task, message)
}

// Write([]byte) (int, error)
// method for io.Writer interface
func (o *WriteCounter) Write(p []byte) (int, error) {
	var n int
	var err error
	select {
	case <-o.Tasker.Branch.Context.Done():
		n = 0
		err = errors.New("Cancel download")
	default:
		if o.Total == 0 {
			o.File.getWebSize()
		}
		n = len(p)
		o.Total += int64(n)
		go o.Progress()
	}
	return n, err
}

// Progress()
// get percent download by teereader
func (o *WriteCounter) Progress() {
	percent := float64(o.Total*100) / float64(o.File.SizeFrom)
	o.Percent = fmt.Sprintf("%4.2f", percent)
	if percent >= 100 {
		o.Done <- true
	}
}

// getWebSize(string) int64
// file size on server side
func getWebSize(url string) int64 {
	client := &http.Client{}
	for i := 1; i < 50; i++ {
		response, err := client.Get(url)
		if err != nil {
			return 0
		}
		time.Sleep(3 * time.Millisecond)
		if response.ContentLength > 0 {
			return response.ContentLength
		}
		defer response.Body.Close()
	}
	return 0
}

// getWebSize()
// file size on server side
func (o *scrFile) getWebSize() {
	client := &http.Client{}
	for i := 1; i < 50; i++ { // may problem with first try
		response, err := client.Get(o.From)
		if err != nil || response.ContentLength == 0 {
			o.SizeFrom = 0
			time.Sleep(3 * time.Millisecond)
			continue
		}
		if response.ContentLength > 0 {
			o.SizeFrom = response.ContentLength
			return
		}
		defer response.Body.Close()
	}
	o.SizeFrom = 0
}

// fileSize()
// trying get file size wrap scrFile
func (o *scrFile) fileSize() {
	for i := 1; i < 50; i++ { // may problem with first try
		fi, err := os.Stat(o.To)
		if err != nil || fi.Size() == 0 {
			o.SizeTo = 0
			time.Sleep(3 * time.Millisecond)
			continue
		}
		o.SizeTo = fi.Size()
		return
	}
	o.SizeTo = 0
}

// initScrFile(string, string) *scrFile
// init struct helper for downloading
func initScrFile(to, from string) *scrFile {
	var err error
	file := scrFile{To: to, From: from}
	file.File, err = os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		toLog(err.Error())
	}
	return &file
}

// fileSize(string) int64
// url filesize from os information
func fileSize(url string) int64 {
	fi, err := os.Stat(url)
	var fileSize int64
	if err == nil {
		fileSize = fi.Size()
	}
	return fileSize
}

// load(string, *io.PipeWriter, int64, *WriteCounter)
// fast downloading method with partial split
func load(url string, w *io.PipeWriter, begin int64, counter *WriteCounter) {
	defer new(Timer).Start().Stop()
	defer w.Close()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		toLog(err.Error())
	}
	end := begin + counter.PartSize //partsize
	header := fmt.Sprintf("bytes=%v-%v", begin, end)
	req.Header.Set("Range", header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusPartialContent { // OR you can use code 206
		current, err := io.Copy(w, io.TeeReader(resp.Body, counter))
		if current == int64(0) || err != nil {
			return
		}
		select {
		case <-counter.Tasker.Branch.Context.Done():
			return
		default:
			load(url, w, begin+current, counter)
		}
	} else {
		if resp.StatusCode != 416 {
			toLog(errors.New(sprintf("%d", resp.StatusCode) + " - error downloading " + counter.File.To).Error())
		}
		return
	}
}

// DownloadFile(*Tasker, string, string, int, Thing, *Message)
// initialization downloading files
func DownloadFile(T *Tasker, from, to string, try int, task Thing, message *Message) {
	select {
	case <-T.Branch.Context.Done():
		return
	default:
	}
	var err error
	pr, pw := io.Pipe()
	// Buttons init
	versions := make(map[string]string)
	versions[message.tr("cancel")] = "/" + commandCancel + message.UUID
	// Counter init
	counter := &WriteCounter{Tasker: T, PartSize: partSize, File: initScrFile(to, from), Done: make(chan bool), Message: Message{MinimalMessage: MinimalMessage{MessageId: message.MessageId, ChatID: message.ChatID, Text: "0 %", ReturnMessageId: message.MessageId, ParseMode: "HTML"}, ReplyMarkup: Button{InlineKeyboard: buttomsMap(versions)}}}
	mask := regexp.MustCompile(`\..+$`)
	counter.Type = mask.FindString(to)
	if counter.Type == "" {
		counter.Type = ".unknown"
	}
	mask = regexp.MustCompile(`[\\/]`)
	tempSplit := strings.Split(mask.ReplaceAllString(to, delimeterString), delimeterString)
	switch {
	case len(tempSplit) >= 1:
		counter.Title = tempSplit[len(tempSplit)-1]
	default:
		counter.Title = message.tr("Error name")
	}
	// Downloading and counting process
	go load(from, pw, int64(0), counter)
	go counter.informMessage(task, message)
	var written int64
	switch {
	case strings.HasSuffix(to, mp4):
		written, err = io.Copy(counter.File.File, pr)
		defer counter.File.File.Close()
	case strings.HasSuffix(to, jpg):
		// Download front JPG
		my_image, err := jpeg.Decode(pr)
		if my_image != nil && err == nil {
			bounds := my_image.Bounds()
			calcXY := func(a, bmin, bmax int, aIsY bool) int {
				if bmin > bmax {
					for b := bmin; b > bmax; b-- {
						aa, bb, cc, dd := a, b, a, b-1
						ee := &dd
						if aIsY {
							aa, bb, cc, dd = b, a, b-1, a
							ee = &cc
						}
						r1, g1, b1, _ := my_image.At(aa, bb).RGBA()
						r2, g2, b2, _ := my_image.At(cc, dd).RGBA()
						if r1 != r2 && g1 != g2 && b1 != b2 {
							return *ee
						}
					}
				} else {
					for b := bmin; b < bmax; b++ {
						aa, bb, cc, dd := a, b, a, b+1
						ee := &dd
						if aIsY {
							aa, bb, cc, dd = b, a, b+1, a
							ee = &cc
						}
						r1, g1, b1, _ := my_image.At(aa, bb).RGBA()
						r2, g2, b2, _ := my_image.At(cc, dd).RGBA()
						if r1 != r2 && g1 != g2 && b1 != b2 {
							return *ee
						}
					}
				}
				return 0
			}
			if err != nil {
				toLog(err.Error())
			}
			my_sub_image := my_image.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(calcXY(bounds.Max.Y/2, bounds.Min.X+1, bounds.Max.X, true), calcXY(bounds.Max.X/2, bounds.Min.Y+1, bounds.Max.Y, false), calcXY(bounds.Max.Y/2, bounds.Max.X-1, bounds.Min.X, true), calcXY(bounds.Max.X/2, bounds.Max.Y-1, bounds.Min.Y, false)))
			err = jpeg.Encode(counter.File.File, my_sub_image, nil)
			if err != nil {
				toLog(err.Error())
			}
			defer counter.File.File.Close()
		} else {
			// if picture not found in normal quality for example will create dummy picture
			myImg := image.NewRGBA(image.Rect(0, 0, 350, 350))
			d := myImg.Bounds().Dx()
			for by := 0; by < 350; by++ {
				for bx := 0; bx < 350; bx++ {
					xx, yy, rr := float64(bx-d/2)+0.5, float64(by-d/2)+0.5, float64(d/2)
					if xx*xx+yy*yy < rr*rr {
						myImg.Set(bx, by, color.RGBA{255, 99, 71, 0xff})
					} else {
						myImg.Set(bx, by, color.RGBA{34, 139, 87, 0xff})
					}
				}
			}
			err = jpeg.Encode(counter.File.File, myImg, nil)
			if err != nil {
				toLog(err.Error())
			}
			defer counter.File.File.Close()
		}
	default:
	}
	if err != nil {
		toLog(err.Error())
		return
	}
	download := true
	if written == 0 && err != nil {
		if try <= tryingDownload {
			select {
			case <-T.Branch.Context.Done():
			default:
				go DownloadFile(T, from, to, try+1, task, message)
			}
		} else {
			download = false
		}
	}
	message.AddCtx(T, to, download)
	if !debug {
		sendFiles(T, task, counter.Type, message)
	}
}

// informMessage(Thing, *Message)
// function send and edit message, when going downloading process
func (o *WriteCounter) informMessage(task Thing, message *Message) {
	taskInput := task.Input.(*JsonPls)
	user := GetCtx[*botUser](o.Tasker, userParam, message)
	if sBool(user.getParameter(paramParam, o.Type)) && !debug {
		text := taskInput.Title + " <b>" + message.tr("Preparing") + " " + o.Type + "</b>"
		r := new(SendMessageReturn)
		telegramQuery(sendParam, o.Message, r, false, "")
		o.Message.MessageId = r.Result.MessageID
		for {
			select {
			case done, ok := <-o.Done:
				select {
				case <-o.Tasker.Branch.Context.Done():
					telegramQuery(delParam, o.Message, new(DeleteMessageReturn), false, "")
					return
				default:
					if done && ok {
						if sBool(user.getParameter(paramParam, mp3)) && o.Type == mp4 {
							o.Message.Text = o.Title + " <b>" + message.tr("Convertation to MP3") + "</b>"
							telegramQuery(editParam, o.Message, new(SendMessageReturn), false, "")
							GetCtx[bool](o.Tasker, taskInput.URLSaved+mp3, message)
						}
						telegramQuery(delParam, o.Message, new(DeleteMessageReturn), false, "")
					}
				}
				return
			default:
				select {
				case <-o.Tasker.Branch.Context.Done():
					telegramQuery(delParam, o.Message, new(DeleteMessageReturn), false, "")
					return
				default:
					o.Message.Text = text + ": <b>" + o.Percent + " %</b>"
					telegramQuery(editParam, o.Message, new(SendMessageReturn), false, "")
					<-time.After(3 * time.Second)
				}
			}
		}
		close(o.Done)
	}
}
