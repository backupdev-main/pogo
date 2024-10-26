// Code generated by gocc; DO NOT EDIT.

package lexer

/*
Let s be the current state
Let r be the current input rune
transitionTable[s](r) returns the next state.
*/
type TransitionTable [NumStates]func(rune) int

var TransTab = TransitionTable{
	// S0
	func(r rune) int {
		switch {
		case r == 9: // ['\t','\t']
			return 1
		case r == 10: // ['\n','\n']
			return 1
		case r == 13: // ['\r','\r']
			return 1
		case r == 32: // [' ',' ']
			return 1
		case r == 33: // ['!','!']
			return 2
		case r == 34: // ['"','"']
			return 3
		case r == 40: // ['(','(']
			return 4
		case r == 41: // [')',')']
			return 5
		case r == 42: // ['*','*']
			return 6
		case r == 43: // ['+','+']
			return 7
		case r == 44: // [',',',']
			return 8
		case r == 45: // ['-','-']
			return 7
		case r == 47: // ['/','/']
			return 9
		case r == 48: // ['0','0']
			return 10
		case 49 <= r && r <= 57: // ['1','9']
			return 11
		case r == 58: // [':',':']
			return 12
		case r == 59: // [';',';']
			return 13
		case r == 60: // ['<','<']
			return 14
		case r == 61: // ['=','=']
			return 15
		case r == 62: // ['>','>']
			return 14
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case r == 96: // ['`','`']
			return 18
		case r == 97: // ['a','a']
			return 16
		case r == 98: // ['b','b']
			return 19
		case 99 <= r && r <= 100: // ['c','d']
			return 16
		case r == 101: // ['e','e']
			return 20
		case r == 102: // ['f','f']
			return 21
		case 103 <= r && r <= 104: // ['g','h']
			return 16
		case r == 105: // ['i','i']
			return 22
		case 106 <= r && r <= 111: // ['j','o']
			return 16
		case r == 112: // ['p','p']
			return 23
		case 113 <= r && r <= 117: // ['q','u']
			return 16
		case r == 118: // ['v','v']
			return 24
		case r == 119: // ['w','w']
			return 25
		case 120 <= r && r <= 122: // ['x','z']
			return 16
		case r == 123: // ['{','{']
			return 26
		case r == 125: // ['}','}']
			return 27
		}
		return NoState
	},
	// S1
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S2
	func(r rune) int {
		switch {
		case r == 61: // ['=','=']
			return 14
		}
		return NoState
	},
	// S3
	func(r rune) int {
		switch {
		case r == 34: // ['"','"']
			return 28
		case r == 92: // ['\','\']
			return 29
		default:
			return 3
		}
	},
	// S4
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S5
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S6
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S7
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S8
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S9
	func(r rune) int {
		switch {
		case r == 42: // ['*','*']
			return 30
		case r == 47: // ['/','/']
			return 31
		}
		return NoState
	},
	// S10
	func(r rune) int {
		switch {
		case r == 46: // ['.','.']
			return 32
		}
		return NoState
	},
	// S11
	func(r rune) int {
		switch {
		case r == 46: // ['.','.']
			return 32
		case 48 <= r && r <= 57: // ['0','9']
			return 33
		}
		return NoState
	},
	// S12
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S13
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S14
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S15
	func(r rune) int {
		switch {
		case r == 61: // ['=','=']
			return 14
		}
		return NoState
	},
	// S16
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
	// S17
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
	// S18
	func(r rune) int {
		switch {
		case r == 96: // ['`','`']
			return 35
		default:
			return 18
		}
	},
	// S19
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 100: // ['a','d']
			return 16
		case r == 101: // ['e','e']
			return 36
		case 102 <= r && r <= 122: // ['f','z']
			return 16
		}
		return NoState
	},
	// S20
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 107: // ['a','k']
			return 16
		case r == 108: // ['l','l']
			return 37
		case r == 109: // ['m','m']
			return 16
		case r == 110: // ['n','n']
			return 38
		case 111 <= r && r <= 122: // ['o','z']
			return 16
		}
		return NoState
	},
	// S21
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 107: // ['a','k']
			return 16
		case r == 108: // ['l','l']
			return 39
		case 109 <= r && r <= 116: // ['m','t']
			return 16
		case r == 117: // ['u','u']
			return 40
		case 118 <= r && r <= 122: // ['v','z']
			return 16
		}
		return NoState
	},
	// S22
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 101: // ['a','e']
			return 16
		case r == 102: // ['f','f']
			return 41
		case 103 <= r && r <= 109: // ['g','m']
			return 16
		case r == 110: // ['n','n']
			return 42
		case 111 <= r && r <= 122: // ['o','z']
			return 16
		}
		return NoState
	},
	// S23
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 113: // ['a','q']
			return 16
		case r == 114: // ['r','r']
			return 43
		case 115 <= r && r <= 122: // ['s','z']
			return 16
		}
		return NoState
	},
	// S24
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case r == 97: // ['a','a']
			return 44
		case 98 <= r && r <= 122: // ['b','z']
			return 16
		}
		return NoState
	},
	// S25
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 103: // ['a','g']
			return 16
		case r == 104: // ['h','h']
			return 45
		case 105 <= r && r <= 122: // ['i','z']
			return 16
		}
		return NoState
	},
	// S26
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S27
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S28
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S29
	func(r rune) int {
		switch {
		case r == 34: // ['"','"']
			return 3
		case r == 110: // ['n','n']
			return 46
		case r == 114: // ['r','r']
			return 46
		case r == 116: // ['t','t']
			return 46
		}
		return NoState
	},
	// S30
	func(r rune) int {
		switch {
		case r == 42: // ['*','*']
			return 47
		default:
			return 30
		}
	},
	// S31
	func(r rune) int {
		switch {
		case r == 10: // ['\n','\n']
			return 48
		default:
			return 31
		}
	},
	// S32
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 49
		}
		return NoState
	},
	// S33
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 33
		}
		return NoState
	},
	// S34
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
	// S35
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S36
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 102: // ['a','f']
			return 16
		case r == 103: // ['g','g']
			return 50
		case 104 <= r && r <= 122: // ['h','z']
			return 16
		}
		return NoState
	},
	// S37
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 114: // ['a','r']
			return 16
		case r == 115: // ['s','s']
			return 51
		case 116 <= r && r <= 122: // ['t','z']
			return 16
		}
		return NoState
	},
	// S38
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 99: // ['a','c']
			return 16
		case r == 100: // ['d','d']
			return 52
		case 101 <= r && r <= 122: // ['e','z']
			return 16
		}
		return NoState
	},
	// S39
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 110: // ['a','n']
			return 16
		case r == 111: // ['o','o']
			return 53
		case 112 <= r && r <= 122: // ['p','z']
			return 16
		}
		return NoState
	},
	// S40
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 109: // ['a','m']
			return 16
		case r == 110: // ['n','n']
			return 54
		case 111 <= r && r <= 122: // ['o','z']
			return 16
		}
		return NoState
	},
	// S41
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
	// S42
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 115: // ['a','s']
			return 16
		case r == 116: // ['t','t']
			return 55
		case 117 <= r && r <= 122: // ['u','z']
			return 16
		}
		return NoState
	},
	// S43
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 104: // ['a','h']
			return 16
		case r == 105: // ['i','i']
			return 56
		case 106 <= r && r <= 110: // ['j','n']
			return 16
		case r == 111: // ['o','o']
			return 57
		case 112 <= r && r <= 122: // ['p','z']
			return 16
		}
		return NoState
	},
	// S44
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 113: // ['a','q']
			return 16
		case r == 114: // ['r','r']
			return 58
		case 115 <= r && r <= 122: // ['s','z']
			return 16
		}
		return NoState
	},
	// S45
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 104: // ['a','h']
			return 16
		case r == 105: // ['i','i']
			return 59
		case 106 <= r && r <= 122: // ['j','z']
			return 16
		}
		return NoState
	},
	// S46
	func(r rune) int {
		switch {
		case r == 34: // ['"','"']
			return 28
		case r == 92: // ['\','\']
			return 29
		default:
			return 3
		}
	},
	// S47
	func(r rune) int {
		switch {
		case r == 42: // ['*','*']
			return 47
		case r == 47: // ['/','/']
			return 60
		default:
			return 30
		}
	},
	// S48
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S49
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 49
		}
		return NoState
	},
	// S50
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 104: // ['a','h']
			return 16
		case r == 105: // ['i','i']
			return 61
		case 106 <= r && r <= 122: // ['j','z']
			return 16
		}
		return NoState
	},
	// S51
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 100: // ['a','d']
			return 16
		case r == 101: // ['e','e']
			return 62
		case 102 <= r && r <= 122: // ['f','z']
			return 16
		}
		return NoState
	},
	// S52
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
	// S53
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case r == 97: // ['a','a']
			return 63
		case 98 <= r && r <= 122: // ['b','z']
			return 16
		}
		return NoState
	},
	// S54
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 98: // ['a','b']
			return 16
		case r == 99: // ['c','c']
			return 64
		case 100 <= r && r <= 122: // ['d','z']
			return 16
		}
		return NoState
	},
	// S55
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
	// S56
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 109: // ['a','m']
			return 16
		case r == 110: // ['n','n']
			return 65
		case 111 <= r && r <= 122: // ['o','z']
			return 16
		}
		return NoState
	},
	// S57
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 102: // ['a','f']
			return 16
		case r == 103: // ['g','g']
			return 66
		case 104 <= r && r <= 122: // ['h','z']
			return 16
		}
		return NoState
	},
	// S58
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
	// S59
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 107: // ['a','k']
			return 16
		case r == 108: // ['l','l']
			return 67
		case 109 <= r && r <= 122: // ['m','z']
			return 16
		}
		return NoState
	},
	// S60
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S61
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 109: // ['a','m']
			return 16
		case r == 110: // ['n','n']
			return 68
		case 111 <= r && r <= 122: // ['o','z']
			return 16
		}
		return NoState
	},
	// S62
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
	// S63
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 115: // ['a','s']
			return 16
		case r == 116: // ['t','t']
			return 55
		case 117 <= r && r <= 122: // ['u','z']
			return 16
		}
		return NoState
	},
	// S64
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
	// S65
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 115: // ['a','s']
			return 16
		case r == 116: // ['t','t']
			return 69
		case 117 <= r && r <= 122: // ['u','z']
			return 16
		}
		return NoState
	},
	// S66
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 113: // ['a','q']
			return 16
		case r == 114: // ['r','r']
			return 70
		case 115 <= r && r <= 122: // ['s','z']
			return 16
		}
		return NoState
	},
	// S67
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 100: // ['a','d']
			return 16
		case r == 101: // ['e','e']
			return 71
		case 102 <= r && r <= 122: // ['f','z']
			return 16
		}
		return NoState
	},
	// S68
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
	// S69
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
	// S70
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case r == 97: // ['a','a']
			return 72
		case 98 <= r && r <= 122: // ['b','z']
			return 16
		}
		return NoState
	},
	// S71
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
	// S72
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 108: // ['a','l']
			return 16
		case r == 109: // ['m','m']
			return 73
		case 110 <= r && r <= 122: // ['n','z']
			return 16
		}
		return NoState
	},
	// S73
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 34
		case 65 <= r && r <= 90: // ['A','Z']
			return 16
		case r == 95: // ['_','_']
			return 17
		case 97 <= r && r <= 122: // ['a','z']
			return 16
		}
		return NoState
	},
}
