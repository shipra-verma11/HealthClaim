const {
    FileSystemWallet,
    Gateway,
    X509WalletMixin
  } = require("fabric-network");
  const path = require("path");
  const fs = require('fs');
  const FabricCAServices = require('fabric-ca-client');
  const ccpPath = path.resolve(__dirname, "..", "..","..", "connection-insurance.json");
  const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
  const ccp = JSON.parse(ccpJSON);
  
exports.enrollAdmin = async (req, res) => {
  console.log(`Wallet path: ${walletPath}`);

     try {
     // Create a new CA client for interacting with the CA.
     const caInfo = ccp.certificateAuthorities['ca.insurance.example.com'];
     const caTLSCACerts = caInfo.tlsCACerts.pem;
     const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);
     // Create a new file system based wallet for managing identities.
     const walletPath = path.join(process.cwd(), 'wallet', 'insurance');
     const wallet = new FileSystemWallet(walletPath);
     console.log(`Wallet path: ${walletPath}`);
     // Check to see if we've already enrolled the admin user.
     const adminExists = await wallet.exists('insurance_admin');
     if (adminExists) {
         console.log('An identity for the admin user "admin" already exists in the wallet');
         return;
     }
     // Enroll the admin user, and import the new identity into the wallet.
     const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
     const identity = X509WalletMixin.createIdentity('InsuranceMSP', enrollment.certificate, enrollment.key.toBytes());
     await wallet.import('insurance_admin', identity);
     console.log('Successfully enrolled admin user "admin" and imported it into the wallet');
     res.json({
         result: "ok",
         message: 'Successfully enrolled admin user "admin" and imported it into the wallet'
       });
 } catch (error) {
     console.error(`Failed to enroll admin : ${error}`);
     res.json({
         result: "failed",
         message: `Failed to enroll admin : ${error}`
       });
      }
}

exports.registerCustomer = async (req, res) => {
    let name = req.body.name;
    let age = req.body.age;
    let sex = req.body.sex;
    let pan = req.body.pan;
    let companyRole = "insurance";
    
    try{
        const wallet = new FileSystemWallet("./wallet/" + companyRole);
        const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
        // Check to see if we've already enrolled the user.
        const fabricIdentity= companyRole + "_admin";
        const userExists = await wallet.exists(fabricIdentity);
        if (!userExists) {
            console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
            console.log('Add the identity to wallet');
            return;
        }
        console.log("connecting to getway");
        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });
        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');
        console.log("connected to getway");
        // Get the contract from the network.
        const contract = network.getContract('healthclaim');
        console.log("received contract");

        const result= await contract.submitTransaction('registerCustomer',name, age, sex, pan);
        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

        res.json({
            status: "success",
            message: `Customer registered for PAN: ${pan} successfully`,
            data:result.toString()
          });

    }
    catch (error) {
        res.json({
          status: "failed",
          message: `Failed to register customer ${pan}: ${error}`
        });
    }    
}

exports.issuePolicy = async (req, res) => {
  let customerId = req.body.customerId;
  let pan = req.body.pan;
  let companyRole = "insurance";

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole+"_admin"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });
      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('mychannel');
      // Get the contract from the network.
      const contract = network.getContract('healthclaim');
      const result=await contract.submitTransaction('issuePolicy', customerId, pan);
      console.log('Transaction has been submitted');
      // Disconnect from the gateway.
      await gateway.disconnect();
      res.json({
          status: "success",
          message: `Policy successfully issued`,
          data:result.toString()
        });
  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to issue of policy: ${error}`
      });
  }
}

exports.claim = async (req, res) => {
  let customerId = req.body.customerId;
  let policyId = req.body.policyId;
  let companyRole = "insurance";

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole+"_admin"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }
      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });
      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('mychannel');
      // Get the contract from the network.
      const contract = network.getContract('healthclaim');
      const result=await contract.submitTransaction('claim', customerId, policyId);
      console.log('Transaction has been submitted');
      // Disconnect from the gateway.
      await gateway.disconnect();
      res.json({
          status: "success",
          message: `Policy claimed successfully`,
          data:result.toString()
        });
  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to claim: ${error}`,
      });
    }
}

exports.queryCustomerRecord = async (req, res) => {
  let pan = req.body.pan;
  let companyRole = "insurance";

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole+"_admin"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }
      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });
      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('mychannel');
      // Get the contract from the network.
      const contract = network.getContract('healthclaim');      
      const result=await contract.submitTransaction('queryCustomerRecord',pan);
      console.log('Transaction has been submitted');
      // Disconnect from the gateway.
      await gateway.disconnect();
      res.json({
          status: "success",
          message: `Customer record:`,
          data:result.toString()
        });
  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to return customer: ${error}`
      });
    }
}

exports.queryAllAsset = async (req, res) => {
  let companyRole = "insurance";

  try{
      const wallet = new FileSystemWallet("./wallet/" + companyRole);
      const connectionProfilePath = path.resolve(__dirname, "..", "..","..", "connection-"+companyRole+".json");
      // Check to see if we've already enrolled the user.
      const fabricIdentity= companyRole+"_admin"
      const userExists = await wallet.exists(fabricIdentity);
      if (!userExists) {
          console.log(`An identity for the company ${companyRole} does not exist in the wallet`);
          console.log('Add the identity to wallet');
          return;
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(connectionProfilePath, { wallet, identity: fabricIdentity, discovery: { enabled: true, asLocalhost: true } });
      // Get the network (channel) our contract is deployed to.
      const network = await gateway.getNetwork('mychannel');
      // Get the contract from the network.
      const contract = network.getContract('healthclaim');
      const result=await contract.submitTransaction('queryAllAsset');
      console.log('Transaction has been submitted');
      // Disconnect from the gateway.
      await gateway.disconnect();
      res.json({
          status: "success",
          message: `Result:`,
          data: result.toString()
        });

  }
  catch (error) {
      res.json({
        status: "failed",
        message: `Failed to fetch customer: ${error}`
      });
  }
}
