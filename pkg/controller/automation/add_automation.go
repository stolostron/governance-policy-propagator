// Copyright Contributors to the Open Cluster Management project

package automation

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, Add)
}
