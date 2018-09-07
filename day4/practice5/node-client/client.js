var PROTO_PATH = __dirname + '/../proto/fizzbuzz.proto';
 
var grpc = require('grpc');
var fizzbuzz_proto = grpc.load(PROTO_PATH).proto;
 
function main() {
  var client = new fizzbuzz_proto.FizzBuzz('localhost:50051',
                                       grpc.credentials.createInsecure());
	const num = process.argv[2]
  client.CalcFizzBuzz({num: num}, function(err, response) {
    console.log(response.res);
  });
}
 
main()
