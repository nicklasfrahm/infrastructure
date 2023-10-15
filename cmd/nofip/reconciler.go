package main

import (
	"sync"
)

// ReconciliationState defines the state of a resource
// before the reconciliation process.
type ReconciliationState struct {
	// Current is the current state of the resource.
	Current map[string]bool
	// Desired is the desired state of the resource.
	Desired map[string]bool
}

// ReconciliationPlan defines the plan to reconcile a
// resource.
type ReconciliationPlan struct {
	// Create is the list of records to create.
	Create map[string]bool
	// Delete is the list of records to delete.
	Delete map[string]bool
}

// Reconciler is stores the current and desired state
// of a resource along with a plan to reconcile them.
type Reconciler struct {
	sync.Mutex
	// State is the state of the resource before the reconciliation.
	State *ReconciliationState
	// Plan is the plan to reconcile the resource.
	Plan *ReconciliationPlan
}

// ShouldUpdate updates the plan to reconcile the resource.
func (r *Reconciler) ShouldUpdate() bool {
	shouldUpdate := false

	// Check if any new records should be created in the new set.
	for key := range r.State.Desired {
		if !r.State.Current[key] {
			shouldUpdate = true
			r.Plan.Create[key] = true
		}
	}

	// Check if any old records should be removed in the new set.
	for key := range r.State.Current {
		if !r.State.Desired[key] {
			shouldUpdate = true
			r.Plan.Delete[key] = true
		}
	}

	return shouldUpdate
}

// NewReconciler initializes a new Reconciler and returns a pointer to it.
func NewReconciler() *Reconciler {
	return &Reconciler{
		State: &ReconciliationState{
			Current: map[string]bool{},
			Desired: map[string]bool{},
		},
		Plan: &ReconciliationPlan{
			Create: map[string]bool{},
			Delete: map[string]bool{},
		},
	}
}
