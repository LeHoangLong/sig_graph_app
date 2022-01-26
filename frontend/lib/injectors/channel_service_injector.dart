import 'package:frontend/services/channel_service_grpc.dart';
import 'package:frontend/services/channel_service_i.dart';
import 'package:frontend/services/connection_service_grpc.dart';
import 'package:frontend/services/connection_service_i.dart';
import 'package:grpc/grpc_connection_interface.dart';
import 'package:injector/injector.dart';

void injectChannelService(Injector injector) {
  injector.registerDependency<ChannelServiceI>(() {
    var channel = injector.get<ClientChannel>();
    return ChannelServiceGrpc(channel: channel);
  });
}
