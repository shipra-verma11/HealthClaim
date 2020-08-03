#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $5)
    local CP=$(one_line_pem $6)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${P1PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        -e "s#\${ORGC}#$7#" \
        ccp-template.json 
}

function yaml_ccp {
    local PP=$(one_line_pem $5)
    local CP=$(one_line_pem $6)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${P1PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${ORGC}#$7#" \
        ccp-template.yaml | sed -e $'s/\\\\n/\\\n        /g'
}

ORG=insurance
ORGC=insurance
P0PORT=7051
P1PORT=7056
CAPORT=7054
PEERPEM=channel/crypto-config/peerOrganizations/insurance.example.com/tlsca/tlsca.insurance.example.com-cert.pem
CAPEM=channel/crypto-config/peerOrganizations/insurance.example.com/ca/ca.insurance.example.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM $ORGC)" > connection-insurance.json
echo "$(yaml_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM $ORGC)" > connection-insurance.yaml

ORG=hospital
ORGC=hospital
P0PORT=8051
P1PORT=8056
CAPORT=8054
PEERPEM=channel/crypto-config/peerOrganizations/hospital.example.com/tlsca/tlsca.hospital.example.com-cert.pem
CAPEM=channel/crypto-config/peerOrganizations/hospital.example.com/ca/ca.hospital.example.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM $ORGC)" > connection-hospital.json
echo "$(yaml_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM $ORGC)" > connection-hospital.yaml