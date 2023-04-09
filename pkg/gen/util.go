package gen

// DataType2Len is to convert data type into length
func DataType2Len(t string) int {
	switch t {
	case "int":
		return 16
	case "bigint":
		return 64
	case "varchar":
		return 511
	case "timestamp":
		return 255
	case "datetime":
		return 255
	case "text":
		return 511
	case "float":
		return 53
	}
	return 16
}
