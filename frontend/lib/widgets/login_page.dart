import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/blocs/user_bloc.dart';
import 'package:frontend/exceptions.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({
    Key? key,
  }) : super(key: key);

  @override
  State createState() => LoginPageState();
}

class LoginPageState extends State<LoginPage> {
  final TextEditingController _usernameController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
  String errorText = "";

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('Username'),
          TextFormField(
            validator: (value) {
              if (value == null || value.isEmpty) {
                return "Cannot be empty";
              } else {
                return null;
              }
            },
            controller: _usernameController,
          ),
          const Text('Password'),
          TextFormField(
            obscureText: true,
            enableSuggestions: false,
            autocorrect: false,
            validator: (value) {
              if (value == null || value.isEmpty) {
                return "Cannot be empty";
              } else {
                return null;
              }
            },
            controller: _passwordController,
          ),
          TextButton(
            onPressed: () async {
              if (_formKey.currentState?.validate() == true) {
                try {
                  await context.read<UserBloc>().login(
                        username: _usernameController.text,
                        password: _passwordController.text,
                      );
                  errorText = "Logged in";
                } on BaseException catch (exception) {
                  errorText = exception.message;
                } finally {
                  setState(() {});
                }
              }
            },
            child: const Text('Log in'),
          ),
          Text(errorText),
        ],
      ),
    );
  }
}
