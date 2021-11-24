#!/usr/bin/make -f

host = ${OYDHOST}
keys := ${OYDKEYS}

clean:
	rm -rf workspace
	mkdir workspace

compile: clean
	cd cmd && GOOS=linux GOARCH=amd64 go build -o ../workspace/callorg.o .

deploy: compile
	ssh -i $(keys) $(host) 'mkdir -p callorg/html'
	scp -rp -i $(keys) workspace/callorg.o $(host):/home/ec2-user/callorg/callorg.o
	cd web && scp -rp -i $(keys) html $(host):/home/ec2-user/callorg/
