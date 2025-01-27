#!/bin/bash

# openssl V3

if ls *.pem 1> /dev/null 2>&1; then
    echo "File with .pem extension exists in the current directory."
    echo "This will overwrite the existing files. Either delete or execute this script in other directory."
    exit 1
fi
if ls *.crt 1> /dev/null 2>&1; then
    echo "File with .crt extension exists in the current directory."
    echo "This will overwrite the existing files. Either delete or execute this script in other directory."
    exit 1
fi
if ls *.csr 1> /dev/null 2>&1; then
    echo "File with .csr extension exists in the current directory."
    echo "This will overwrite the existing files. Either delete or execute this script in other directory."
    exit 1
fi

# 1. Generate a RSA private key for the Local CA
openssl genpkey -algorithm RSA -out ca-key.pem -pkeyopt rsa_keygen_bits:2048
if [ $? -ne 0 ]; then
    echo 'Failed to generate Local CA RSA private key'
    exit 1
fi

# 2. Generate the X509 certificate for the Local CA
#openssl req -new -x509 -days 365 -key ca-key.pem -out ca-cert.pem -subj "/C=IN/ST=Karnataka/L=Bangalore/O=Local CA Organisation/OU=IT Department/CN=localcaorganisation.com"
openssl req -new -x509 -days 365 -key ca-key.pem -out ca-cert.crt -subj "/C=IN/ST=Karnataka/L=Bangalore/O=Local CA Organisation/OU=IT Department/CN=localcaorganisation.com"
if [ $? -ne 0 ]; then
    echo 'Failed to generate Local CA X509 certificate'
    exit 1
fi

# 3. Generate the RSA private key for server
openssl genpkey -algorithm RSA -out server-key.pem -pkeyopt rsa_keygen_bits:2048
if [ $? -ne 0 ]; then
    echo 'Failed to generate Server RSA private key'
    exit 1
fi

# 4. Generate certificate request
#openssl req -new -key server-key.pem -out server-request.pem -subj "/C=IN/ST=Karnataka/L=Bangalore/O=Local server Organisation/OU=IT Department/CN=localserverorganisation.com"
openssl req -new -key server-key.pem -out server-request.csr -subj "/C=IN/ST=Karnataka/L=Bangalore/O=Local server Organisation/OU=IT Department/CN=localserverorganisation.com"
if [ $? -ne 0 ]; then
    echo 'Failed to generate certificate request'
    exit 1
fi

# 5. Generate the X509 certificate
#openssl x509 -req -days 365 -in server-request.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem
openssl x509 -req -days 365 -in server-request.csr -CA ca-cert.crt -CAkey ca-key.pem -CAcreateserial -out server-cert.crt -extfile <(printf "subjectAltName=DNS:localhost")
if [ $? -ne 0 ]; then
    echo 'Failed to generate certificate'
    exit 1
fi

# 6. Verify the certificate
#openssl verify -CAfile ca-cert.pem server-cert.pem
openssl verify -CAfile ca-cert.crt server-cert.crt
if [ $? -ne 0 ]; then
    echo 'Failed to verify server certificate'
    exit 1
fi
