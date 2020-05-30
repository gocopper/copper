#!/bin/zsh

http -v \
  POST http://localhost:7450/api/auth/email/login \
  email="user@email.com" password="password"
