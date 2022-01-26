///
//  Generated code. Do not modify.
//  source: connection.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,unnecessary_const,non_constant_identifier_names,library_prefixes,unused_import,unused_shown_name,return_of_invalid_type,unnecessary_this,prefer_final_fields,deprecated_member_use_from_same_package

import 'dart:core' as $core;
import 'dart:convert' as $convert;
import 'dart:typed_data' as $typed_data;
@$core.Deprecated('Use connectionProfileDescriptor instead')
const ConnectionProfile$json = const {
  '1': 'ConnectionProfile',
  '2': const [
    const {'1': 'data', '3': 1, '4': 1, '5': 9, '10': 'data'},
  ],
};

/// Descriptor for `ConnectionProfile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionProfileDescriptor = $convert.base64Decode('ChFDb25uZWN0aW9uUHJvZmlsZRISCgRkYXRhGAEgASgJUgRkYXRh');
@$core.Deprecated('Use getConnectionProfileRequestDescriptor instead')
const GetConnectionProfileRequest$json = const {
  '1': 'GetConnectionProfileRequest',
};

/// Descriptor for `GetConnectionProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConnectionProfileRequestDescriptor = $convert.base64Decode('ChtHZXRDb25uZWN0aW9uUHJvZmlsZVJlcXVlc3Q=');
@$core.Deprecated('Use saveConnectionProfileRequestDescriptor instead')
const SaveConnectionProfileRequest$json = const {
  '1': 'SaveConnectionProfileRequest',
  '2': const [
    const {'1': 'path', '3': 1, '4': 1, '5': 9, '10': 'path'},
  ],
};

/// Descriptor for `SaveConnectionProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List saveConnectionProfileRequestDescriptor = $convert.base64Decode('ChxTYXZlQ29ubmVjdGlvblByb2ZpbGVSZXF1ZXN0EhIKBHBhdGgYASABKAlSBHBhdGg=');
@$core.Deprecated('Use connectRequestDescriptor instead')
const ConnectRequest$json = const {
  '1': 'ConnectRequest',
};

/// Descriptor for `ConnectRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectRequestDescriptor = $convert.base64Decode('Cg5Db25uZWN0UmVxdWVzdA==');
@$core.Deprecated('Use connectResponseDescriptor instead')
const ConnectResponse$json = const {
  '1': 'ConnectResponse',
};

/// Descriptor for `ConnectResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectResponseDescriptor = $convert.base64Decode('Cg9Db25uZWN0UmVzcG9uc2U=');
