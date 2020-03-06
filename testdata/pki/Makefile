all: ca.crt inventory.example.com.crt

reset:
	rm -f ca.srl *.crt *.key *.csr

%.key:
	openssl genrsa -out $@

%.crt: %.key ca.key
	openssl req -new -subj /CN=$* -key $*.key -out $*.csr -days 3650
	if [ -e ca.crt ]; then \
		openssl x509 -in $*.csr -out $*.crt -req -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650; \
	else  \
		openssl x509 -in $*.csr -out $*.crt -req -signkey ca.key -days 3650; \
	fi
	rm $*.csr

.PRECIOUS: %.key %.crt
