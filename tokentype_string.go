// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package parser

import "fmt"

const _TokenType_name = "ILLEGALEOFWSINTIDENTASTERISKCOMMAPERIODLPARENRPARENADDSUBEQUALLSSGTRQUESTIONSELECTINSERTUPDATEDELETECREATEALTERDROPFROMASSETINTOWHEREJOINLEFTRIGHTFULLOUTERINNERONGROUPBYORDERBYHAVINGONDUPLICATEKEYUPDATEDESCASCNULLPRIMARYKEYANDORIFNOTEXISTCOLUMNDEFAULTDATABASETABLEASSERTIONINDEXCHECKREFERENCEUNIQUEINTEGERSERIALVARCHAR"

var _TokenType_index = [...]uint16{0, 7, 10, 12, 15, 20, 28, 33, 39, 45, 51, 54, 57, 62, 65, 68, 76, 82, 88, 94, 100, 106, 111, 115, 119, 121, 124, 128, 133, 137, 141, 146, 150, 155, 160, 162, 169, 176, 182, 202, 206, 209, 213, 223, 226, 228, 230, 233, 238, 244, 251, 259, 264, 273, 278, 283, 292, 298, 305, 311, 318}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return fmt.Sprintf("TokenType(%d)", i)
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
