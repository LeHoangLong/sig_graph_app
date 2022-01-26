import 'package:frontend/blocs/connection_bloc.dart';
import 'package:frontend/services/connection_service_i.dart';
import 'package:injector/injector.dart';

void injectConnectionBloc(Injector injector) {
  injector.registerSingleton<ConnectionBloc>(() {
    var connectionService = injector.get<ConnectionServiceI>();
    return ConnectionBloc(connectionService: connectionService);
  });
}
