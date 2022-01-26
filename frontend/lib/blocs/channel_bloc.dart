import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/blocs/connection_bloc.dart';
import 'package:frontend/exceptions.dart';
import 'package:frontend/services/channel_service_i.dart';
import 'package:frontend/states/channel_state.dart';
import 'package:frontend/states/connection_state.dart';
import 'package:frontend/states/operation_state.dart';

class ChannelBloc extends Cubit<ChannelState> {
  ChannelServiceI service;
  var channels = <String>[];
  String? _previousConnectionProfile;

  ChannelBloc({
    required this.service,
    required ConnectionBloc bloc,
  }) : super(ChannelState(channels: [], state: OperationState.init)) {
    if (bloc.state.connected) {
      _getChannels(bloc.state);
    }

    bloc.stream.listen(_getChannels);
  }

  _getChannels(ConnectionState state) async {
    emit(ChannelState(
      channels: channels,
      state: OperationState.inProgress,
    ));
    try {
      if (state.connected &&
          state.connectionProfile != null &&
          state.connectionProfile != _previousConnectionProfile) {
        var channels = await service.getChannels();
        emit(ChannelState(
          channels: channels,
          state: OperationState.idle,
        ));
      } else if (!state.connected) {
        channels = [];
        emit(
          ChannelState(
            channels: channels,
            state: OperationState.idle,
          ),
        );
      }
    } on BaseException catch (exception) {
      emit(
        ChannelState(
          channels: channels,
          error: exception,
          state: OperationState.idle,
        ),
      );
    } catch (exception) {
      emit(
        ChannelState(
          channels: channels,
          error: BaseException(message: exception.toString()),
          state: OperationState.idle,
        ),
      );
    } finally {
      _previousConnectionProfile = state.connectionProfile;
    }
  }
}
