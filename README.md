# goods-remaining-bot

余量机器人*1

# ! 以下是最原始的版本完成时实现的教程，仅供参考，如果有误完全是我的问题 !

# 功能

帮助(/余量帮助)，存图(/添加余量)，看图(/看余量)，删图(/删余量)



# 搭建方法:

## 直接在服务器上跑(运维(比如账号掉登录)比较方便)


### 1. 改机器人所属者qq(用于显示余量图编号)

前往bot/main.go，修改14行的MasterQQ的值


### 2. 生成go-cqhttp所需的文件

前往data/go-cqhttp，修改config.yml第4行的uin为机器人的qq号

执行go-cqhttp.exe(双击打开也应该可以，提示请勿直接打开的话就一直点是，生成bat文件后双击打开该文件)并登录成功，保证生成了session.token和device.json
```bash
data/go-cqhttp $ ./go-cqhttp.exe
```

### 3. 把现在的所有文件丢到linux系统中(从该步骤开始我们转到linux侧进行操作)

一定要使用自带systemd的linux发行版，因为我也不太会装
传文件可以用WinSCP什么的，略


### 4. 安装go-cqhttp的rpm包,还有其他要的包

rpm包在data/go-cqhttp中

```bash
$ rpm -ivh data/go-cqhttp/go-cqhttp_linux_amd64.rpm
```

还需要cgo要用到的gcc以及一些工具

使用apt的:

```bash
$ apt install -y build-essential
```

使用yum的:

```bash
$ yum group install "Development Tools"
```

使用pacman(Arch)或者apk(Alpine)等等其他的用户应该都比我nb，我正好也不会用，略


### 5. 复制systemd中的自启动文件，设置开机自启动

```bash
$ cp systemd/* /etc/systemd/system
$ systemctl daemon-reload
$ systemctl enable go-cqhttp
$ systemctl enable goods-remaining-bot
```


### 6. 安装go

官方教程: https://golang.google.cn/doc/install

理论上装1.18的任意版本都可以

```bash
$ curl https://golang.google.cn/dl/go1.18.2.linux-amd64.tar.gz -o go1.18.2.linux-amd64.tar.gz
$ rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.2.linux-amd64.tar.gz
```

添加环境变量
```bash
$ vim /etc/profile  // 如果没有vim的话就试试vi，还没有的话就下载一个
```

在结尾添加该内容后，保存退出(按键盘的a进入insert模式,esc退出该模式,退出后输入```:wq```保存文件)
```bash
export PATH=$PATH:/usr/local/go/bin
```

然后使修改的文件即生效
```bash
$ source /etc/profile
```

此时输入go version，应该就能看到版本号了

```bash
$ go version
go version go1.18.2 linux/amd64
```

### 7. 打包go程序

```bash
$ go mod download && go mod verify
$ go build -v .
```


### 8. 首次运行

```bash
$ systemctl start go-cqhttp
$ systemctl start goods-remaining-bot
```

之后理论上重启就会自动运行



## 使用Docker(可能更简单)

### 1. 改机器人所属者qq(用于显示余量图编号)

前往bot/main.go，修改14行的MasterQQ


### 2. 生成go-cqhttp所需的文件

前往data/go-cqhttp，修改config.yml第4行的uin为机器人的qq号

执行go-cqhttp.exe(双击打开也应该可以，提示请勿直接打开的话就一直点是，生成bat文件后双击打开该文件)并登录成功，保证生成了session.token和device.json
```bash
data/go-cqhttp $ ./go-cqhttp.exe
```


### 3. 把现在的所有文件丢到linux系统中(从该步骤开始我们转到linux侧进行操作)

比如用WinSCP什么的，略


### 4. 安装docker

这个我真没有过经历，还是参考官方吧 https://docs.docker.com/desktop/linux/install/


### 5. 在目录里build该镜像

使用`docker build -t`构建镜像，具体的在下面

```bash
$ docker build -t goods-remaining-bot .
<省略很多输出>
 => exporting to image                                                                            0.8s
 => => exporting layers                                                                           0.7s 
 => => writing image sha256:4a1b16b5911be1fde8ec44d83a096cba27697337e42feac5450d7d0cac22afec      0.0s 
 => => naming to docker.io/library/goods-remaining-bot                                            0.0s
```

完成后可以在`docker image ls`中看到REPOSITORY一个goods-remaining-bot

```bash
$ docker image ls
REPOSITORY                                                    TAG             IMAGE ID       CREATED          SIZE
goods-remaining-bot                                           latest          4a1b16b5911b   24 minutes ago   1.1GB
```


### 6. 运行

```bash
$ docker run -itd --name robo --privileged=true goods-remaining-bot:latest
```

此时可以查看容器有没有跑起来

```bash
$ docker container ps
CONTAINER ID   IMAGE                        COMMAND        CREATED         STATUS         PORTS     NAMES
a2af35a66599   goods-remaining-bot:latest   "/sbin/init"   6 seconds ago   Up 2 seconds             robo
```