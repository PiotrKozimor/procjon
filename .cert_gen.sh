if [ -z "$1" ] 
then
    CERT_DIR=.certs
else
    CERT_DIR=$1
fi
mkdir $CERT_DIR
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
IP.2 = 127.0.0.1" > $CERT_DIR/certificate.conf
openssl genrsa -out $CERT_DIR/ca.key 4096
openssl req -new -x509 -key $CERT_DIR/ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out $CERT_DIR/ca.pem
openssl genrsa -out $CERT_DIR/procjon.key 4096
openssl genrsa -out $CERT_DIR/procjonagent.key 4096
openssl req -new -key $CERT_DIR/procjon.key      -out $CERT_DIR/procjon.csr      -config $CERT_DIR/certificate.conf
openssl req -new -key $CERT_DIR/procjonagent.key -out $CERT_DIR/procjonagent.csr -config $CERT_DIR/certificate.conf
openssl x509 -req -in $CERT_DIR/procjon.csr      -CA $CERT_DIR/ca.pem -CAkey $CERT_DIR/ca.key -CAcreateserial -out $CERT_DIR/procjon.pem      -days 365 -sha256 -extfile $CERT_DIR/certificate.conf -extensions req_ext
openssl x509 -req -in $CERT_DIR/procjonagent.csr -CA $CERT_DIR/ca.pem -CAkey $CERT_DIR/ca.key -CAcreateserial -out $CERT_DIR/procjonagent.pem -days 365 -sha256 -extfile $CERT_DIR/certificate.conf -extensions req_ext
chmod -R 700 $CERT_DIR
