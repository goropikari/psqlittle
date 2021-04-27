package core

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
