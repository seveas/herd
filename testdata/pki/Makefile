all: ca.crt inventory.example.com.crt

reset:
	rm -f ca.srl *.crt *.key *.csr

%.key:
	openssl genrsa -out $@

%.crt: %.key ca.key
	openssl req -new -subj /CN=$* -addext "subjectAltName = DNS:$*" -key $*.key -out $*.csr -days 3650
	if [ -e ca.crt ]; then \
		SAN=DNS:$* openssl x509 -in $*.csr -out $*.crt -req -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 -extfile openssl.cnf -extensions san_env; \
	else  \
		openssl x509 -in $*.csr -out $*.crt -req -signkey ca.key -days 3650; \
	fi
	# rm $*.csr

.PRECIOUS: %.key %.crt
