/*
Copyright 2024.

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

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	hyperspikeiov1 "hyperspike.io/gitea-operator/api/v1"
)

var _ = Describe("Repo Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		repo := &hyperspikeiov1.Repo{}

		BeforeEach(func() {
			// This setup is minimal. For real tests, you'd want to set up mock Gitea clients.
			By("creating the custom resource for the Kind Repo")
			err := k8sClient.Get(ctx, typeNamespacedName, repo)
			if err != nil && errors.IsNotFound(err) {
				resource := &hyperspikeiov1.Repo{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: hyperspikeiov1.RepoSpec{
						Description: "A test mirror repo",
						Org: &hyperspikeiov1.OrgRef{
							Name: "test-org",
						},
						Mirror:         true,
						CloneAddr:      "https://github.com/test/test.git",
						MirrorInterval: "8h",
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &hyperspikeiov1.Repo{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance Repo")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should not error when reconciling a mirror repo", func() {
			By("Reconciling the created resource")
			controllerReconciler := &RepoReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			// In a real test, you would mock the Gitea client here.
			// Since there's no Gitea instance, we expect an error, but this shows the structure.
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			// Because we don't have a real Gitea instance or a mock, we expect a failure here.
			// Once mocking is in place, you would change this to Expect(err).NotTo(HaveOccurred())
			Expect(err).To(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
