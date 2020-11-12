echo "[req]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn
[dn]
C = PL
O = Foo, Inc.
CN = localhost
[req_ext]
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
IP.1 = ::1
IP.2 = 127.0.0.1" > .certs/certificate.conf
openssl genrsa -out .certs/ca.key 4096
openssl req -new -x509 -key .certs/ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out .certs/ca.pem
openssl genrsa -out .certs/procjon.key 4096
openssl genrsa -out .certs/procjonagent.key 4096
openssl req -new -key .certs/procjon.key      -out .certs/procjon.csr      -config .certs/certificate.conf
openssl req -new -key .certs/procjonagent.key -out .certs/procjonagent.csr -config .certs/certificate.conf
openssl x509 -req -in .certs/procjon.csr      -CA .certs/ca.pem -CAkey .certs/ca.key -CAcreateserial -out .certs/procjon.pem      -days 365 -sha256 -extfile .certs/certificate.conf -extensions req_ext
openssl x509 -req -in .certs/procjonagent.csr -CA .certs/ca.pem -CAkey .certs/ca.key -CAcreateserial -out .certs/procjonagent.pem -days 365 -sha256 -extfile .certs/certificate.conf -extensions req_ext
