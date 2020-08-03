const express = require('express');
const router = express.Router();
const chaincodeMethods = require('../controllers/chaincodemethods');

router.get('/', (req,res,next) => {
    res.json({
        result: "ok",
        message: "initialized"
    })
});
router.post('/registerCustomer', chaincodeMethods.registerCustomer);
router.post('/issuePolicy', chaincodeMethods.issuePolicy);
router.post('/claim', chaincodeMethods.claim);
router.post('/queryAllAsset', chaincodeMethods.queryAllAsset);
router.post('/enrollAdmin', chaincodeMethods.enrollAdmin);
router.post('/queryCustomerRecord', chaincodeMethods.queryCustomerRecord);

module.exports = router;