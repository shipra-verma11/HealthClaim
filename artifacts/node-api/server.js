const http =require('http');
const app = require('./app');
const port  = process.env.PORT || 9092;

const server  = http.createServer(app);
server.listen(port,(err,res)=>{console.log('server started')});    
