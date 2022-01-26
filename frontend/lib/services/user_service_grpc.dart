import 'package:frontend/exceptions.dart';
import 'package:frontend/grpc/user.pbgrpc.dart';
import 'package:frontend/services/user_service_i.dart';
import 'package:grpc/grpc_connection_interface.dart';

class UserServiceGrpc extends UserServiceI {
  late UserGrpcClient _client;
  UserServiceGrpc({
    required ClientChannel channel,
  }) {
    _client = UserGrpcClient(channel);
  }

  @override
  Future<void> createUser({
    required String username,
    required String password,
    required String organizationMspId,
    required String certificatePath,
    required String keyPath,
  }) async {
    try {
      var request = CreateUserRequest(
        username: username,
        password: password,
        organizationMspId: organizationMspId,
        certPath: certificatePath,
        certKey: keyPath,
      );
      await _client.createUser(request);
    } on GrpcError catch (exception) {
      var message = exception.message;
      if (message != null) {
        throw BaseException(message: message);
      } else {
        throw BaseException(message: exception.codeName);
      }
    }
  }

  @override
  Future<void> login({
    required String username,
    required String password,
  }) async {
    try {
      var request = LoginRequest(
        username: username,
        password: password,
      );
      await _client.login(request);
    } on GrpcError catch (exception) {
      var message = exception.message;
      if (message != null) {
        throw BaseException(message: message);
      } else {
        throw BaseException(message: exception.codeName);
      }
    }
  }
}
