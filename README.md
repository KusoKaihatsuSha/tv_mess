[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![godoc](https://godoc.org/github.com/KusoKaihatsuSha/tv_mess?status.svg)](https://godoc.org/github.com/KusoKaihatsuSha/tv_mess) [![Go Report Card](https://goreportcard.com/badge/github.com/KusoKaihatsuSha/tv_mess)](https://goreportcard.com/report/github.com/KusoKaihatsuSha/tv_mess)

# Telegram video & music easy self saving (*TV & MESS*)

App consist of functionality Telegram Bot and functionality for download mp4, convert to mp3 using ffmpeg, crop JPG and then add to ID3(MP3). Current tested variants for deploy you can saw below.

### Notes: 
> Sometimes you want test some content as binary files (load from ***youtube*** or ***music.youtube*** mosquito noise, for example). And you too pity spend time for find some site, filled adv. Plus you want ***App*** in your control and on your side.
If you need get ***mp4*** or ***mp3*** file for testing by ***Telegram Bot*** use this repo.

## Choose your variant (if exist):

### 1 (Orange Pi 3 LTS):

**1.** You buy cheapest single-board PC and PC with Windows 10/11(for example):

   > Orange Pi 3 LTS

**2.** Use **balenaEtcher** and flash microSD with **Debian** from official site Orange Pi

   > http://www.orangepi.org/downloadresources/

**3.** Pick microSD into slot on board OPI and hold power button, until you saw Debian loader.

**4.** login as **orangepi**/**orangepi**

**5.** Relocate Debian Boot into board memory(8GB internal) with command:

   ```sh
$ sudo nand-sata-install
   ```

**6.** Choose 2 and format memory in ext4. Wait complete, power off, pickup microSD from slot, power on.

**7.** Install Docker:

   ```sh
$ sudo apt-get update
$ sudo apt-get upgrade
   ```

   From official site you may find tiny command for install Docker if board have internet access:

   ```sh
$ curl -fsSL test.docker.com -o get-docker.sh && sh get-docker.sh
   ```

   Add current user for Docker using:

   ```sh
$ sudo usermod -aG docker $USER
   ```

   Reboot singe-board OPI:

   ```sh
$ sudo reboot
   ```

   Test Docker after reboot:

   ```sh
$ docker run hello-world
   ```

**8.** Install Docker on Windows 10/11 PC (if not exist). 

**9.** For using ssh-agent on Windows 10/11 PC (if not)

   Run as Administrator:

   ```sh
$ Get-Service -Name ssh-agent | Set-Service -StartupType Manual
   ```

   ssh-agent need startup on windows start or run manually:

   ```sh
$ ssh-agent start
   ```

**10.** Generate key for accessing enter without credentials:

   ```sh
$ ssh-keygen -t ecdsa -b 521 -f remoteKeytoremoteNameIp
   ```

   Rename key as you wish, in example key named as "remoteKeytoremoteNameIp" and edit data public key: *user@computer* => *remoteOpiUser@remoteOpiIp*
    
   Then pull key into OPI:

   ```sh
$ scp C:/Users/USERFOLDER/.ssh/remoteKeytoremoteNameIp.pub remoteOpiUser@remoteOpiIp:/home/remoteOpiUser/.ssh/authorized_keys
   ```

   If problem with folder authorized_keys on OPI:
    

   ```sh
$ ssh remoteOpiUser@remoteOpiIp "touch /home/remoteOpiUser/.ssh/authorized_keys && chmod 600 /home/remoteOpiUser/.ssh/authorized_keys"
   ```

**11.** Create new docker context for using docker from PC windows 10/11:

   ```sh
$ docker context ls
$ docker context create someRemoteName --docker "host=ssh://remoteOpiUser@remoteOpiIp"
$ docker context use someRemoteName
$ docker context ls
   ```

**12.** Create folder for deploy by docker-compose

**13.** Create file with environment vars

   > .env

   ```dockerfile
GAPI=keyGoogleApiv3
TAPI=telegramBotApi
PORT=8910
COUNTTASK=512
DEBUG=0
GIT=https://github.com/KusoKaihatsuSha/tv_mess.git
WEBHOOK=0
HOST=null
   ```

**14.** Create file docker-compose.yml

   > docker-compose.yml

   ```dockerfile
version: "3.8"
services:
      git:
         container_name: tv_mess
         build: ${GIT}
      ports:
         - ${PORT}:${PORT}
      restart: always
      expose:
         - ${PORT}
      environment:
         GAPI: ${GAPI}
         TAPI: ${TAPI}
         DEBUG: ${DEBUG}
         COUNTTASK: ${COUNTTASK}
         HOST: ${HOST}
         WEBHOOK: ${WEBHOOK}
   ```

**15.** Run command in this folder (**NOT NEED GIT CLONE**):

   ```sh
$ docker-compose up -d --build    
   ```

   For stop: 
    

   ```sh
$ docker-compose down
   ```

### 2 (Heroku):

**1.** Install Heroku and login in Heroku CLI

   ```sh
$ heroku login
   ```

**2.** Run commands:

   ```sh
$ git clone https://github.com/KusoKaihatsuSha/tv_mess.git
$ cd tv_mess
$ rmdir /s .git
$ heroku create -a tv-mess
$ heroku config:set GAPI=googleYoutubeApiKey3 -a tv-mess
$ heroku config:set TAPI=telegramApiKey -a tv-mess
$ heroku config:set PORT=8910 -a tv-mess
$ heroku config:set COUNTTASK=512 -a tv-mess
$ heroku config:set DEBUG=0 -a tv-mess
$ heroku config:set WEBHOOK=1 -a tv-mess
$ heroku config:set HOST=tv-mess.herokuapp.com -a tv-mess
$ git init . && git add * && git commit -am "init"
$ git remote add origin https://git.heroku.com/tv-mess.git
$ heroku stack:set container -a tv-mess
$ git push origin master
   ```

### This repo using:


> github.com/boltdb/bolt
>
> github.com/kkdai/youtube/v2
>
> github.com/google/uuid


### Screenshots:


<div style="width:50%">
<img src="/pictures/001.gif" >
</div>
