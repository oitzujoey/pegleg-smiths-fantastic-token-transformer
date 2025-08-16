registerMachine(
	{
		machine = "microcomp",
		word_bit_width = 8,
		address_spaces = {
			{name = "program", address_bit_width = 16},
			{name = "data", address_bit_width = 16}
		}
	}
)
registerMachine(
	{
		machine = "ATmega328",
		word_bit_width = 8,
		address_spaces = {
			{name = "program", address_bit_width = 16},
			{name = "data", address_bit_width = 16},
			{name = "EEPROM", address_bit_width = 8}
		}
	}
)
registerMachine(
	{
		machine = "x86_64",
		word_bit_width = 64,
		address_spaces = {{name = nil, address_bit_width = 64}}
	}
)
