package validate

func ID(s string) bool { return s != "" && len(s) <= 128 }

func Worker(s string) bool { return s != "" }
