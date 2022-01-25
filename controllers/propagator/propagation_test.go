// Copyright Contributors to the Open Cluster Management project

package propagator

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	appsv1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/placementrule/v1"

	policiesv1 "github.com/stolostron/governance-policy-propagator/api/v1"
)

func TestInitializeAttempts(t *testing.T) {
	tests := []struct {
		envVarValue string
		expected    int
	}{
		{"", attemptsDefault},
		{fmt.Sprint(attemptsDefault + 2), attemptsDefault + 2},
		{"0", attemptsDefault},
		{"-3", attemptsDefault},
	}

	for _, test := range tests {
		t.Run(
			fmt.Sprintf(`%s="%s"`, attemptsEnvName, test.envVarValue),
			func(t *testing.T) {
				defer func() {
					// Reset to the default values
					attempts = 0

					err := os.Unsetenv(attemptsEnvName)
					if err != nil {
						t.Fatalf("failed to unset the environment variable: %v", err)
					}
				}()

				err := os.Setenv(attemptsEnvName, test.envVarValue)
				if err != nil {
					t.Fatalf("failed to set the environment variable: %v", err)
				}

				var k8sInterface kubernetes.Interface
				Initialize(&rest.Config{}, &k8sInterface)

				if attempts != test.expected {
					t.Fatalf("Expected attempts=%d, got %d", test.expected, attempts)
				}
			},
		)
	}
}

func TestInitializeRequeueErrorDelay(t *testing.T) {
	tests := []struct {
		envVarValue string
		expected    int
	}{
		{"", requeueErrorDelayDefault},
		{fmt.Sprint(requeueErrorDelayDefault + 2), requeueErrorDelayDefault + 2},
		{"0", requeueErrorDelayDefault},
		{"-3", requeueErrorDelayDefault},
	}

	for _, test := range tests {
		t.Run(
			fmt.Sprintf(`%s="%s"`, requeueErrorDelayEnvName, test.envVarValue),
			func(t *testing.T) {
				defer func() {
					// Reset to the default values
					requeueErrorDelay = 0

					err := os.Unsetenv(requeueErrorDelayEnvName)
					if err != nil {
						t.Fatalf("failed to unset the environment variable: %v", err)
					}
				}()

				err := os.Setenv(requeueErrorDelayEnvName, test.envVarValue)
				if err != nil {
					t.Fatalf("failed to set the environment variable: %v", err)
				}
				var k8sInterface kubernetes.Interface
				Initialize(&rest.Config{}, &k8sInterface)

				if requeueErrorDelay != test.expected {
					t.Fatalf("Expected requeueErrorDelay=%d, got %d", test.expected, attempts)
				}
			},
		)
	}
}

func TestInitializeConcurrencyPerPolicyEnvName(t *testing.T) {
	tests := []struct {
		envVarValue string
		expected    int
	}{
		{"", concurrencyPerPolicyDefault},
		{fmt.Sprint(concurrencyPerPolicyDefault + 2), concurrencyPerPolicyDefault + 2},
		{"0", concurrencyPerPolicyDefault},
		{"-3", concurrencyPerPolicyDefault},
	}

	for _, test := range tests {
		t.Run(
			fmt.Sprintf(`%s="%s"`, concurrencyPerPolicyEnvName, test.envVarValue),
			func(t *testing.T) {
				defer func() {
					// Reset to the default values
					concurrencyPerPolicy = 0

					err := os.Unsetenv(concurrencyPerPolicyEnvName)
					if err != nil {
						t.Fatalf("failed to unset the environment variable: %v", err)
					}
				}()

				err := os.Setenv(concurrencyPerPolicyEnvName, test.envVarValue)
				if err != nil {
					t.Fatalf("failed to set the environment variable: %v", err)
				}
				var k8sInterface kubernetes.Interface
				Initialize(&rest.Config{}, &k8sInterface)

				if concurrencyPerPolicy != test.expected {
					t.Fatalf("Expected concurrencyPerPolicy=%d, got %d", test.expected, attempts)
				}
			},
		)
	}
}

