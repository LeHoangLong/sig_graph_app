import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/services/user_service_i.dart';
import 'package:frontend/states/user_state.dart';

class UserBloc extends Cubit<UserState> {
  UserServiceI userService;
  bool isLoggedIn = false;

  UserBloc({
    required this.userService,
  }) : super(UserState(isLoggedIn: false));

  createUser({
    required String username,
    required String password,
    required String organizationMspId,
    required String certificatePath,
    required String keyPath,
  }) async {
    await userService.createUser(
      username: username,
      password: password,
      organizationMspId: organizationMspId,
      certificatePath: certificatePath,
      keyPath: keyPath,
    );
  }

  login({
    required String username,
    required String password,
  }) async {
    try {
      await userService.login(username: username, password: password);
      isLoggedIn = true;
      emit(UserState(isLoggedIn: isLoggedIn));
    } catch (exception) {
      /// Do not change login state if login fail
    }
  }
}
