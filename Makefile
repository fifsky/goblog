build:
	go build -o blog main.go

install:
	go run main.go install

stop:
	cat blog.pid | xargs kill

start:
	./blog > /dev/null 2>&1 &

restart:
	cat blog.pid | xargs kill
	./blog > /dev/null 2>&1 &
