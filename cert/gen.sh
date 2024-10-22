rm *.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=DK/ST=Denmark/L=Copenhagen/O=ITU/OU=Education/CN=Casper/emailAddress=csfr@itu.dk"

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=DK/ST=Denmark/L=Copenhagen/O=ITU/OU=Education/CN=Server/emailAddress=csfr@itu.dk"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in server-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client_1-key.pem -out client_1-req.pem -subj "/C=DK/ST=Denmark/L=Copenhagen/O=PC Client/OU=Computer/CN=client_1/emailAddress=csfr@itu.dk"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client_1-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client_1-cert.pem -extfile client-ext.cnf

echo "Client's signed certificate"
openssl x509 -in client_1-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client_2-key.pem -out client_2-req.pem -subj "/C=DK/ST=Denmark/L=Copenhagen/O=PC Client/OU=Computer/CN=client_2/emailAddress=csfr@itu.dk"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client_2-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client_2-cert.pem -extfile client-ext.cnf

echo "Client's signed certificate"
openssl x509 -in client_2-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client_3-key.pem -out client_3-req.pem -subj "/C=DK/ST=Denmark/L=Copenhagen/O=PC Client/OU=Computer/CN=client_3/emailAddress=csfr@itu.dk"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client_3-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client_3-cert.pem -extfile client-ext.cnf

echo "Client's signed certificate"
openssl x509 -in client_3-cert.pem -noout -text