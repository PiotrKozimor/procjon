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
IP.2 = 127.0.0.1" > certificate.conf
openssl genrsa -out ca.key 4096
openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out ca.cert
openssl genrsa -out procjon.key 4096
openssl req -new -key procjon.key -out procjon.csr -config certificate.conf
openssl x509 -req -in procjon.csr -CA ca.cert -CAkey ca.key -CAcreateserial -out procjon.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext
for dir in procjonagent cmd/elastic cmd/systemd
do
    cp procjon.key $dir/
    cp procjon.pem $dir/
    cp ca.cert $dir/
done