#!/bin/zsh

http -v \
  -a "user_uuid:session_token" \
  POST http://localhost:7450/api/auth/email/verify \
  verification_code="abc123"
