// Code generated by "stringer -type=Instruction"; DO NOT EDIT.

package vm

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PUSH0-0]
	_ = x[PUSHF-0]
	_ = x[PUSHBYTES1-1]
	_ = x[PUSHBYTES2-2]
	_ = x[PUSHBYTES3-3]
	_ = x[PUSHBYTES4-4]
	_ = x[PUSHBYTES5-5]
	_ = x[PUSHBYTES6-6]
	_ = x[PUSHBYTES7-7]
	_ = x[PUSHBYTES8-8]
	_ = x[PUSHBYTES9-9]
	_ = x[PUSHBYTES10-10]
	_ = x[PUSHBYTES11-11]
	_ = x[PUSHBYTES12-12]
	_ = x[PUSHBYTES13-13]
	_ = x[PUSHBYTES14-14]
	_ = x[PUSHBYTES15-15]
	_ = x[PUSHBYTES16-16]
	_ = x[PUSHBYTES17-17]
	_ = x[PUSHBYTES18-18]
	_ = x[PUSHBYTES19-19]
	_ = x[PUSHBYTES20-20]
	_ = x[PUSHBYTES21-21]
	_ = x[PUSHBYTES22-22]
	_ = x[PUSHBYTES23-23]
	_ = x[PUSHBYTES24-24]
	_ = x[PUSHBYTES25-25]
	_ = x[PUSHBYTES26-26]
	_ = x[PUSHBYTES27-27]
	_ = x[PUSHBYTES28-28]
	_ = x[PUSHBYTES29-29]
	_ = x[PUSHBYTES30-30]
	_ = x[PUSHBYTES31-31]
	_ = x[PUSHBYTES32-32]
	_ = x[PUSHBYTES33-33]
	_ = x[PUSHBYTES34-34]
	_ = x[PUSHBYTES35-35]
	_ = x[PUSHBYTES36-36]
	_ = x[PUSHBYTES37-37]
	_ = x[PUSHBYTES38-38]
	_ = x[PUSHBYTES39-39]
	_ = x[PUSHBYTES40-40]
	_ = x[PUSHBYTES41-41]
	_ = x[PUSHBYTES42-42]
	_ = x[PUSHBYTES43-43]
	_ = x[PUSHBYTES44-44]
	_ = x[PUSHBYTES45-45]
	_ = x[PUSHBYTES46-46]
	_ = x[PUSHBYTES47-47]
	_ = x[PUSHBYTES48-48]
	_ = x[PUSHBYTES49-49]
	_ = x[PUSHBYTES50-50]
	_ = x[PUSHBYTES51-51]
	_ = x[PUSHBYTES52-52]
	_ = x[PUSHBYTES53-53]
	_ = x[PUSHBYTES54-54]
	_ = x[PUSHBYTES55-55]
	_ = x[PUSHBYTES56-56]
	_ = x[PUSHBYTES57-57]
	_ = x[PUSHBYTES58-58]
	_ = x[PUSHBYTES59-59]
	_ = x[PUSHBYTES60-60]
	_ = x[PUSHBYTES61-61]
	_ = x[PUSHBYTES62-62]
	_ = x[PUSHBYTES63-63]
	_ = x[PUSHBYTES64-64]
	_ = x[PUSHBYTES65-65]
	_ = x[PUSHBYTES66-66]
	_ = x[PUSHBYTES67-67]
	_ = x[PUSHBYTES68-68]
	_ = x[PUSHBYTES69-69]
	_ = x[PUSHBYTES70-70]
	_ = x[PUSHBYTES71-71]
	_ = x[PUSHBYTES72-72]
	_ = x[PUSHBYTES73-73]
	_ = x[PUSHBYTES74-74]
	_ = x[PUSHBYTES75-75]
	_ = x[PUSHDATA1-76]
	_ = x[PUSHDATA2-77]
	_ = x[PUSHDATA4-78]
	_ = x[PUSHM1-79]
	_ = x[PUSH1-81]
	_ = x[PUSHT-81]
	_ = x[PUSH2-82]
	_ = x[PUSH3-83]
	_ = x[PUSH4-84]
	_ = x[PUSH5-85]
	_ = x[PUSH6-86]
	_ = x[PUSH7-87]
	_ = x[PUSH8-88]
	_ = x[PUSH9-89]
	_ = x[PUSH10-90]
	_ = x[PUSH11-91]
	_ = x[PUSH12-92]
	_ = x[PUSH13-93]
	_ = x[PUSH14-94]
	_ = x[PUSH15-95]
	_ = x[PUSH16-96]
	_ = x[NOP-97]
	_ = x[JMP-98]
	_ = x[JMPIF-99]
	_ = x[JMPIFNOT-100]
	_ = x[CALL-101]
	_ = x[RET-102]
	_ = x[APPCALL-103]
	_ = x[SYSCALL-104]
	_ = x[TAILCALL-105]
	_ = x[DUPFROMALTSTACK-106]
	_ = x[TOALTSTACK-107]
	_ = x[FROMALTSTACK-108]
	_ = x[XDROP-109]
	_ = x[XSWAP-114]
	_ = x[XTUCK-115]
	_ = x[DEPTH-116]
	_ = x[DROP-117]
	_ = x[DUP-118]
	_ = x[NIP-119]
	_ = x[OVER-120]
	_ = x[PICK-121]
	_ = x[ROLL-122]
	_ = x[ROT-123]
	_ = x[SWAP-124]
	_ = x[TUCK-125]
	_ = x[CAT-126]
	_ = x[SUBSTR-127]
	_ = x[LEFT-128]
	_ = x[RIGHT-129]
	_ = x[SIZE-130]
	_ = x[INVERT-131]
	_ = x[AND-132]
	_ = x[OR-133]
	_ = x[XOR-134]
	_ = x[EQUAL-135]
	_ = x[INC-139]
	_ = x[DEC-140]
	_ = x[SIGN-141]
	_ = x[NEGATE-143]
	_ = x[ABS-144]
	_ = x[NOT-145]
	_ = x[NZ-146]
	_ = x[ADD-147]
	_ = x[SUB-148]
	_ = x[MUL-149]
	_ = x[DIV-150]
	_ = x[MOD-151]
	_ = x[SHL-152]
	_ = x[SHR-153]
	_ = x[BOOLAND-154]
	_ = x[BOOLOR-155]
	_ = x[NUMEQUAL-156]
	_ = x[NUMNOTEQUAL-158]
	_ = x[LT-159]
	_ = x[GT-160]
	_ = x[LTE-161]
	_ = x[GTE-162]
	_ = x[MIN-163]
	_ = x[MAX-164]
	_ = x[WITHIN-165]
	_ = x[SHA1-167]
	_ = x[SHA256-168]
	_ = x[HASH160-169]
	_ = x[HASH256-170]
	_ = x[CHECKSIG-172]
	_ = x[VERIFY-173]
	_ = x[CHECKMULTISIG-174]
	_ = x[ARRAYSIZE-192]
	_ = x[PACK-193]
	_ = x[UNPACK-194]
	_ = x[PICKITEM-195]
	_ = x[SETITEM-196]
	_ = x[NEWARRAY-197]
	_ = x[NEWSTRUCT-198]
	_ = x[NEWMAP-199]
	_ = x[APPEND-200]
	_ = x[REVERSE-201]
	_ = x[REMOVE-202]
	_ = x[HASKEY-203]
	_ = x[KEYS-204]
	_ = x[VALUES-205]
	_ = x[THROW-240]
	_ = x[THROWIFNOT-241]
}

