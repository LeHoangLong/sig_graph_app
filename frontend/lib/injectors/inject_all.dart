import 'package:frontend/injectors/channel_bloc_injector.dart';
import 'package:frontend/injectors/channel_service_injector.dart';
import 'package:frontend/injectors/connection_bloc_injector.dart';
import 'package:frontend/injectors/connection_service_injector.dart';
import 'package:frontend/injectors/grpc_channel_injector.dart';
import 'package:frontend/injectors/user_bloc_injector.dart';
import 'package:frontend/injectors/user_service_injector.dart';
import 'package:injector/injector.dart';

void injectAll(Injector injector) {
  injectGrpcChannel(injector, "localhost", 8000);

  /// connection
  injectConnectionService(injector);
  injectConnectionBloc(injector);

  /// user
  injectUserService(injector);
  injectUserBloc(injector);

  /// channel
  injectChannelService(injector);
  injectChannelBloc(injector);
}
