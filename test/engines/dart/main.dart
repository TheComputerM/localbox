import 'dart:io';
import 'dart:convert';

void main() {
  stdin.transform(utf8.decoder).listen((data) {
    stdout.write(data);
  });
}