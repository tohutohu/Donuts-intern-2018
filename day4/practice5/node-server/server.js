var PROTO_PATH = __dirname + '/../proto/fizzbuzz.proto';
 
var grpc = require('grpc');
var fizzbuzz_proto = grpc.load(PROTO_PATH).proto;
 
function calcFizzBuzz(call, callback) {
  console.log(call)
	const num = Number(call.request.num)
  let res
	if (num%15 == 0) {
		res = "fizzbuzz"
	} else if (num%5 == 0) {
		res = "buzz"
	} else if (num%3 == 0) {
		res = "fizz"
	} else {
		res = num
	}
  console.log(res)
	callback(null, {res:res});
}
 
function main() {
	const server = new grpc.Server();
	server.addService(fizzbuzz_proto.FizzBuzz.service, {CalcFizzBuzz: calcFizzBuzz});
	server.bind('0.0.0.0:50051', grpc.ServerCredentials.createInsecure());
	server.start();
}
 
main();
