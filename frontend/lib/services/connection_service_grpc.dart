import 'package:frontend/exceptions.dart';
import 'package:frontend/grpc/connection.pbgrpc.dart';
import 'package:frontend/services/connection_service_i.dart';
import 'package:grpc/grpc_connection_interface.dart';

class ConnectionServiceGrpc implements ConnectionServiceI {
  late ConnectionGrpcClient _client;
  ConnectionServiceGrpc({
    required ClientChannel channel,
  }) {
    _client = ConnectionGrpcClient(channel);
  }

  @override
  Future<String> getConnection() async {
    try {
      var request = GetConnectionProfileRequest();
      var response = await _client.getConnectionProfile(request);
      return response.data;
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
  Future<String> saveConnection(String path) async {
    try {
      var request = SaveConnectionProfileRequest(
        path: path,
      );
      var response = await _client.saveConnectionProfile(request);
      return response.data;
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
  Future<void> connect() async {
    try {
      var request = ConnectRequest();
      await _client.connect(request);
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
