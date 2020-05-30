#!/bin/zsh

http -v \
  -a "user_uuid:session_token" \
  POST http://localhost:7450/api/auth/email/resend-verification-code
