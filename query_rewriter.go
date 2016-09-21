//  Copyright (c) 2016 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package bleve

import (
	"fmt"

	"golang.org/x/net/context"
)

const queryRewriterContextKey = "bleve.queryRewriter"

// WithQueryRewriter associates a query rewriter to a context.
func WithQueryRewriter(parent context.Context, queryRewriter QueryRewriter) context.Context {
	return context.WithValue(parent, queryRewriterContextKey, queryRewriter)
}

// The QueryRewriter interface allows applications to provide their
// own, optional query optimization logic.
type QueryRewriter interface {
	// The input variables should be treated as read-only / immutable.
	RewriteQuery(ctx context.Context, index Index, req *SearchRequest, query Query) (Query, error)
}

// A StandardQueryRewriter implements the QueryRewriter interface and
// performs only a basic set of standard query optimizations.
type StandardQueryRewriter struct{}

func (qr *StandardQueryRewriter) RewriteQuery(ctx context.Context, index Index, req *SearchRequest, query Query) (Query, error) {
	fmt.Printf("rewriting? %+v\n", query)

	bquery, ok := query.(*booleanQuery)
	if ok {
		if bquery.Must == nil && bquery.MustNot == nil && bquery.Should != nil {
			fmt.Printf("rewriting!!!\n")

			return qr.RewriteQuery(ctx, index, req, bquery.Should)
		}
	}

	return query, nil
}

var standardQueryRewriter = &StandardQueryRewriter{}
