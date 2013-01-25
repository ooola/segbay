FILES      := /Users/ola/go
RSYNC      := rsync -rltgoD --delete --force -e "ssh -o ClearAllForwardings=yes"
SCP        := scp -o ClearAllForwardings=yes
BUILDHOST  := scrape
BUILDCMD   := ssh -o ClearAllForwardings=yes $(BUILDHOST) "(cd ~/go/src/segbay/server && go build)"
BUILDBIN   := /home/ola/go/src/segbay/server/server
TESTHOST   := segbay
TESTBIN    := /usr/local/bidserver/bidserver
TESTPY     := /usr/local/bidserver/bid.py
FIXUPBIN   := ssh -o ClearAllForwardings=yes $(TESTHOST) "chmod 755 $(TESTBIN) $(TESTPY) && chown root:root $(TESTBIN) $(TESTPY)"
BIN        := server
PYBIN      := bid.py

LOCALHTML  := /Users/ola/Dropbox/Sources/ebay-buy/segbaysnipe/public_html/
REMOTEHTML := /home/segbaysn/public_html/
FIXUPHTML  := ssh -o ClearAllForwardings=yes $(TESTHOST) "(cd /home/segbaysn/public_html/ && chown -R segbaysn:segbaysn * .ftpquota .htaccess && chown segbaysn:segbaysn ../public_html && find . -type f -exec chmod 644 {} \; && find . -type d -exec chmod 755 {} \;)"

LOCALRC	   := bidserver
REMOTERC   := /etc/init.d/bidserver
FIXUPRC    := ssh -o ClearAllForwardings=yes $(TESTHOST) "chmod 755 $(REMOTERC) && chown root:root $(REMOTERC)"
RESTART    := ssh -o ClearAllForwardings=yes $(TESTHOST) "$(REMOTERC) restart"

default: sync build upload

.PHONY: sync build upload install restart html 
sync:
	rm -f $(BIN)
	rm -f $(FILES)/src/segbay/server/server
	$(RSYNC) $(FILES) $(BUILDHOST):.

build: sync
	$(BUILDCMD)

upload: build
	$(RSYNC) $(BUILDHOST):$(BUILDBIN) $(BIN)
	$(RSYNC) $(BIN) $(TESTHOST):$(TESTBIN)
	$(RSYNC) $(PYBIN) $(TESTHOST):$(TESTPY)
	$(FIXUPBIN)

install: upload
	$(RSYNC) $(LOCALRC) $(TESTHOST):$(REMOTERC)
	$(FIXUPRC)

restart: install
	$(RESTART)

html:
	$(RSYNC) $(LOCALHTML) $(TESTHOST):$(REMOTEHTML)
	$(FIXUPHTML)

# TODO: add chkconfig
# root@seg [/etc/init.d]# chkconfig --add bidserver
# root@seg [/etc/init.d]# chkconfig --level 35 bidserver on
# root@seg [/etc/init.d]# chkconfig --list bidserver
# bidserver      	0:off	1:off	2:off	3:on	4:off	5:on	6:off
# root@seg [/etc/init.d]# chkconfig --list apache
# error reading information on service apache: No such file or directory
# root@seg [/etc/init.d]# chkconfig --list httpd
# httpd          	0:off	1:off	2:off	3:on	4:off	5:on	6:off
