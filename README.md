# HealthClaim
HealthClaim is tool to track the claim records and allows individual to self claim from the portal.

## Steps to run the network and API:
Note: Make sure to re-generate cryptogen and configtx. `configtx.yaml` and `cryptogen.yaml` is present under `./artifacts/channel` (remove .tx and .block files if any exists)
```
cd ./artifacts/channel
docker-compose up -d
```
For node API
```
cd ./artifacts/node-api
npm install
node server.js
```

## Information on CLI and NODE API:
**1. CLI:**
    
        Functions available within chaincode:
            - registerCustomer(register customer within insurance organization)
            - issuePolicy (issue policy if customer is registered)
            - claim (claims policy if customer and policy exists and not-claimed)
            - queryCustomerRecord (fetch customer record if available)
            - queryAllAsset (fetch all record)

***CLI commands:***

`registerCustomer` (ARGS: Name string, Age string, Sex string, Pan string):
```
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n healthclaim -c '{"Args":["registerCustomer","Ragul","23","male","RAGPAN123"]}'
```
`issuePolicy` (ARGS: customerId string, Pan string):
```
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n healthclaim -c '{"Args":["issuePolicy","1","RAGPAN123"]}'
```
`queryAllAsset`
```
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n healthclaim -c '{"Args":["queryAllAsset"]}'
```
`claim` (ARGS: customerId string, policyId string):
```
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n healthclaim -c '{"Args":["claim","1","AZ56"]}'
```
`queryCustomerRecord` (ARGS: Pan string):
```
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n healthclaim -c '{"Args":["queryCustomerRecord", "RAGPAN123"]}'
```


**2. NODE API**

`enrollAdmin`
```
POST http://localhost:9092/enrollAdmin
```

`registerCustomer`
```
POST http://localhost:9092/registerCustomer

JSON PAYLOAD:
{
    "name": "Ragul",
    "age": "23",
    "sex": "male",
    "pan": "RAGPAN123",
    "companyRole": "insurance"
}
```
`issuePolicy`
```
POST http://localhost:9092/issuePolicy

JSON PAYLOAD:
{
    "customerId": "1",
    "pan": "RAGPAN123"
}
```
`queryAllAsset`
```
POST http://localhost:9092/queryAllAsset
```
`claim`
```
POST http://localhost:9092/claim

JSON PAYLOAD:
{
    "customerId":"1",
    "policyId":"AZ60"
}
```
`queryCustomerRecord`
```
POST http://localhost:9092/queryCustomerRecord

JSON PAYLOAD:
{
    "pan": "RAGPAN123"
}
```