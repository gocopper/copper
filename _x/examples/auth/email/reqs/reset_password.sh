#!/bin/zsh

http -v \
  POST http://localhost:7450/api/auth/email/reset-password \
  email="user@email.com"
