#!/bin/bash

rm *.pem

openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=SomeOrganization/OU=Test/CN=*.localhost.local/emailAddress=admin@localhost"
openssl x509 -in ca-cert.pem -noout -text
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=SomeOrganization/OU=Test/CN=*.localhost.local/emailAddress=admin@localhost"
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile config/server.cnf
openssl x509 -in server-cert.pem -noout -text
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=SomeOrganization/OU=Test/CN=*.localhost.local/emailAddress=admin@localhost"
openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile config/client.cnf
openssl x509 -in client-cert.pem -noout -text