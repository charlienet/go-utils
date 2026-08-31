/*
Package stringx provides utilities for string naming convention conversions (camelCase, snake_case, PascalCase, etc.).

The package includes functions for converting between different naming conventions commonly used in programming,
such as camelCase, PascalCase, snake_case, and UPPER_SNAKE_CASE. It also provides helper functions for
capitalizing the first letter of strings.

Note: Current implementation only supports ASCII characters. For non-ASCII characters (e.g., Chinese, accented letters),
behavior may not be as expected. Case conversion applies to ASCII letters only; non-ASCII characters pass through unchanged.

Statelessness:
All conversion functions are pure: they hold no internal state, perform no caching, and are safe for concurrent use from multiple goroutines.

Exported Functions:
  - Pascal2Camel(name string) string - Convert PascalCase to camelCase (e.g., "UserName" → "userName")
  - Camel2Pascal(name string) string - Convert camelCase to PascalCase (e.g., "userName" → "UserName")
  - Pascal2Snake(name string) string - Convert PascalCase to snake_case (e.g., "UserName" → "User_Name")
  - Pascal2UpperSnake(name string) string - Convert PascalCase to UPPER_SNAKE_CASE (e.g., "UserName" → "USER_NAME")
  - Camel2UpperSnake(name string) string - Convert camelCase to UPPER_SNAKE_CASE (e.g., "userName" → "USER_NAME")
  - Camel2Snake(name string) string - Convert camelCase to snake_case (e.g., "userName" → "user_Name")
  - Snake2Pascal(name string) string - Convert snake_case to PascalCase (e.g., "user_name" → "UserName")
  - Snake2Camel(name string) string - Convert snake_case to camelCase (e.g., "user_name" → "userName")
  - ToSnake(name string) string - Convert PascalCase/camelCase to lower-case snake_case (e.g., "UserName" → "user_name")
  - Ucfirst(str string) string - Capitalize the first letter AND lowercase all remaining letters (destructive: "helloWorld" → "Helloworld"). See Lcfirst for a non-destructive variant.
  - Lcfirst(str string) string - Lowercase the first letter of a string (e.g., "Hello" → "hello") (non-destructive: only the first letter is changed)

Examples:

	// Convert from PascalCase to other formats
	result := stringx.Pascal2Snake("UserName")        // "User_Name"
	result := stringx.Pascal2Camel("UserName")        // "userName"
	result := stringx.Pascal2UpperSnake("UserName")   // "USER_NAME"

	// Convert from camelCase to other formats
	result := stringx.Camel2Snake("userName")         // "user_Name"
	result := stringx.Camel2Pascal("userName")        // "UserName"
	result := stringx.ToSnake("UserName")             // "user_name"

	// Convert from snake_case to other formats
	result := stringx.Snake2Pascal("user_name")       // "UserName"
	result := stringx.Snake2Camel("user_name")        // "userName"

	// Capitalization helpers
	result := stringx.Ucfirst("hello")                // "Hello"
	result := stringx.Lcfirst("Hello")                // "hello"

Note:

	// Current implementation preserves initial capitalization in snake_case conversions:
	result := stringx.Pascal2Snake("UserName")        // "User_Name" (not "user_name")
	result := stringx.Camel2Snake("userName")         // "user_Name" (first word lowercase, rest capitalized)

Known Limitations:

	// 1. camelCase input undergoes lossy conversion: abbreviation capitalization cannot be restored
	// For example, "xmlParser" cannot restore the "XMLParser" abbreviation casing:
	result := stringx.Pascal2Snake("xmlParser")       // "Xml_Parser" (not "XML_Parser")
	result := stringx.Pascal2Snake("XMLParser")       // "XML_Parser" (different from above)

	// 2. Bidirectional conversion does not guarantee restoration of original input
	// For example, snake_case → PascalCase loses abbreviation information:
	result := stringx.Snake2Pascal("xml_parser")      // "XmlParser" (not "XMLParser")
	// Converting back would not restore the original "XMLParser"

	// 3. Acronyms followed directly by a lower-case word are split lossily.
	// The splitter follows the community-standard boundary rule: a word starts
	// at an upper-case letter, and within a run of consecutive upper-case
	// letters the LAST one belongs to the next word ("XMLParser" -> XML|Parser).
	// Under this rule an acronym immediately followed by a lower-case word
	// cannot be detected and loses its last letter:
	//
	//	result := stringx.Pascal2Snake("HTTPserver")   // "HTT_Pserver"  (not "HTTP_server")
	//	result := stringx.Pascal2Snake("XMLparser")    // "XM_Lparser"   (not "XML_parser")
	//	result := stringx.Pascal2Snake("ABcdef")       // "A_Bcdef"      (not "AB_cdef")
	//	result := stringx.Pascal2Snake("HTTPSserver")  // "HTTP_Sserver" (not "HTTPS_server")
	//
	// This ambiguity is provably unresolvable without a dictionary:
	// "HTTPserver" (HTTP+server) and "XMLParsers" (XML+Parsers) share the same
	// case pattern yet require different splits. Prefer well-formed input with
	// an upper-case word boundary ("HTTPServer" -> "HTTP_Server") to avoid the loss.

	// 4. Inputs are assumed to be identifiers of a single naming convention,
	// composed of ASCII letters, digits and (for snake_case inputs)
	// underscores. Behavior for any other input is undefined:
	//   - characters other than letters/digits/underscores pass through
	//     unchanged: Pascal2Snake("user-name") -> "User-name"
	//   - consecutive, leading or trailing underscores are silently collapsed
	//     by Snake2Pascal/Snake2Camel: Snake2Pascal("user__name") -> "UserName"
	//   - UPPER_SNAKE input fed to Pascal2Snake/Pascal2UpperSnake/ToSnake
	//     yields doubled underscores: Pascal2Snake("USER_NAME") -> "USER__NAME"
	// Sanitize input before conversion if these cases are possible.
*/
package stringx
