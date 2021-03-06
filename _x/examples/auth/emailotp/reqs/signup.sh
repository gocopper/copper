#!/bin/zsh

http -v \
  POST http://localhost:7450/api/auth/email-otp/signup \
  email="test@email.com"
