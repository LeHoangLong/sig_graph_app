abstract class UserServiceI {
  Future<void> createUser({
    required String username,
    required String password,
    required String organizationMspId,
    required String certificatePath,
    required String keyPath,
  });
  Future<void> login({
    required String username,
    required String password,
  });
}
