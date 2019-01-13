package dynauthconst

// Db returns a slice with the needed db information
func Db() []string {
	dbinfo := []string{"mysql", DatabaseUser + ":" + DatabasePass + "@/" + DatabaseName}
	return dbinfo
}
