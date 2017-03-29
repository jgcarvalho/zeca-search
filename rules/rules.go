package rules

type Config struct {
	Input  string `toml:"input"`
	Output string `toml:"output"`
}

// type Pattern [3]string
// type Rule map[Pattern]string

type Rule [NumStates][NumStates][NumStates]State

type State uint8

const NumStates = 69

const (
	S__ State = iota
	S_A       //1
	S_C
	S_D
	S_E
	S_F
	S_G
	S_H
	S_I
	S_K
	S_L
	S_M
	S_N
	S_P
	S_Q
	S_R
	S_S
	S_T
	S_V
	S_Y
	S_W //20
	S_c
	S_h
	S_e
	// rose hydrophobicity
	S_An //24
	S_Cn
	S_Fn
	S_In
	S_Ln
	S_Mn
	S_Vn
	S_Wn
	S_cn //32
	S_hn
	S_en //34
	// polar
	S_Dp
	S_Ep
	S_Gp //37
	S_Hp
	S_Kp
	S_Np
	S_Pp //41
	S_Qp
	S_Rp
	S_Sp
	S_Tp
	S_Yp
	S_cp //47
	S_hp
	S_ep //49
	// special residues
	// gly
	S_GG //50
	S_cG //51
	S_hG
	S_eG //53
	// pro
	S_PP //54
	S_cP
	S_hP
	S_eP //57
	// charged residues
	// negative
	S_Dneg
	S_Eneg
	S_cneg
	S_hneg
	S_eneg
	// positive
	S_Kpos
	S_Rpos
	S_cpos
	S_hpos
	S_epos
	// S_init return to init state
	S_init
	// total of 69 states
)

// doing transition
func Transition(c State, st State) State {
	var s State
	if st == 3 {
		s = S_init
	} else {
		if c == S__ {
			s = S__
		} else if (c >= S_A) && (c <= S_e) {
			s = S_c + st
		} else if (c >= S_An) && (c <= S_en) {
			s = S_cn + st
		} else if (c >= S_Dp) && (c <= S_ep) {
			s = S_cp + st
		} else if (c >= S_GG) && (c <= S_eG) {
			s = S_cG + st
		} else if (c >= S_PP) && (c <= S_eP) {
			s = S_cP + st
		} else if (c >= S_Dneg) && (c <= S_eneg) {
			s = S_cneg + st
		} else if (c >= S_Kpos) && (c <= S_epos) {
			s = S_cpos + st
		} else {
			panic("Transition state not defined")
		}
	}
	return s
}

func SS(c State) string {
	switch c {
	case S_h, S_hn, S_hp, S_hG, S_hP, S_hpos, S_hneg:
		return "helix"
	case S_e, S_en, S_ep, S_eG, S_eP, S_epos, S_eneg:
		return "strand"
	case S_c, S_cn, S_cp, S_cG, S_cP, S_cpos, S_cneg:
		return "coil"
	default:
		return "none"
	}
}

func String2State(st string) State {
	var s State
	switch st {
	case "#":
		s = S__
	case "A":
		s = S_A
	case "C":
		s = S_C
	case "D":
		s = S_D
	case "E":
		s = S_E
	case "F":
		s = S_F
	case "G":
		s = S_G
	case "H":
		s = S_H
	case "I":
		s = S_I
	case "K":
		s = S_K
	case "L":
		s = S_L
	case "M":
		s = S_M
	case "N":
		s = S_N
	case "P":
		s = S_P
	case "Q":
		s = S_Q
	case "R":
		s = S_R
	case "S":
		s = S_S
	case "T":
		s = S_T
	case "V":
		s = S_V
	case "Y":
		s = S_Y
	case "W":
		s = S_W
	case "_":
		s = S_c
	case "*":
		s = S_h
	case "|":
		s = S_e
	case "An":
		s = S_An
	case "Cn":
		s = S_Cn
	case "Dp":
		s = S_Dp
	case "Fn":
		s = S_Fn
	case "Ep":
		s = S_Ep
	case "Gp":
		s = S_Gp
	case "Hp":
		s = S_Hp
	case "In":
		s = S_In
	case "Kp":
		s = S_Kp
	case "Ln":
		s = S_Ln
	case "Mn":
		s = S_Mn
	case "Np":
		s = S_Np
	case "Pp":
		s = S_Pp
	case "Qp":
		s = S_Qp
	case "Rp":
		s = S_Rp
	case "Sp":
		s = S_Sp
	case "Tp":
		s = S_Tp
	case "Vn":
		s = S_Vn
	case "Yp":
		s = S_Yp
	case "Wn":
		s = S_Wn
	case "_n":
		s = S_cn
	case "*n":
		s = S_hn
	case "|n":
		s = S_en
	case "_p":
		s = S_cp
	case "*p":
		s = S_hp
	case "|p":
		s = S_ep
	case "GG":
		s = S_GG
	case "PP":
		s = S_PP
	case "_G":
		s = S_cG
	case "*G":
		s = S_hG
	case "|G":
		s = S_eG
	case "_P":
		s = S_cP
	case "*P":
		s = S_hP
	case "|P":
		s = S_eP
	case "D-":
		s = S_Dneg
	case "E-":
		s = S_Eneg
	case "K+":
		s = S_Kpos
	case "R+":
		s = S_Rpos
	case "_-":
		s = S_cneg
	case "*-":
		s = S_hneg
	case "|-":
		s = S_eneg
	case "_+":
		s = S_cpos
	case "*+":
		s = S_hpos
	case "|+":
		s = S_epos
	case "?", "??", "?n", "?p", "?G", "?P", "?-", "?+":
		s = S_init
	}
	return s
}

