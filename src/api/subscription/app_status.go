// Package subscription defines pub/sub primitives for the GraphQL API. This
// file provides helpers for tracking application status and resides in
// src/api/subscription.
package subscription

// ApplicationStatus reflects the current overall status of the application.
var ApplicationStatus AppStatus = AppStatusUp
