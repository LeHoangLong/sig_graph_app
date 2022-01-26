import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/blocs/user_bloc.dart';
import 'package:frontend/widgets/create_user_page.dart';
import 'package:frontend/widgets/login_page.dart';

class UserPage extends StatefulWidget {
  const UserPage({
    Key? key,
  }) : super(key: key);

  @override
  State createState() => UserPageState();
}

class UserPageState extends State<UserPage> {
  ScrollController _controller = ScrollController();
  final _signUpKey = GlobalKey();
  final _logInKey = GlobalKey();

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Row(
          children: [
            TextButton(
              onPressed: () {
                if (_logInKey.currentContext != null) {
                  Scrollable.ensureVisible(
                    _logInKey.currentContext!,
                    duration: const Duration(
                      milliseconds: 200,
                    ),
                    curve: Curves.easeInOut,
                  );
                }
              },
              child: const Text('Log in'),
            ),
            TextButton(
              onPressed: () {
                if (_signUpKey.currentContext != null) {
                  Scrollable.ensureVisible(
                    _signUpKey.currentContext!,
                    duration: const Duration(
                      milliseconds: 200,
                    ),
                    curve: Curves.easeInOut,
                  );
                }
              },
              child: const Text('Sign up'),
            ),
          ],
        ),
        Expanded(
          child: LayoutBuilder(
            builder: (context, constraints) {
              return SingleChildScrollView(
                controller: _controller,
                physics: const NeverScrollableScrollPhysics(),
                scrollDirection: Axis.horizontal,
                child: Row(
                  children: [
                    ConstrainedBox(
                      key: _logInKey,
                      constraints: BoxConstraints(
                        maxWidth: constraints.maxWidth,
                      ),
                      child: const LoginPage(),
                    ),
                    ConstrainedBox(
                      key: _signUpKey,
                      constraints: BoxConstraints(
                        maxWidth: constraints.maxWidth,
                      ),
                      child: const CreateUserPage(),
                    ),
                  ],
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}
