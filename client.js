const http2 = require('http2');
const fs = require('fs');
process.env['NODE_TLS_REJECT_UNAUTHORIZED'] = 0
const client = http2.connect('https://localhost:8100');

var buffer = JSON.stringify({
    "Login":{
      "username":"sudeep.dasgupta",
      "password":"kaihiwatari"
    },
    "GetProfile":{},
    "SchemesMaster":{
      "category":"EQUITY"
    }
  });

// Must not specify the ':path' and ':scheme' headers
// for CONNECT requests or an error will be thrown.
const req = client.request({
	//':path': '/streaming',
  	// ':method': 'POST'

  	[http2.constants.HTTP2_HEADER_SCHEME]: "http",

    [http2.constants.HTTP2_HEADER_METHOD]: http2.constants.HTTP2_METHOD_POST,

    [http2.constants.HTTP2_HEADER_PATH]: `/streaming`,

    "Content-Type": "application/json",

    "Content-Length": buffer.length,

});

req.on('response', (headers) => {
  console.log(headers[http2.constants.HTTP2_HEADER_STATUS]);
});
let data = '';
req.setEncoding('utf8');
req.on('data', (chunk) => {
	data += chunk;
	console.log("chunk: "+chunk);
});
req.on('end', () => {
  console.log(`The server says: ${data}`);
  client.close();
});
req.end(buffer);