const (
	_Instruction_name_0 = "PUSH0PUSHBYTES1PUSHBYTES2PUSHBYTES3PUSHBYTES4PUSHBYTES5PUSHBYTES6PUSHBYTES7PUSHBYTES8PUSHBYTES9PUSHBYTES10PUSHBYTES11PUSHBYTES12PUSHBYTES13PUSHBYTES14PUSHBYTES15PUSHBYTES16PUSHBYTES17PUSHBYTES18PUSHBYTES19PUSHBYTES20PUSHBYTES21PUSHBYTES22PUSHBYTES23PUSHBYTES24PUSHBYTES25PUSHBYTES26PUSHBYTES27PUSHBYTES28PUSHBYTES29PUSHBYTES30PUSHBYTES31PUSHBYTES32PUSHBYTES33PUSHBYTES34PUSHBYTES35PUSHBYTES36PUSHBYTES37PUSHBYTES38PUSHBYTES39PUSHBYTES40PUSHBYTES41PUSHBYTES42PUSHBYTES43PUSHBYTES44PUSHBYTES45PUSHBYTES46PUSHBYTES47PUSHBYTES48PUSHBYTES49PUSHBYTES50PUSHBYTES51PUSHBYTES52PUSHBYTES53PUSHBYTES54PUSHBYTES55PUSHBYTES56PUSHBYTES57PUSHBYTES58PUSHBYTES59PUSHBYTES60PUSHBYTES61PUSHBYTES62PUSHBYTES63PUSHBYTES64PUSHBYTES65PUSHBYTES66PUSHBYTES67PUSHBYTES68PUSHBYTES69PUSHBYTES70PUSHBYTES71PUSHBYTES72PUSHBYTES73PUSHBYTES74PUSHBYTES75PUSHDATA1PUSHDATA2PUSHDATA4PUSHM1"
	_Instruction_name_1 = "PUSH1PUSH2PUSH3PUSH4PUSH5PUSH6PUSH7PUSH8PUSH9PUSH10PUSH11PUSH12PUSH13PUSH14PUSH15PUSH16NOPJMPJMPIFJMPIFNOTCALLRETAPPCALLSYSCALLTAILCALLDUPFROMALTSTACKTOALTSTACKFROMALTSTACKXDROP"
	_Instruction_name_2 = "XSWAPXTUCKDEPTHDROPDUPNIPOVERPICKROLLROTSWAPTUCKCATSUBSTRLEFTRIGHTSIZEINVERTANDORXOREQUAL"
	_Instruction_name_3 = "INCDECSIGN"
	_Instruction_name_4 = "NEGATEABSNOTNZADDSUBMULDIVMODSHLSHRBOOLANDBOOLORNUMEQUAL"
	_Instruction_name_5 = "NUMNOTEQUALLTGTLTEGTEMINMAXWITHIN"
	_Instruction_name_6 = "SHA1SHA256HASH160HASH256"
	_Instruction_name_7 = "CHECKSIGVERIFYCHECKMULTISIG"
	_Instruction_name_8 = "ARRAYSIZEPACKUNPACKPICKITEMSETITEMNEWARRAYNEWSTRUCTNEWMAPAPPENDREVERSEREMOVEHASKEYKEYSVALUES"
	_Instruction_name_9 = "THROWTHROWIFNOT"
)

