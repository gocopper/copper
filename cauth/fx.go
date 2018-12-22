// Package cauth provides tools for authentication and authorization. It has end-to-end API flows for user registration,
// and password management. Authorization is provided with roles and permissions.
package cauth

import (
	"go.uber.org/fx"
)

// Fx module for the cauth package that provides the SQL implementation for all services.
// Additionally, it registers the routes for authentication flows.
var Fx = fx.Provide(
	newSQLUserRepo,
	newUsersSvc,

	newRouter,
	newSignupRoute,
	newLoginRoute,
	newVerifyUserRoute,
	newResendVerificationCodeRoute,
	newResetPasswordRoute,

	newAuthMiddleware,
)

// RunMigrations can be used with fx.Invoke to run the db migrations for SQL implementation of cauth services.
var RunMigrations = fx.Invoke(
	runMigrations,
)
