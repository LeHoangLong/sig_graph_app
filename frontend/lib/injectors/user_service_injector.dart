import 'package:frontend/services/user_service_grpc.dart';
import 'package:frontend/services/user_service_i.dart';
import 'package:grpc/grpc_connection_interface.dart';
import 'package:injector/injector.dart';

void injectUserService(Injector injector) {
  injector.registerDependency<UserServiceI>(() {
    var channel = injector.get<ClientChannel>();
    return UserServiceGrpc(channel: channel);
  });
}