func TestInitializeEncryptionEnabledEnvName(t *testing.T) {
	tests := []struct {
		envVarValue string
		expected    bool
	}{
		{"true", true},
		{"false", false},
		{"something else", encryptionEnabledEnvDefault},
	}

	for _, test := range tests {
		t.Run(
			fmt.Sprintf(`%s="%s"`, encryptionEnabledEnvName, test.envVarValue),
			func(t *testing.T) {
				defer func() {
					// Reset to the default values
					encryptionEnabled = true

					err := os.Unsetenv(encryptionEnabledEnvName)
					if err != nil {
						t.Fatalf("failed to unset the environment variable: %v", err)
					}
				}()

				err := os.Setenv(encryptionEnabledEnvName, test.envVarValue)
				if err != nil {
					t.Fatalf("failed to set the environment variable: %v", err)
				}
				var k8sInterface kubernetes.Interface
				Initialize(&rest.Config{}, &k8sInterface)

				if encryptionEnabled != test.expected {
					t.Fatalf("Expected encryptionEnabled=%v, got %v", test.expected, attempts)
				}
			},
		)
	}
}

// A mock implementation of the PolicyReconciler for the handleDecisionWrapper function.
type MockPolicyReconciler struct {
	Err error
}

func (r MockPolicyReconciler) handleDecision(
	instance *policiesv1.Policy, decision appsv1.PlacementDecision,
) error {
	return r.Err
}

func TestHandleDecisionWrapper(t *testing.T) {
	// Simulate running Initialize. This is required for the retry library to actually run
	// the handleDecision method.
	attempts = 1
	defer func() {
		// Reset to the default value
		attempts = 0
	}()

	tests := []struct {
		Error         error
		ExpectedError bool
	}{
		{nil, false},
		{errors.New("some error"), true},
	}

	for _, test := range tests {
		// Simulate three placement decisions for the policy.
		decisions := []appsv1.PlacementDecision{
			{ClusterName: "cluster1", ClusterNamespace: "cluster1"},
			{ClusterName: "cluster2", ClusterNamespace: "cluster2"},
			{ClusterName: "cluster3", ClusterNamespace: "cluster3"},
		}
		policy := policiesv1.Policy{
			ObjectMeta: metav1.ObjectMeta{Name: "gambling-age", Namespace: "laws"},
		}

		// Load up the decisionsChan channel with all the decisions so that handleDecisionWrapper
		// will call handleDecision with each.
		decisionsChan := make(chan appsv1.PlacementDecision, len(decisions))

		for _, decision := range decisions {
			decisionsChan <- decision
		}

		resultsChan := make(chan decisionResult, len(decisions))

		// Instantiate the mock PolicyReconciler to pass to handleDecisionWrapper.
		reconciler := MockPolicyReconciler{Err: test.Error}

		go func() {
			start := time.Now()
			// Wait until handleDecisionWrapper has completed its work. Then close
			// the channel so that handleDecisionWrapper returns. This times out
			// after five seconds.
			for len(resultsChan) != len(decisions) {
				if time.Since(start) > (time.Second * 5) {
					close(decisionsChan)
				}
			}
			close(decisionsChan)
		}()

		handleDecisionWrapper(reconciler, &policy, decisionsChan, resultsChan)

		// Expect a 1x1 mapping of results to decisions.
		if len(resultsChan) != len(decisions) {
			t.Fatalf(
				"Expected the results channel length of %d, got %d", len(decisions), len(resultsChan),
			)
		}

		// Ensure all the results from the channel are as expected.
		for i := 0; i < len(decisions); i++ {
			result := <-resultsChan
			if test.ExpectedError {
				if result.Err == nil {
					t.Fatal("Expected an error but didn't get one")
				} else if result.Err != test.Error { // nolint: errorlint
					t.Fatalf("Expected the error %v but got: %v", test.Error, result.Err)
				}
			} else if result.Err != nil {
				t.Fatalf("Didn't expect but got: %v", result.Err)
			}

			expectedIdentifier := fmt.Sprintf("cluster%d/cluster%d", i+1, i+1)
			if result.Identifier != expectedIdentifier {
				t.Fatalf("Expected the identifier %s, got %s", result.Identifier, expectedIdentifier)
			}
		}
		close(resultsChan)
	}
}