var (
	_Instruction_index_0 = [...]uint16{0, 5, 15, 25, 35, 45, 55, 65, 75, 85, 95, 106, 117, 128, 139, 150, 161, 172, 183, 194, 205, 216, 227, 238, 249, 260, 271, 282, 293, 304, 315, 326, 337, 348, 359, 370, 381, 392, 403, 414, 425, 436, 447, 458, 469, 480, 491, 502, 513, 524, 535, 546, 557, 568, 579, 590, 601, 612, 623, 634, 645, 656, 667, 678, 689, 700, 711, 722, 733, 744, 755, 766, 777, 788, 799, 810, 821, 830, 839, 848, 854}
	_Instruction_index_1 = [...]uint8{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 51, 57, 63, 69, 75, 81, 87, 90, 93, 98, 106, 110, 113, 120, 127, 135, 150, 160, 172, 177}
	_Instruction_index_2 = [...]uint8{0, 5, 10, 15, 19, 22, 25, 29, 33, 37, 40, 44, 48, 51, 57, 61, 66, 70, 76, 79, 81, 84, 89}
	_Instruction_index_3 = [...]uint8{0, 3, 6, 10}
	_Instruction_index_4 = [...]uint8{0, 6, 9, 12, 14, 17, 20, 23, 26, 29, 32, 35, 42, 48, 56}
	_Instruction_index_5 = [...]uint8{0, 11, 13, 15, 18, 21, 24, 27, 33}
	_Instruction_index_6 = [...]uint8{0, 4, 10, 17, 24}
	_Instruction_index_7 = [...]uint8{0, 8, 14, 27}
	_Instruction_index_8 = [...]uint8{0, 9, 13, 19, 27, 34, 42, 51, 57, 63, 70, 76, 82, 86, 92}
	_Instruction_index_9 = [...]uint8{0, 5, 15}
)

func (i Instruction) String() string {
	switch {
	case 0 <= i && i <= 79:
		return _Instruction_name_0[_Instruction_index_0[i]:_Instruction_index_0[i+1]]
	case 81 <= i && i <= 109:
		i -= 81
		return _Instruction_name_1[_Instruction_index_1[i]:_Instruction_index_1[i+1]]
	case 114 <= i && i <= 135:
		i -= 114
		return _Instruction_name_2[_Instruction_index_2[i]:_Instruction_index_2[i+1]]
	case 139 <= i && i <= 141:
		i -= 139
		return _Instruction_name_3[_Instruction_index_3[i]:_Instruction_index_3[i+1]]
	case 143 <= i && i <= 156:
		i -= 143
		return _Instruction_name_4[_Instruction_index_4[i]:_Instruction_index_4[i+1]]
	case 158 <= i && i <= 165:
		i -= 158
		return _Instruction_name_5[_Instruction_index_5[i]:_Instruction_index_5[i+1]]
	case 167 <= i && i <= 170:
		i -= 167
		return _Instruction_name_6[_Instruction_index_6[i]:_Instruction_index_6[i+1]]
	case 172 <= i && i <= 174:
		i -= 172
		return _Instruction_name_7[_Instruction_index_7[i]:_Instruction_index_7[i+1]]
	case 192 <= i && i <= 205:
		i -= 192
		return _Instruction_name_8[_Instruction_index_8[i]:_Instruction_index_8[i+1]]
	case 240 <= i && i <= 241:
		i -= 240
		return _Instruction_name_9[_Instruction_index_9[i]:_Instruction_index_9[i+1]]
	default:
		return "Instruction(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