func State2String(s State) string {
	var state string
	switch s {
	case S__:
		state = "#"
	case S_A:
		state = "A"
	case S_C:
		state = "C"
	case S_D:
		state = "D"
	case S_E:
		state = "E"
	case S_F:
		state = "F"
	case S_G:
		state = "G"
	case S_H:
		state = "H"
	case S_I:
		state = "I"
	case S_K:
		state = "K"
	case S_L:
		state = "L"
	case S_M:
		state = "M"
	case S_N:
		state = "N"
	case S_P:
		state = "P"
	case S_Q:
		state = "Q"
	case S_R:
		state = "R"
	case S_S:
		state = "S"
	case S_T:
		state = "T"
	case S_V:
		state = "V"
	case S_Y:
		state = "Y"
	case S_W:
		state = "W"
	case S_c:
		state = "_"
	case S_h:
		state = "*"
	case S_e:
		state = "|"
	// rose hydrophobicity:
	case S_An:
		state = "An"
	case S_Cn:
		state = "Cn"
	case S_Fn:
		state = "Fn"
	case S_In:
		state = "In"
	case S_Ln:
		state = "Ln"
	case S_Mn:
		state = "Mn"
	case S_Vn:
		state = "Vn"
	case S_Wn:
		state = "Wn"
	case S_cn:
		state = "_n"
	case S_hn:
		state = "*n"
	case S_en:
		state = "|n"
	// polar:
	case S_Dp:
		state = "Dp"
	case S_Ep:
		state = "Ep"
	case S_Gp:
		state = "Gp"
	case S_Hp:
		state = "Hp"
	case S_Kp:
		state = "Kp"
	case S_Np:
		state = "Np"
	case S_Pp:
		state = "Pp"
	case S_Qp:
		state = "Qp"
	case S_Rp:
		state = "Rp"
	case S_Sp:
		state = "Sp"
	case S_Tp:
		state = "Tp"
	case S_Yp:
		state = "Yp"
	case S_cp:
		state = "_p"
	case S_hp:
		state = "*p"
	case S_ep:
		state = "|p"
	// special residues:
	// gly:
	case S_GG:
		state = "GG"
	case S_cG:
		state = "_G"
	case S_hG:
		state = "*G"
	case S_eG:
		state = "|G"
	// pro:
	case S_PP:
		state = "PP"
	case S_cP:
		state = "_P"
	case S_hP:
		state = "*P"
	case S_eP:
		state = "|P"
	// charged residues:
	// negative:
	case S_Dneg:
		state = "D-"
	case S_Eneg:
		state = "E-"
	case S_cneg:
		state = "_-"
	case S_hneg:
		state = "*-"
	case S_eneg:
		state = "|-"
	// positive
	case S_Kpos:
		state = "K+"
	case S_Rpos:
		state = "R+"
	case S_cpos:
		state = "_+"
	case S_hpos:
		state = "*+"
	case S_epos:
		state = "|+"
		// S_init state =  to init state
	case S_init:
		state = "?"
	}
	return state
}
