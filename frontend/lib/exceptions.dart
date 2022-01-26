class BaseException extends Error {
  String message;
  BaseException({
    this.message = "",
  });

  @override
  String toString() {
    if (stackTrace != null) {
      return "Uncaught error : " +
          message +
          ".\nStack trace: " +
          stackTrace.toString();
    } else {
      return "Uncaught error : " + message;
    }
  }
}
