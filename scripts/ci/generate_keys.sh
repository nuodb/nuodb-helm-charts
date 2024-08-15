#!/usr/bin/env bash

# (C) Copyright NuoDB, Inc. 2019-2021  All Rights Reserved
# This file is licensed under the BSD 3-Clause License.
# See https://github.com/nuodb/nuodb-helm-charts/blob/master/LICENSE

ME=`basename $0`
SCRIPT_DIR=`dirname $0`

: ${SELF_ROOT:=`dirname ${SCRIPT_DIR}`}
pushd ${SELF_ROOT} >/dev/null
SELF_ROOT=`pwd`
popd >/dev/null

: ${KEYS_DIR:="${SELF_ROOT}/../keys"}
: ${PASSWORD:="changeIt"}

# nuo way...
if [ ! -d ${KEYS_DIR} ]; then
    docker run --rm -d --name create-tls-keys nuodb/nuodb:5.0 -- tail -f /dev/null
    docker exec -it create-tls-keys bash -c "mkdir /tmp/keys && cd /tmp/keys && DEFAULT_PASSWORD=changeIt setup-keys.sh"
    docker cp create-tls-keys:/tmp/keys ${KEYS_DIR}
    docker stop create-tls-keys
fi

# uncomment the following to test with OpenSSL generated certs...
# rm -fr ${KEYS_DIR}

# standard way...
if [ ! -d ${KEYS_DIR} ]; then

    mkdir -p ${KEYS_DIR}

    # -----------------------------------------------------
    # CA.CERT

    cat > ${KEYS_DIR}/openssl.cnf <<'EOF'
[ req ]
default_md              = sha256
distinguished_name      = req_distinguished_name
x509_extensions = usr_cert
[ req_distinguished_name ]
[ usr_cert ]
basicConstraints=critical,CA:TRUE
subjectKeyIdentifier=hash
EOF

    # References:
    # - https://www.digitalocean.com/community/tutorials/openssl-essentials-working-with-ssl-certificates-private-keys-and-csrs

    # OPTION 1: (works)
    #   generate a private key and certificate

    # openssl req -config openssl.cnf -newkey rsa:2048 -nodes -x509 -days 365 \
    #     -subj '/C=US/ST=MA/L=Boston/O=NuoDB/OU=Eng/CN=ca.nuodb.com' \
    #     -keyout ca.key -x509 -days 365 -out ca.cert

    # OPTIONS 2:
    #   generate a private key and csr
    #   generate a certificate from a private key and csr

    openssl req -newkey rsa:2048 -nodes \
        -subj '/C=US/ST=MA/L=Boston/O=NuoDB/OU=Eng/CN=ca.nuodb.com' \
        -keyout ${KEYS_DIR}/ca.key -out ${KEYS_DIR}/ca.csr

    # test
    # openssl req -text -in ${KEYS_DIR}/ca.csr

    cat > ${KEYS_DIR}/ssl-extensions-x509.cnf <<'EOF'
[v3_ca]
basicConstraints=critical,CA:TRUE
subjectKeyIdentifier=hash
EOF

    openssl x509 -extensions v3_ca -extfile ${KEYS_DIR}/ssl-extensions-x509.cnf -req \
        -in ${KEYS_DIR}/ca.csr -signkey ${KEYS_DIR}/ca.key -days 365 -out ${KEYS_DIR}/ca.cert

    # -----------------------------------------------------
    # nuocmd.pem

    openssl req -newkey rsa:2048 -nodes \
        -subj '/C=US/ST=MA/L=Boston/CN=nuocmd.nuodb.com/O=NuoDB/OU=Eng' \
        -keyout ${KEYS_DIR}/client.key -out ${KEYS_DIR}/client.csr

    # test
    # openssl req -text -in ${KEYS_DIR}/client.csr

    cat > ${KEYS_DIR}/ssl-extensions-x509.cnf <<'EOF'
