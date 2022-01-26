///
//  Generated code. Do not modify.
//  source: channel.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,unnecessary_const,non_constant_identifier_names,library_prefixes,unused_import,unused_shown_name,return_of_invalid_type,unnecessary_this,prefer_final_fields

import 'dart:async' as $async;

import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'channel.pb.dart' as $0;
export 'channel.pb.dart';

class ChannelGrpcClient extends $grpc.Client {
  static final _$getChannels =
      $grpc.ClientMethod<$0.GetChannelsRequest, $0.GetChannelsResponse>(
          '/dashboard.ChannelGrpc/GetChannels',
          ($0.GetChannelsRequest value) => value.writeToBuffer(),
          ($core.List<$core.int> value) =>
              $0.GetChannelsResponse.fromBuffer(value));

  ChannelGrpcClient($grpc.ClientChannel channel,
      {$grpc.CallOptions? options,
      $core.Iterable<$grpc.ClientInterceptor>? interceptors})
      : super(channel, options: options, interceptors: interceptors);

  $grpc.ResponseFuture<$0.GetChannelsResponse> getChannels(
      $0.GetChannelsRequest request,
      {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$getChannels, request, options: options);
  }
}

abstract class ChannelGrpcServiceBase extends $grpc.Service {
  $core.String get $name => 'dashboard.ChannelGrpc';

  ChannelGrpcServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.GetChannelsRequest, $0.GetChannelsResponse>(
            'GetChannels',
            getChannels_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetChannelsRequest.fromBuffer(value),
            ($0.GetChannelsResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.GetChannelsResponse> getChannels_Pre($grpc.ServiceCall call,
      $async.Future<$0.GetChannelsRequest> request) async {
    return getChannels(call, await request);
  }

  $async.Future<$0.GetChannelsResponse> getChannels(
      $grpc.ServiceCall call, $0.GetChannelsRequest request);
}
