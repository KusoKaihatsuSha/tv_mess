package main

import (
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ConvertToMp3(*Tasker, Thing, int, *Message)
// Using for convert mp4 to mp3
func ConvertToMp3(T *Tasker, task Thing, try int, message *Message) {
	defer new(Timer).Start().Stop()
	v := task.Input.(*JsonPls)
	select {
	case <-T.Branch.Context.Done():
		return
	default:
	}
	if GetCtx[bool](T, v.URLSaved+jpg, message) && GetCtx[bool](T, v.URLSaved+mp4, message) && existFile(v.URLSaved+jpg) != "" && existFile(v.URLSaved+mp4) != "" {
		if existFile(v.URLSaved+mp3) != "" {
			os.Remove(v.URLSaved + mp3)
		}
		var args []string
		args = append(args, "-i")
		args = append(args, v.URLSaved+mp4)
		args = append(args, "-i")
		args = append(args, v.URLSaved+jpg)
		args = append(args, "-map")
		args = append(args, "0")
		args = append(args, "-map")
		args = append(args, "1")
		args = append(args, "-metadata")
		args = append(args, `title=`+v.Song)
		args = append(args, "-metadata")
		args = append(args, `artist=`+v.Artist)
		args = append(args, "-metadata")
		args = append(args, `track=`+v.ID)
		args = append(args, v.URLSaved+mp3)
		err := exec.Command("ffmpeg", args...).Run()
		if err != nil {
			v.toLog(mp3, true)
			if try <= tryingDownload {
				ConvertToMp3(T, task, try+1, message)
			} else {
				message.AddCtx(T, v.URLSaved+mp3, true)
			}
		} else {
			v.toLog(mp3)
			message.AddCtx(T, v.URLSaved+mp3, true)
		}
	}
}

// SplitMp(*Tasker, Thing, int, int64, string, *Message)
// Using for partialing files, if size more than limit
func SplitMp(T *Tasker, task Thing, try int, limit int64, format string, message *Message) {
	defer new(Timer).Start().Stop()
	v := task.Input.(*JsonPls)
	select {
	case <-T.Branch.Context.Done():
		return
	default:
	}
	if GetCtx[bool](T, v.URLSaved+jpg, message) && GetCtx[bool](T, v.URLSaved+format, message) && existFile(v.URLSaved+jpg) != "" && existFile(v.URLSaved+format) != "" {
		var args []string
		args = append(args, "-i")
		args = append(args, v.URLSaved+format)
		args = append(args, "-c")
		args = append(args, "copy")
		args = append(args, "-map")
		args = append(args, "0")
		args = append(args, "-segment_time")
		args = append(args, FfprobeSize(T, v.URLSaved+format, limit))
		args = append(args, "-f")
		args = append(args, "segment")
		args = append(args, "-reset_timestamps")
		args = append(args, "1")
		args = append(args, v.URLSaved+"__%04d"+format)
		err := exec.Command("ffmpeg", args...).Run()
		if err != nil {
			v.toLog(format+" fail split", true)
			if try <= tryingDownload {
				SplitMp(T, task, try+1, limit, format, message)
			} else {
			}
		} else {
			v.toLog(format + " split")
			for _, val := range searchFiles(v.URLSaved+"__", v.UUID, format) {
				message.AddCtx(T, val, true)
			}
			go func() {
				if !debug {
					os.Remove(v.URLSaved + jpg)
					timeout := false
					var notfound = make(chan struct{})
					go func() {
						for {
							select {
							case <-T.Branch.Context.Done():
								break
							default:
								if os.Remove(v.URLSaved+format) == nil || timeout {
									notfound <- struct{}{}
									close(notfound)
									return
								}
							}
						}
					}()
					select {
					case <-time.After(10 * time.Minute):
						timeout = true
					case <-notfound:
					}
					os.Remove(v.URLSaved + mp4)
					os.Remove(v.UUID)
				}
			}()
		}
	}
}

// ffprobeSize(*Tasker, string, int64) string
// function use for calc file size
func FfprobeSize(T *Tasker, url string, limit int64) string {
	defer new(Timer).Start().Stop()
	select {
	case <-T.Branch.Context.Done():
		return ""
	default:
	}
	var chunks int64
	if fileSize(url)/limit != 0 {
		chunks = (fileSize(url) / limit) + 1
	} else {
		chunks = (fileSize(url) / limit)
	}
	var args []string
	args = append(args, "-v")
	args = append(args, "error")
	args = append(args, "-show_entries")
	args = append(args, "format=duration")
	args = append(args, "-of")
	args = append(args, "default=noprint_wrappers=1:nokey=1")
	args = append(args, url)
	size, _ := exec.Command("ffprobe", args...).Output()
	i, _ := strconv.ParseFloat(strings.TrimSpace(string(size)), 10)
	j := int64(math.Round(i)) / chunks
	return sprintf("%02d:%02d:%02d", j/3600, (j % 3600 / 60), ((j % 3600) % 60))
}
