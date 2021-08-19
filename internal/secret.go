/*
 * Copyright 2021 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package internal

import "regexp"

var validSecretKey = regexp.MustCompile(`^[A-Za-z0-9\-_.]+$`)

// IsValidSecretKey tests whether a string is a valid Kubernetes Secret key:
//https://kubernetes.io/docs/concepts/configuration/secret/#overview-of-secrets.
func IsValidSecretKey(key string) bool {
	return validSecretKey.MatchString(key)
}
