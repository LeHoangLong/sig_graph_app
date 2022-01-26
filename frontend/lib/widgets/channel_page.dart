import 'package:flutter/widgets.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/blocs/channel_bloc.dart';
import 'package:frontend/states/channel_state.dart';

class ChannelPage extends StatefulWidget {
  const ChannelPage({
    Key? key,
  }) : super(key: key);

  @override
  State createState() => ChannelPageState();
}

class ChannelPageState extends State<ChannelPage> {
  List<Widget> displayChannels(ChannelState state) {
    var ret = <Widget>[];
    for (var channel in state.channels) {
      ret.add(Text(channel));
    }
    return ret;
  }

  Widget displayErrorText(ChannelState state) {
    if (state.error == null) {
      return Container();
    } else {
      return Text(state.error.toString());
    }
  }

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<ChannelBloc, ChannelState>(
      builder: (context, state) {
        return ListView(
          children: displayChannels(state)
            ..add(
              displayErrorText(state),
            ),
        );
      },
    );
  }
}
