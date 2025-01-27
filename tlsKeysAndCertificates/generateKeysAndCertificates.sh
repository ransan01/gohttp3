#!/bin/bash

# openssl V3

if ls *.pem 1> /dev/null 2>&1; then
    echo "File with .pem extension exists in the current directory."
    echo "This will overwrite the existing files. Either delete it or move it to other directory."
    exit 1
fi
if ls *.crt 1> /dev/null 2>&1; then
    echo "File with .crt extension exists in the current directory."
    echo "This will overwrite the existing files. Either delete it or move it to other directory."
    exit 1
fi
if ls *.csr 1> /dev/null 2>&1; then
    echo "File with .csr extension exists in the current directory."
    echo "This will overwrite the existing files. Either delete it or move it to other directory."
    exit 1
fi

# 1. Generate a RSA private key for the Local CA
openssl genpkey -algorithm RSA -out local-ca-key.pem -pkeyopt rsa_keygen_bits:2048
if [ $? -ne 0 ]; then
    echo 'Failed to generate Local CA RSA private key'
    exit 1
fi

# 2. Generate the X509 certificate for the Local CA
openssl req -new -x509 -days 365 -key local-ca-key.pem -out local-ca-cert.crt -subj "/C=ex/ST=example_st/L=example_l/O=example ca organisation/OU=example_ou/CN=www.example.com/emailAddress=someone@example.com"
if [ $? -ne 0 ]; then
    echo 'Failed to generate Local CA X509 certificate'
    exit 1
fi

# 3. Generate the RSA private key for server
openssl genpkey -algorithm RSA -out my-server-key.pem -pkeyopt rsa_keygen_bits:2048
if [ $? -ne 0 ]; then
    echo 'Failed to generate Server RSA private key'
    exit 1
fi

# 4. Generate certificate request
openssl req -new -key my-server-key.pem -out my-server-certificate-signing-request.csr -subj "/C=ex/ST=example_st/L=example_l/O=example server organisation/OU=example_ou/CN=www.example.com/emailAddress=someone@example.com"
if [ $? -ne 0 ]; then
    echo 'Failed to generate certificate request'
    exit 1
fi

# 5. Generate the X509 certificate
openssl x509 -req -days 365 -in my-server-certificate-signing-request.csr -CA local-ca-cert.crt -CAkey local-ca-key.pem -CAcreateserial -out my-server-cert.crt -extfile <(printf "subjectAltName=DNS:localhost")
if [ $? -ne 0 ]; then
    echo 'Failed to generate certificate'
    exit 1
fi

# 6. Verify the certificate
openssl verify -CAfile local-ca-cert.crt my-server-cert.crt
if [ $? -ne 0 ]; then
    echo 'Failed to verify server certificate'
    exit 1
fi


# 7. Install local-ca-cert.crt as Trusted CA based on your Operating System