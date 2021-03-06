// Package cauth provides the primitives and the service layer for authentication. It supports multiple forms of
// authentication including username/password, email/password, phone otp, email magic links, etc. It includes a
// VerifySessionMiddleware that can be used in chttp.Route to protect it with a valid user session.
package cauth
