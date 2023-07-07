build:
	go build ./pimouse.go
install:
	go build ./pimouse.go
	cp ./pimouse /usr/local/bin/pimouse
	mkdir -p /etc/pimouse/
	cp ./default.yaml /etc/pimouse/default.yaml
	cp ./pimouse.service /etc/systemd/system/pimouse.service
	systemctl daemon-reload
rebuild:
	systemctl stop pimouse.service
	./tools/remove-gadget.sh
	go build ./pimouse.go
	systemctl start pimouse.service
