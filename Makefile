NAME=around


all: build

build:
	go build -o ${NAME} .

clean:
	rm -f ${NAME}

re: clean build


.PHONY: all build clean re