[v3_ca]
subjectKeyIdentifier=hash
EOF

    openssl x509 -extensions v3_ca -extfile ${KEYS_DIR}/ssl-extensions-x509.cnf -req \
        -in ${KEYS_DIR}/client.csr -CA ${KEYS_DIR}/ca.cert -CAkey ${KEYS_DIR}/ca.key \
        -CAcreateserial -days 1024 -sha256 -out ${KEYS_DIR}/client.cert

    cat ${KEYS_DIR}/client.cert > ${KEYS_DIR}/nuocmd.pem
    cat ${KEYS_DIR}/client.key >> ${KEYS_DIR}/nuocmd.pem

    # -----------------------------------------------------
    # nuoadmin.p12

    openssl req -newkey rsa:2048 -nodes \
        -subj '/C=US/ST=MA/L=Boston/CN=nuoadmin.nuodb.com/O=NuoDB/OU=Eng' \
        -keyout ${KEYS_DIR}/server.key -out ${KEYS_DIR}/server.csr

    cat > ${KEYS_DIR}/ssl-extensions-x509.cnf <<'EOF'
[v3_ca]
authorityKeyIdentifier=keyid,issuer
basicConstraints=critical,CA:TRUE
subjectKeyIdentifier=hash
EOF

    openssl x509 -extensions v3_ca -extfile ${KEYS_DIR}/ssl-extensions-x509.cnf -req \
        -in ${KEYS_DIR}/server.csr -CA ${KEYS_DIR}/ca.cert -CAkey ${KEYS_DIR}/ca.key \
        -CAcreateserial -days 1024 -sha256 -out ${KEYS_DIR}/server.cert

    openssl pkcs12 -inkey ${KEYS_DIR}/server.key -in ${KEYS_DIR}/server.cert -export \
        -out ${KEYS_DIR}/nuoadmin.p12 -passout pass:${PASSWORD}

    keytool -importcert -file ${KEYS_DIR}/ca.cert -trustcacerts -alias ca-cert \
        -keystore ${KEYS_DIR}/nuoadmin-truststore.p12 -storepass ${PASSWORD} -storetype pkcs12 -noprompt
    keytool -importcert -file ${KEYS_DIR}/server.cert -alias nuoadmin \
        -keystore ${KEYS_DIR}/nuoadmin-truststore.p12 -storepass ${PASSWORD} -storetype PKCS12 -noprompt
    keytool -importcert -file ${KEYS_DIR}/client.cert -alias nuocmd \
        -keystore ${KEYS_DIR}/nuoadmin-truststore.p12 -storepass ${PASSWORD} -storetype PKCS12 -noprompt
fi

# save off the OOTB certficates for comparisons in tests...
docker run --rm -it nuodb/nuodb:5.0 nuocmd show certificate --keystore /etc/nuodb/keys/nuoadmin.p12 --store-password changeIt > ${KEYS_DIR}/default.certificate

# (aka ca.cert)
kubectl delete secret nuodb-ca-cert
kubectl create secret generic nuodb-ca-cert \
  --from-file=ca.cert=${KEYS_DIR}/ca.cert

# The PEM file containing the certificate used to verify admin certificates (aka nuocmd.pem)
kubectl delete secret nuodb-client-pem
kubectl create secret generic nuodb-client-pem \
  --from-file=nuocmd.pem=${KEYS_DIR}/nuocmd.pem

# The keystore for the NuoDB Admin process; contains only the admin key and
# certificate (aka nuoadmin.p12).
kubectl delete secret nuodb-keystore
kubectl create secret generic nuodb-keystore \
  --from-file=nuoadmin.p12=${KEYS_DIR}/nuoadmin.p12 --from-literal=password=${PASSWORD}

# Contains the certificate used to verify admin certificates and the client
# certificate (aka nuoadmin-truststore.p12).
kubectl delete secret nuodb-truststore
kubectl create secret generic nuodb-truststore \
  --from-file=nuoadmin-truststore.p12=${KEYS_DIR}/nuoadmin-truststore.p12 --from-literal=password=${PASSWORD}
