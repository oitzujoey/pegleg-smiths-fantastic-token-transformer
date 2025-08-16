function load_input(path)
	bytes = file_read_bytes(path)
	return ArrayLiteral(bytes, Integer(8))
end
