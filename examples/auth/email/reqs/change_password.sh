#!/bin/zsh

http -v \
  POST http://localhost:7450/api/auth/email/change-password \
  email="user@email.com" old_password="old_password" new_password="new_password"
