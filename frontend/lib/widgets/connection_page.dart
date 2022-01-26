import 'dart:io';

import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/blocs/connection_bloc.dart';
import 'package:frontend/states/connection_state.dart' as custom;

class ConnectionProfilePage extends StatefulWidget {
  const ConnectionProfilePage({
    Key? key,
  }) : super(key: key);

  @override
  State createState() => ConnectionProfilePageState();
}

class ConnectionProfilePageState extends State<ConnectionProfilePage> {
  File? _selectedFile;
  String errorText = "";
  bool inProgress = false;

  onSaveButtonClicked(BuildContext context) async {
    if (_selectedFile == null) {
      errorText = "Please choose a file";
      return;
    } else {
      setState(() {
        inProgress = true;
      });
      try {
        await context
            .read<ConnectionBloc>()
            .saveConnection(_selectedFile!.path);
        setState(() {
          errorText = "";
        });
      } catch (exception) {
        setState(() {
          errorText = exception.toString();
        });
      } finally {
        setState(() {
          inProgress = false;
        });
      }
    }
  }

  onConnectButtonClicked(BuildContext context, String configStr) async {
    if (configStr == "") {
      errorText = "Please save a connection file";
      return;
    }

    setState(() {
      inProgress = true;
    });
    try {
      await context.read<ConnectionBloc>().connect();

      setState(() {
        errorText = "";
      });
    } catch (exception) {
      setState(() {
        errorText = exception.toString();
      });
    } finally {
      setState(() {
        inProgress = false;
      });
    }
  }

  Widget showConnectStatus(custom.ConnectionState state) {
    if (state.connected) {
      return const Icon(Icons.check);
    } else {
      return Container();
    }
  }

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<ConnectionBloc, custom.ConnectionState>(
        builder: (context, state) {
      return Column(
        children: [
          Row(
            children: [
              showConnectStatus(state),
              TextButton(
                onPressed: () async {
                  FilePickerResult? result =
                      await FilePicker.platform.pickFiles(
                    allowMultiple: false,
                    type: FileType.custom,
                    allowedExtensions: ['json'],
                  );
                  if (result != null) {
                    setState(() {
                      _selectedFile = File(result.files.single.path!);
                    });
                  }
                },
                child: inProgress
                    ? const CircularProgressIndicator()
                    : const Text('Load'),
              ),
              Text(_selectedFile?.path.split('/').last ?? ""),
              TextButton(
                child: const Text('Save'),
                onPressed: () => onSaveButtonClicked(context),
              ),
              TextButton(
                onPressed: () {
                  onConnectButtonClicked(
                    context,
                    state.connectionProfile ?? "",
                  );
                },
                child: const Text('Connect'),
              )
            ],
          ),
          Text(errorText),
          Expanded(
            child: SingleChildScrollView(
              child: Text(
                state.connectionProfile ?? "",
              ),
            ),
          ),
        ],
      );
    });
  }
}
