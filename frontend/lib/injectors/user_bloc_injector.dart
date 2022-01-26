import 'package:frontend/blocs/user_bloc.dart';
import 'package:frontend/services/user_service_i.dart';
import 'package:injector/injector.dart';

void injectUserBloc(Injector injector) {
  injector.registerSingleton<UserBloc>(() {
    var userService = injector.get<UserServiceI>();
    return UserBloc(userService: userService);
  });
}
