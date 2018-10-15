const PROTO_PATH = __dirname + '/../calculator/calculatorpb/calculator.proto';

const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});
const calc_proto = grpc.loadPackageDefinition(packageDefinition).calculator;

function main() {
  const client = new calc_proto.CalculatorService(
    'localhost:50052',
    grpc.credentials.createInsecure()
  );

  client.Sum({ first_number: 1, second_number: 2 }, (err, res) => {
    console.log('Res from server: ', res);
  });

  const call = client.PrimeNumberDecomposition({ number: 12 });
  call.on('data', res => {
    console.log('yo', res);
  });
}

main();
