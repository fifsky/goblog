install:
	go run main.go install

stop:
	cat blog.pid | xargs kill

start:
	./blog > /dev/null 2>&1 &
	caddy

restart:
	make stop
	make start

build:
	go build -o blog main.go