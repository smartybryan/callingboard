#!/usr/bin/make -f

host = ${OYDHOST}
keys := ${OYDKEYS}

clean:
	rm -rf workspace
	mkdir workspace

compile: clean
	cd cmd && go build -o ../workspace/callorg.o .

deploy: compile
	scp -rp -i $(keys) workspace/callorg.o $(host):/home/ec2-user/callorg
	cd web && scp -rp -i $(keys) html $(host):/home/ec2-user/callorg/
