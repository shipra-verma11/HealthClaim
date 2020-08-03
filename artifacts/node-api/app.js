const express = require('express');
const app = express();
const bodyParser = require('body-parser'); 

var cors = require("cors");
app.use(cors());

const routes = require('./api/routes/routes');

app.use(bodyParser.urlencoded({extended : false}));
app.use(bodyParser.json());

app.use('/',routes);

app.use((req,res,next)=>{
    res.status(200).json({
        message : 'It Works!'
    });
});

app.use((error,req,res,next) =>{
    res.status(error.status ||500);
    res.json({
        error : {
            message : error.message
        }
    });
});

module.exports= app;