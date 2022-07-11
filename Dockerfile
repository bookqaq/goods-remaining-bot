FROM golang:1.18.2-buster

RUN apt update

RUN apt install -y systemd-sysv build-essential rpm

WORKDIR /usr/src/goods-remaining-bot

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.io,direct

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN rpm -ivh data/go-cqhttp/go-cqhttp_linux_amd64.rpm

RUN mv systemd/go-cqhttp.service /etc/systemd/system/go-cqhttp.service
RUN mv systemd/goods-remaining-bot.service /etc/systemd/system/goods-remaining-bot.service
RUN systemctl enable go-cqhttp.service
RUN systemctl enable goods-remaining-bot.service 

RUN go build -v . 

CMD ["/sbin/init"]