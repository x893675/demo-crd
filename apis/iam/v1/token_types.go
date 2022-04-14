/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TokenStatus defines the observed state of Token
type TokenStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +optional
	ExpiresAt *metav1.Time `json:"expiresAt,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:JSONPath=".metadata.labels.iam\\.x893675\\.io/user",name="User",type=string
//+kubebuilder:printcolumn:name="IssueAt",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="ExpiresAt",type="date",JSONPath=".status.expiresAt"
//+kubebuilder:resource:categories="all",scope="Cluster",shortName="tk",singular="token"

// Token is the Schema for the tokens API
type Token struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	// access token
	AccessToken string `json:"accessToken"`

	// refresh token
	// +optional
	RefreshToken string `json:"refreshToken,omitempty"`

	// +kubebuilder:validation:Minimum=0
	// Expire time (second) for the token. 0 means no ttl
	// +optional
	ExpiresIn *int64      `json:"expiresIn,omitempty"`
	Status    TokenStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TokenList contains a list of Token
type TokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Token `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Token{}, &TokenList{})
}
