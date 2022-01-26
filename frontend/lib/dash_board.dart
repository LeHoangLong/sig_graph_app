import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:frontend/grpc/connection.pb.dart';
import 'package:frontend/widgets/channel_page.dart';
import 'package:frontend/widgets/connection_page.dart';
import 'package:frontend/widgets/user_page.dart';

class DashBoard extends StatefulWidget {
  @override
  State createState() => DashBoardState();
}

class _Tile extends StatelessWidget {
  final String title;
  final void Function() onClick;

  _Tile({
    required this.title,
    required this.onClick,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 25, horizontal: 10),
      child: TextButton(
        child: Text(
          title,
        ),
        onPressed: onClick,
      ),
    );
  }
}

class DashBoardState extends State<DashBoard> {
  final Widget connectionPage = const ConnectionProfilePage();
  final Widget userPage = const UserPage();
  final Widget channelPage = const ChannelPage();
  late Widget child;

  @override
  void initState() {
    super.initState();
    child = connectionPage;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Row(
        children: [
          Column(
            children: [
              _Tile(
                title: 'User',
                onClick: () {
                  setState(() {
                    child = userPage;
                  });
                },
              ),
              _Tile(
                title: 'Connection',
                onClick: () {
                  setState(() {
                    child = connectionPage;
                  });
                },
              ),
              _Tile(
                title: 'Channels',
                onClick: () {
                  setState(() {
                    child = channelPage;
                  });
                },
              ),
            ],
          ),
          Expanded(
            child: child,
          ),
        ],
      ),
    );
  }
}
