all: ca.crt inventory.example.com.crt

reset:
	rm -f ca.srl *.crt *.key *.csr

%.key:
	openssl genrsa -out $@

ca.crt: ca.key
	openssl req -new -subj /CN=ca -key ca.key -out ca.csr -days 3650
	openssl x509 -in ca.csr -out ca.crt -req -signkey ca.key -days 3650
	rm ca.csr

%.crt: %.key ca.crt
	openssl req -new -subj /CN=$* -addext "subjectAltName = DNS:$*" -key $*.key -out $*.csr -days 3650
	SAN=DNS:$* openssl x509 -in $*.csr -out $*.crt -req -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 -extfile openssl.cnf -extensions san_env; \
	rm $*.csr

install-ca:
	cp ca.crt /etc/ssl/certs
	ln -s ca.crt /etc/ssl/$$(openssl x509 -in ca.crt -hash -noout).crt

.PRECIOUS: %.key %.crt
