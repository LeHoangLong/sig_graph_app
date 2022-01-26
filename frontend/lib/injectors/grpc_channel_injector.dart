import 'dart:ffi';

import 'package:grpc/grpc.dart';
import 'package:grpc/grpc_connection_interface.dart';
import 'package:injector/injector.dart';

void injectGrpcChannel(Injector injector, String hostname, int port) {
  injector.registerSingleton<ClientChannel>(() {
    return ClientChannel(
      hostname,
      port: port,
      options: const ChannelOptions(credentials: ChannelCredentials.insecure()),
    );
  });
}
