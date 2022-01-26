import 'package:frontend/blocs/channel_bloc.dart';
import 'package:frontend/blocs/connection_bloc.dart';
import 'package:frontend/services/channel_service_i.dart';
import 'package:frontend/services/connection_service_i.dart';
import 'package:injector/injector.dart';

void injectChannelBloc(Injector injector) {
  injector.registerSingleton<ChannelBloc>(() {
    var connectionService = injector.get<ChannelServiceI>();
    var connectionBloc = injector.get<ConnectionBloc>();
    return ChannelBloc(
      service: connectionService,
      bloc: connectionBloc,
    );
  });
}
