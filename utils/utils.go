package utils

import "regexp"


func ExtractEmailMentions(input string) ([]string , error){
	re, err := regexp.Compile(`\B@[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}\b`)
	if err != nil {
		return nil, err
	}

	matches := re.FindAllString(input, -1)
	for i, match := range matches {
		matches[i] = match[1:]
	}

	return matches, nil
}


func RemoveDuplicates(elements []string) []string {
    encountered := map[string]bool{}
    result := []string{}

    for v := range elements {
        if encountered[elements[v]] {
            // Do not add duplicate.
        } else {
            encountered[elements[v]] = true
            result = append(result, elements[v])
        }
    }
    return result
}