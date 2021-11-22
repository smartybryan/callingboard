#!/usr/bin/make -f

host := ec2-user@ec2-3-141-170-241.us-east-2.compute.amazonaws.com
keys := /Users/bryan/oldyellowdoor/oldyellowdoor.pem

clean:
	rm -rf workspace
	mkdir workspace

compile: clean
	cd cmd && go build -o ../workspace/callorg.o .

deploy: compile
	scp -rp -i $(keys) workspace/callorg.o $(host):/home/ec2-user/callorg
	cd web && scp -rp -i $(keys) html $(host):/home/ec2-user/callorg/
