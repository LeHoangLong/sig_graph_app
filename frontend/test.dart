import 'package:grpc/grpc.dart';
import 'package:frontend/grpc/connection.pbgrpc.dart';

Future<void> main(List<String> args) async {
  print('start');
  final channel = ClientChannel(
    'localhost',
    port: 8000,
    options: const ChannelOptions(credentials: ChannelCredentials.insecure()),
  );

  final stub = ConnectionGrpcClient(channel);

  final name = args.isNotEmpty ? args[0] : 'world';

  try {
    var response =
        await stub.getConnectionProfile(GetConnectionProfileRequest());
    print('Greeter client received: ${response.data}');
  } catch (e) {
    print('Caught error: $e');
  }
  await channel.shutdown();
}
