package cmailer

import "go.uber.org/fx"

var AWSFx = fx.Provide(
	NewAWSMailer,
)
