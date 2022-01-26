///
//  Generated code. Do not modify.
//  source: connection.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,unnecessary_const,non_constant_identifier_names,library_prefixes,unused_import,unused_shown_name,return_of_invalid_type,unnecessary_this,prefer_final_fields

import 'dart:async' as $async;

import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'connection.pb.dart' as $0;
export 'connection.pb.dart';

class ConnectionGrpcClient extends $grpc.Client {
  static final _$getConnectionProfile =
      $grpc.ClientMethod<$0.GetConnectionProfileRequest, $0.ConnectionProfile>(
          '/dashboard.ConnectionGrpc/GetConnectionProfile',
          ($0.GetConnectionProfileRequest value) => value.writeToBuffer(),
          ($core.List<$core.int> value) =>
              $0.ConnectionProfile.fromBuffer(value));
  static final _$saveConnectionProfile =
      $grpc.ClientMethod<$0.SaveConnectionProfileRequest, $0.ConnectionProfile>(
          '/dashboard.ConnectionGrpc/SaveConnectionProfile',
          ($0.SaveConnectionProfileRequest value) => value.writeToBuffer(),
          ($core.List<$core.int> value) =>
              $0.ConnectionProfile.fromBuffer(value));
  static final _$connect =
      $grpc.ClientMethod<$0.ConnectRequest, $0.ConnectResponse>(
          '/dashboard.ConnectionGrpc/Connect',
          ($0.ConnectRequest value) => value.writeToBuffer(),
          ($core.List<$core.int> value) =>
              $0.ConnectResponse.fromBuffer(value));

  ConnectionGrpcClient($grpc.ClientChannel channel,
      {$grpc.CallOptions? options,
      $core.Iterable<$grpc.ClientInterceptor>? interceptors})
      : super(channel, options: options, interceptors: interceptors);

  $grpc.ResponseFuture<$0.ConnectionProfile> getConnectionProfile(
      $0.GetConnectionProfileRequest request,
      {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$getConnectionProfile, request, options: options);
  }

  $grpc.ResponseFuture<$0.ConnectionProfile> saveConnectionProfile(
      $0.SaveConnectionProfileRequest request,
      {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$saveConnectionProfile, request, options: options);
  }

  $grpc.ResponseFuture<$0.ConnectResponse> connect($0.ConnectRequest request,
      {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$connect, request, options: options);
  }
}

abstract class ConnectionGrpcServiceBase extends $grpc.Service {
  $core.String get $name => 'dashboard.ConnectionGrpc';

  ConnectionGrpcServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.GetConnectionProfileRequest,
            $0.ConnectionProfile>(
        'GetConnectionProfile',
        getConnectionProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetConnectionProfileRequest.fromBuffer(value),
        ($0.ConnectionProfile value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SaveConnectionProfileRequest,
            $0.ConnectionProfile>(
        'SaveConnectionProfile',
        saveConnectionProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SaveConnectionProfileRequest.fromBuffer(value),
        ($0.ConnectionProfile value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ConnectRequest, $0.ConnectResponse>(
        'Connect',
        connect_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ConnectRequest.fromBuffer(value),
        ($0.ConnectResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.ConnectionProfile> getConnectionProfile_Pre(
      $grpc.ServiceCall call,
      $async.Future<$0.GetConnectionProfileRequest> request) async {
    return getConnectionProfile(call, await request);
  }

  $async.Future<$0.ConnectionProfile> saveConnectionProfile_Pre(
      $grpc.ServiceCall call,
      $async.Future<$0.SaveConnectionProfileRequest> request) async {
    return saveConnectionProfile(call, await request);
  }

  $async.Future<$0.ConnectResponse> connect_Pre(
      $grpc.ServiceCall call, $async.Future<$0.ConnectRequest> request) async {
    return connect(call, await request);
  }

  $async.Future<$0.ConnectionProfile> getConnectionProfile(
      $grpc.ServiceCall call, $0.GetConnectionProfileRequest request);
  $async.Future<$0.ConnectionProfile> saveConnectionProfile(
      $grpc.ServiceCall call, $0.SaveConnectionProfileRequest request);
  $async.Future<$0.ConnectResponse> connect(
      $grpc.ServiceCall call, $0.ConnectRequest request);
}
