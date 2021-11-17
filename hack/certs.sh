#!/bin/bash
# usage: ./certs.sh <cmd> <output-directory> [DNS]
# example: ./cert.sh genCA .
# example: ./certs.sh genCSR . /CN=unicorn server dev.kubeedge.edgefarm.dev
# example: ./certs.sh genCert . server dev.kubeedge.edgefarm.dev

SECRET=secretPasswordUnicornFancyStyle

genCA() {
    echo generate a certificate authority
    echo $1 $2 $3
    openssl genrsa -des3 -out ${1}/rootCa.key -passout pass:$SECRET 4096
    echo openssl req -x509 -new -nodes -key ${1}/rootCa.key -sha256 -days 3650 -subj '/CN=kubeedge' -out ${1}/rootCa.pem -passin pass:$SECRET
    openssl req -x509 -new -nodes -key ${1}/rootCa.key -sha256 -days 3650 -subj '/CN=kubeedge' -out ${1}/rootCa.pem -passin pass:$SECRET
}

# genCert arguments:
# 1: relative path to rootCa.pem
# 2: relative path to rootCa.key
# 3: output path
# 4: certificate name
# 5: DNS name
# example: genCert ./rootCa.pem ./RootCa.key . server app.example.com
genCert() {
    echo generate a cert
    # openssl x509 -req -extfile <(printf "subjectAltName=DNS:"${3}) -days 365 -in  ${1}/${2}.csr -CA  ${1}/rootCa.pem -CAkey  ${1}/rootCa.key -CAcreateserial -out  ${1}/${2}.pem -passin pass:$SECRET
    openssl x509 -req -extfile <(printf "subjectAltName=DNS:"${5}) -days 365 -in  ${3}/${4}.csr -CA ${1} -CAkey ${2} -CAcreateserial -out  ${3}/${4}.pem -passin pass:$SECRET
}

genCSR() {
    echo generate a csr
    openssl genrsa -out ${1}/${3}.key 4096
    openssl req -new -key ${1}/${3}.key -subj ${2} -addext "subjectAltName = DNS:"${4} -out ${1}/${3}.csr
}

# if ! command openssl &> /dev/null
# then
#     echo opennssl could not be found
#     exit 1
# fi
echo $(pwd)
case ${1} in
    genCA)
        if [ -z ${2} ]; then
            echo Missing output directory
            exit 1
        fi
        genCA ${2}
    ;;
    genCert)
        if [ -z ${2} ]; then
            echo Missing output directory
            exit 1
        fi
        if [ -z ${3} ]; then
            echo You have to add a name to this functionality
            exit 1
        fi
        if [ -z ${4} ]; then
            echo You have to add a name to this functionality
            exit 1
        fi
        if [ -z ${5} ]; then
            echo You have to add a name to this functionality
            exit 1
        fi
        if [ -z ${6} ]; then
            echo You have to add a name to this functionality
            exit 1
        fi
        genCert ${2} ${3} ${4} ${5} ${6}
    ;;
    genCSR)
        if [ -z ${2} ]; then
            echo Missing output directory
            exit 1
        fi
        if [ -z ${3} ]; then
            echo You have to add a name to this functionality
            exit 1
        fi
        if [ -z ${4} ]; then
            echo You have to add a name to this functionality
            exit 1
        fi
        if [ -z ${5} ]; then
            echo Missing FQDN
            exit 1
        fi
        genCSR ${2} ${3} ${4} ${5}
    ;;
    *)
        echo unknown
    ;;
esac
exit 0
