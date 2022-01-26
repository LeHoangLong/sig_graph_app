import 'package:frontend/exceptions.dart';

abstract class ConnectionServiceI {
  Future<String> getConnection();
  Future<String> saveConnection(String path);
  Future<void> connect();
}
