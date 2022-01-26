import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/services/connection_service_i.dart';
import 'package:frontend/states/connection_state.dart';

class ConnectionBloc extends Cubit<ConnectionState> {
  ConnectionServiceI connectionService;
  String _connection = "";
  bool isConnected = false;
  ConnectionBloc({
    required this.connectionService,
  }) : super(ConnectionState()) {
    _getConnection();
  }

  Future<String> _getConnection() async {
    _connection = await connectionService.getConnection();
    if (_connection != "") {
      connect();
    }
    emit(ConnectionState(
      connectionProfile: _connection,
    ));
    return _connection;
  }

  Future<void> saveConnection(String path) async {
    _connection = await connectionService.saveConnection(path);
    emit(ConnectionState(
      connectionProfile: _connection,
    ));
  }

  Future<void> connect() async {
    await connectionService.connect();
    isConnected = true;
    emit(ConnectionState(
      connectionProfile: _connection,
      connected: isConnected,
    ));
  }
}
