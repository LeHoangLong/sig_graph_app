import 'package:frontend/states/operation_state.dart';

class ChannelState {
  final Error? error;
  final List<String> channels;
  final OperationState state;
  ChannelState({
    required this.channels,
    this.error,
    required this.state,
  });
}
