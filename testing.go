package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

var testid = "tO-vtgZxPl0"
var testgapi = os.Getenv("GAPI")
var testname = `Audio Library — Music for content creators__Coral – LiQWYD (No Copyright Music)__tO-vtgZxPl0`

func hash_files_tester(T *Tasker, task Thing, message *Message) {
	v := task.Input.([]any)
	v1 := v[0].(int)
	t := v[1].(*testing.T)
	fmt.Println(">>>", t.Name(), v1)
	idd := testid + sprintf("__%d", v1)
	var hashSum = func(path, h string) error {
		file, err := os.Open(path)
		if err != nil {
			return errors.New("test error open file " + path)
		}
		defer file.Close()
		hash := sha1.New()
		if _, err := io.Copy(hash, file); err != nil {
			return errors.New("test error calc sum sha1 file " + path)
		}
		hash_ := fmt.Sprintf("%x", hash.Sum(nil))
		if existFile(path) != "" && hash_ == h {
			return nil
		} else {
			return errors.New("test error sum sha1 file " + path + " (" + hash_ + " not " + h + ")")
		}
	}
	f1 := hashSum(filepath.Join(idd, testname+mp4), "aafc22c54c796cbaf580a72111ec4f1b53594860")
	f2 := hashSum(filepath.Join(idd, testid+".json"), "e2f4acae9cb9817b149b483ee733b67a5bbed0e6")
	//f3 := hashSum(filepath.Join(idd, testname+jpg), "512c99e2b55eeec16537ead181918d5f83bb774d")//in ram
	f4 := hashSum(filepath.Join(idd, testname+mp3), "a6f9546be69047e1b0eef85cbe9f0757af18c1cf")
	filepath.Join(idd, testname+mp4)
	filepath.Join(idd, testid+".json")
	filepath.Join(idd, testname+jpg)
	filepath.Join(idd, testname+mp3)
	if f1 != nil {
		t.Error(f1)
	} else {
		fmt.Println("success", filepath.Join(idd, testname+mp4))
	}
	if f2 != nil {
		t.Error(f2)
	} else {
		fmt.Println("success", filepath.Join(idd, testid+".json"))
	}
	// if f3 != nil {
	// 	t.Error(f3)
	// } else {
	// 	fmt.Println("success", filepath.Join(idd, testname+jpg))
	// }
	if f4 != nil {
		t.Error(f4)
	} else {
		fmt.Println("success", filepath.Join(idd, testname+mp3))
	}
	fmt.Println(">>>", t.Name(), " - done", v1)
}

func downloading_tasker_tester(T *Tasker, task Thing, message *Message) {
	v := task.Input.([]any)
	v1 := v[0].(int)
	t := v[1].(*testing.T)
	fmt.Println(t.Name())
	fmt.Println(">>>", t.Name(), v1)
	idd := testid + sprintf("__%d", v1)
	debug = true
	playlistQ :=
		resource +
			"playlistItems" + startDelimeter +
			"key=" + testgapi + delimeter +
			"playlistId=" + "%s" + delimeter +
			"part=contentDetails" + delimeter +
			"maxResults=" + maxResults + delimeter +
			"pageToken="
	videoQ :=
		resource +
			"videos" + startDelimeter +
			"key=" + testgapi + delimeter +
			"id=" + "%s" + delimeter +
			"part=snippet,contentDetails"
	obj := new(Action)
	obj.Db = new(DataBase)
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Error(err.Error())
	}
	obj.Db.Open(filepath.Join(path, idd))
	defer obj.Db.Close()
	tmp := new(Query)
	tmp.Tasker = T
	message.UUID = idd
	message.AddCtx(T, "log_path", filepath.Join(idd, testid+".json"))
	usr := new(botUser)
	usr.Name = testid
	usr.New(obj.Db, 0)
	usr.setParameter(paramParam, mp4, true)
	usr.setParameter(paramParam, mp3, true)
	usr.setParameter(paramParam, paramTypeVideo, "140")
	usr.setParameter(paramParam, "uuid", message.UUID)
	message.AddCtx(T, userParam, usr)
	tmp.M = new(sync.RWMutex)
	tmp.playlistQ = playlistQ
	tmp.videoQ = videoQ
	tmp.GetInformationVideo(T, testid, message)
	fmt.Println(">>>", t.Name(), " - done", v1)
}
