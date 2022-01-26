import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/blocs/user_bloc.dart';

class CreateUserPage extends StatefulWidget {
  const CreateUserPage({
    Key? key,
  }) : super(key: key);

  @override
  State createState() => CreateUserPageState();
}

class CreateUserPageState extends State<CreateUserPage> {
  var username = TextEditingController();
  var password = TextEditingController();
  var organizationMspId = TextEditingController();
  String certificatePath = "";
  String keyPath = "";
  String errorText = "";
  GlobalKey<FormState> key = GlobalKey<FormState>();

  onCreateButtonClicked() async {
    try {
      await context.read<UserBloc>().createUser(
            username: username.text,
            password: password.text,
            organizationMspId: organizationMspId.text,
            certificatePath: certificatePath,
            keyPath: keyPath,
          );
      setState(() {
        username.clear();
        password.clear();
        organizationMspId.clear();
        certificatePath = "";
        keyPath = "";
        errorText = "Created";
      });
    } catch (exception) {
      errorText = exception.toString();
    }
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    return Form(
      key: key,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('Username'),
              TextField(
                controller: username,
              ),
            ],
          ),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('Password'),
              TextField(
                obscureText: true,
                enableSuggestions: false,
                autocorrect: false,
                controller: password,
              ),
            ],
          ),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('Organization MSP ID'),
              TextField(
                controller: organizationMspId,
              ),
            ],
          ),
          Row(
            children: [
              TextButton(
                onPressed: () async {
                  var result = await FilePicker.platform.pickFiles(
                    allowMultiple: false,
                    type: FileType.custom,
                    allowedExtensions: ['pem'],
                  );
                  setState(() {
                    if (result != null && result.files.isNotEmpty) {
                      certificatePath = result.files[0].path!;
                    }
                  });
                },
                child: const Text('Choose public certificate'),
              ),
              Expanded(child: Text(certificatePath.split('/').last)),
            ],
          ),
          Row(
            children: [
              TextButton(
                onPressed: () async {
                  var result = await FilePicker.platform.pickFiles(
                    allowMultiple: false,
                    type: FileType.custom,
                    allowedExtensions: ['pem'],
                  );
                  setState(() {
                    if (result != null && result.files.isNotEmpty) {
                      keyPath = result.files[0].path!;
                    }
                  });
                },
                child: const Text('Choose certificate key'),
              ),
              Expanded(child: Text(keyPath.split('/').last)),
            ],
          ),
          TextButton(
            onPressed: () {
              if (key.currentState!.validate()) {
                onCreateButtonClicked();
              }
            },
            child: const Text('Create'),
          ),
          Text(errorText),
        ],
      ),
    );
  }
}
