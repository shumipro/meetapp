MeetApp
==============

## ローカル環境

### Go環境ない人向けのGo環境構築（Mac & brew版）

```
$ brew update
$ brew install go
$ export GOPATH=$HOME
$ export PATH=$GOPATH/bin:$PATH
```

exportはbashなりzshに追記してください。

### install & build

```
$ go get -u github.com/shumipro/meetapp
$ cd $GOPATH/src/github.com/shumipro/meetapp
$ go build
$ ./meetapp
```

### try

[http://localhost:8000](http://localhost:8000)

## heroku環境

自分用のheroku作りたい人向け

```
$ heroku create -b https://github.com/kr/heroku-buildpack-go.git meetapp-xxx
$ git push heroku master
```

- `xxx`: は適当にかぶらない文字列（いっそランダムでもいいけど）


