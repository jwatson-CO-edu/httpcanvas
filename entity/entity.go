/*
entity.go
James Watson, 2021-01
Entity-Component system for basic 2D animations
*/

/*** DEV PLAN ***
[ ] Basic Entity-Component System with the following elements
	[ ] Basic Entity Type
		[ ] arbitrary array of component pointers
	[ ] Component Interface that enforces minimal structure to support system
	[ ] Special Data component that holds data common to all components
	[ ] Registry that maps componenents to handlers
		[ ] Offers a way to iterate over each type of component, in turn
		[ ] Baseline attempt at spatial and temporal locality with component arrays
*/

package entity

type Entity struct{
	// Base Entity Struct	
}

type Registry struct {
	// Mapping between `Component`s and their handlers
}

type Component interface{
	//
	GetDataRef() *Component
}

type BaseComponent struct{
	// Base Component Struct
	Data *Component
}