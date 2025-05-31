package ratelimit

import "go.uber.org/ratelimit"

var PlaymatchRatelimit ratelimit.Limiter = ratelimit.New(4)
var HasheousRatelimit ratelimit.Limiter = ratelimit.New(4)
var RommRatelimit ratelimit.Limiter = ratelimit.New(4)
