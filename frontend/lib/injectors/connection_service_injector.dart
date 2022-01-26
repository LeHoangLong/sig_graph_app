import 'package:frontend/services/connection_service_grpc.dart';
import 'package:frontend/services/connection_service_i.dart';
import 'package:grpc/grpc_connection_interface.dart';
import 'package:injector/injector.dart';

void injectConnectionService(Injector injector) {
  injector.registerDependency<ConnectionServiceI>(() {
    var channel = injector.get<ClientChannel>();
    return ConnectionServiceGrpc(channel: channel);
  });
}
