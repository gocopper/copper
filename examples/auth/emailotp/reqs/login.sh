#!/bin/zsh

http -v \
  POST http://localhost:7450/api/auth/email-otp/login \
  email="test@email.com" verification_code:=9315
