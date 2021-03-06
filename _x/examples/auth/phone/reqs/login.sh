#!/bin/zsh

http -v \
  POST http://localhost:7450/api/auth/phone/login \
  phone_number="+11234567890" verification_code:=1234
