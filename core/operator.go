package core

import "fmt"

// Not negates x
func Not(x Value) Value {
	if x == Null {
		return Null
	}
	if x == True {
		return False
	}
	if x == False {
		return True
	}

	// Shouldn't reach this
	fmt.Println("should n't reach this xxxxxxxxxxxxxxxxxxxxxxxxxx")
	return Null
}

// OR calculates x OR y
func OR(x, y Value) Value {
	// memo:
	// True or Null -> True
	// Null or True -> True
	// False or Null -> Null
	// Null or False -> Null
	if x == True || y == True {
		return True
	}
	if x == Null || y == Null {
		return Null
	}
	if x == False && y == False {
		return False
	}
	return True
}

// AND calculates x AND y
func AND(x, y Value) Value {
	// memo:
	// True and Null -> Null
	// Null and True -> Null
	// False and Null -> False
	// Null and False -> False
	if x == True && y == True {
		return True
	}
	if x == False || y == False {
		return False
	}
	return Null
}

func LessForSort(x, y Value, sortDir int) bool {
	if y == Null {
		return true
	}
	if x == Null {
		return false
	}

	if sortDir < 3 {
		return x.(int) < y.(int)
	}
	return x.(int) > y.(int)
}
