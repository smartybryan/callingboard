#!/usr/bin/make -f

host = ${OYDHOST}
keys := ${OYDKEYS}

clean:
	rm -rf workspace
	mkdir workspace

compile: clean
	cd cmd && GOOS=linux GOARCH=amd64 go build -o ../workspace/callingboard.o .

deploy: compile
	ssh -i $(keys) $(host) 'mkdir -p callingboard/html'
	scp -rp -i $(keys) workspace/callingboard.o $(host):/home/ec2-user/callingboard/callingboard.o
	cd web && scp -rp -i $(keys) html $(host):/home/ec2-user/callingboard/
