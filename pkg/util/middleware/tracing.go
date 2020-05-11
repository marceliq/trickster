/*
 * Copyright 2018 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package middleware

import (
	"net/http"

	"github.com/tricksterproxy/trickster/pkg/proxy/request"
	"github.com/tricksterproxy/trickster/pkg/tracing"
	tspan "github.com/tricksterproxy/trickster/pkg/tracing/span"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
)

// Trace attaches a Tracer to an HTTP request
func Trace(tr *tracing.Tracer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r, span := tspan.PrepareRequest(r, tr)
		if span != nil {
			defer span.End()

			rsc := request.GetResources(r)
			if rsc != nil &&
				rsc.OriginConfig != nil &&
				rsc.PathConfig != nil &&
				rsc.CacheConfig != nil {
				tspan.SetAttributes(tr, span,
					[]core.KeyValue{
						key.String("origin.name", rsc.OriginConfig.Name),
						key.String("origin.type", rsc.OriginConfig.OriginType),
						key.String("router.path", rsc.PathConfig.Path),
						key.String("cache.name", rsc.CacheConfig.Name),
						key.String("cache.type", rsc.CacheConfig.CacheType),
					}...,
				)
			}

		}
		next.ServeHTTP(w, r)
	})
}