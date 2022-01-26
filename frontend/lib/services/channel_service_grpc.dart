import 'package:frontend/exceptions.dart';
import 'package:frontend/grpc/channel.pbgrpc.dart';
import 'package:frontend/services/channel_service_i.dart';
import 'package:grpc/grpc_connection_interface.dart';

class ChannelServiceGrpc implements ChannelServiceI {
  late ChannelGrpcClient _client;
  ChannelServiceGrpc({
    required ClientChannel channel,
  }) {
    _client = ChannelGrpcClient(channel);
  }

  @override
  Future<List<String>> getChannels() async {
    try {
      var request = GetChannelsRequest();
      var response = await _client.getChannels(request);
      return response.channels;
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
