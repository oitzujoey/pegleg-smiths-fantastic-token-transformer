package main

import ("fmt"; "io/fs"; "regexp")

// Parameters: A list of files
// Dispatch a goroutine for each file.
//   Read the file.
//   Begin parsing.
//   Message back AST.
// Merge ASTs?

type parser_result struct {
	bytes []byte
	error error
}

type parser_configuration struct {
	// operator_tokens []string
	token_regexp *regexp.Regexp
	lookahead int
}

type parse_error struct {
	message string
}

func (e *parse_error) Error() string {
	return e.message
}


func parse(buildConfiguration *buildConfiguration) {
	file_names := buildConfiguration.files
	var parser_configuration parser_configuration
	// parser_configuration.operator_tokens = ;
	// Prepare operator regexp.
	{
		operator_tokens := []string{"//", "/*", "*/", "~=", "!=", "%=", "^=", "&=", "*=", "+=", "|=", ":=", "<=", ">=", "/=", "~", "!", "%", "^", "&", "*", "(", ")", "-", "+", "=", "{", "}", "|", ":", "\"", "<", ">", "[", "]", "\\", ";", "'", ",", ".", "/", "rn", "n", "r", "`"}
		{
			operator_token_regexp_string := ""
			for index, operator := range operator_tokens {
				operator_length := len(operator)
				if index != 0 {operator_token_regexp_string += "|"}
				if operator_length > 1 {operator_token_regexp_string += "("}
				for _, character := range operator {
					operator_token_regexp_string += "\\" + string(character)
				}
				if operator_length > 1 {operator_token_regexp_string += ")"}
			}
			token_regexp_string := "^(([a-zA-Z_][0-9a-zA-Z_]*)|([0-9_]+)|([\t ]+)|" + operator_token_regexp_string + ")"
			fmt.Println(token_regexp_string)
			token_regexp, e := regexp.Compile(token_regexp_string)
			if e != nil {
				exit(1, "Failed to compile tokenization regexp.")
			}
			parser_configuration.token_regexp = token_regexp
		}
		{
			max_length := 0
			for _, token := range operator_tokens {
				token_length := len(token)
				if max_length < token_length {max_length = token_length}
			}
			parser_configuration.lookahead = max_length
		}
	}
	// Run parsers.
	parser_channels := make([]chan parser_result, len(file_names))
	for files_index, file_name := range file_names {
		parser_channel := make(chan parser_result)
		go goroutine_parse(parser_channel, &parser_configuration, buildConfiguration, file_name)
		parser_channels[files_index] = parser_channel
	}
	for _, parser_channel := range parser_channels {
		parser_result := <-parser_channel
		fmt.Printf("%v\n", parser_result.error)
	}
}

func goroutine_parse(parser_result_channel chan parser_result, parser_configuration *parser_configuration, buildConfiguration *buildConfiguration, file_name string) {
	var e error = nil

	fmt.Printf("Parsing \"%v\".\n", file_name)

	file_bytes, e := fs.ReadFile(buildConfiguration.fileSystem, file_name)
	if e != nil {
		exit(1, "Could not read file.")
	}

	tokens, e := tokenize(parser_configuration, file_bytes, file_name)
	if e != nil {
		exit(1, "Tokenization failed.")
	}
	{
		for _, token := range tokens {
			token_string := string(token)
			if token_string == "\n" || token_string == "\r\n" || token_string == "\r" {
				fmt.Println("\\n")
			} else if token_string == "\\" {
				fmt.Println("\\\\")
			} else if token_string == " " {
				fmt.Println("\\ ")
			} else {
				fmt.Printf("%v\n", token_string)
			}
		}
	}

	var parser_result parser_result
	parser_result.bytes = file_bytes
	parser_result.error = e
	parser_result_channel <- parser_result
}

func tokenize(parser_configuration *parser_configuration, source []byte, file_name string) ([][]byte, error) {
	// Maximum munch.
	var tokens [][]byte
	for len(source) > 0 {
		token := parser_configuration.token_regexp.Find(source)
		if token == nil {
			return nil, &parse_error{fmt.Sprintf("Failed to tokenize \"%v\".", file_name)}
		}
		tokens = append(tokens, token)
		source = source[len(token):]
	}
	return tokens, nil
}
