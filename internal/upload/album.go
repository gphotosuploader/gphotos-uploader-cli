package upload

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// albumName returns the album name based on the configured parameter.
// albumName returns the album name based on the configured parameter.
func (job *UploadFolderJob) albumName(filePath, absoluteFilePath string) string {
	before, after, found := strings.Cut(job.Album, ":")
	if !found {
		return ""
	}
	if before == "name" {
		return after
	}

	if before == "template" {
		val, err := albumNameUsingTemplate(after, filePath, absoluteFilePath)
		if err != nil {
			panic("invalid Albums name template format - " + err.Error())
		}

		return val
	}

	if before == "auto" {
		switch after {
		case "folderPath":
			return albumNameUsingFolderPath(filePath)
		case "folderName":
			return albumNameUsingFolderName(filePath)
		default:
			panic("invalid Albums parameter")
		}
	}

	return ""
}

// albumNameUsingFolderPath returns an AlbumID name using the full Path of the given folder.
func albumNameUsingFolderPath(path string) string {
	p := filepath.Dir(path)
	if p == "." {
		return ""
	}

	p = strings.ReplaceAll(p, "/", "_")

	// In path starts with '/' remove it before.
	if p[0] == '_' {
		return p[1:]
	}
	return p
}

// albumNameUsingFolderName returns an AlbumID name using the name of the given folder.
func albumNameUsingFolderName(path string) string {
	p := filepath.Dir(path)
	if p == "." {
		return ""
	}
	return filepath.Base(p)
}

// albumNameUsingTemplate returns an AlbumID name using the given template.
func albumNameUsingTemplate(template, filePath, absoluteFilePath string) (string, error) {
	//TODO: implement pass time creation
	return parseAlbumNameTemplate(template, filePath, time.Now())
}

// Recursively parse the template and replace the tokens with the corresponding values.
func parseAlbumNameTemplate(template, filePath string, fileCreateTime time.Time) (string, error) {
	outputs := ""
	i := 0
	for {
		tokenName := getTokenName(template[i:])
		if tokenName != "" {
			tokenNameLen := len(tokenName) + 3
			val, err := replaceTemplateToken(tokenName, filePath, fileCreateTime)
			if err != nil {
				return "", err
			}

			outputs += val
			i += tokenNameLen
		}

		functionName := getTemplateFunctionName(template[i:])
		if functionName != "" {
			functionArgStart := i + len(functionName) + 2
			i = functionArgStart
			functionDepth := 1
			args := []string{}
			currentArg := ""
			for i < len(template) {
				if template[i] == '(' {
					functionDepth++
				}

				if template[i] == ')' {
					functionDepth--
				}

				if (template[i] == ',' && functionDepth == 1) || functionDepth == 0 {
					val, err := parseAlbumNameTemplate(currentArg, filePath, fileCreateTime)
					if err != nil {
						return "", err
					}

					args = append(args, val)
					currentArg = ""
				} else {
					currentArg += string(template[i])
				}

				i++
				if functionDepth == 0 {
					//empty first argument is 0 args
					if len(args) == 1 && args[0] == "" {
						args = []string{}
					}
					val, err := runTemplateFunction(functionName, args)
					if err != nil {
						return "", err
					}

					outputs += val
					break
				}
			}

			if functionDepth != 0 {
				return "", fmt.Errorf("function missing closing parenthesis")
			}
		}

		if i == len(template) {
			break
		}
		outputs += string(template[i])
		i++
	}

	return outputs, nil
}

func getTemplateFunctionName(template string) string {
	// perf optimization to avoid regex if not needed
	if (len(template) < 4) || (template[0] != '$') {
		return ""
	}

	re := regexp.MustCompile(`^\$\b([a-zA-Z]+)\b\(`)
	match := re.FindStringSubmatch(template)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func getTokenName(template string) string {
	// perf optimization to avoid regex if not needed
	if (len(template) < 4) || (template[0] != '%') || (template[1] != '_') {
		return ""
	}

	re := regexp.MustCompile(`^%_([a-zA-Z_]+)%`)
	match := re.FindStringSubmatch(template)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func runTemplateFunction(name string, args []string) (string, error) {
	switch name {
	case "cutLeft", "cutRight":
		if len(args) != 2 {
			return "", fmt.Errorf("%s requires 2 arguments", name)
		}

		cutN, err := strconv.Atoi(strings.TrimSpace(args[1]))

		if err != nil {
			return "", fmt.Errorf("%s requires a number as second argument", name)
		}

		if cutN >= len(args[0]) {
			return "", nil
		}

		if name == "cutLeft" {
			return args[0][cutN:], nil
		}

		if name == "cutRight" {
			return args[0][:len(args[0])-cutN], nil
		}
	case "lower", "upper", "sentence", "title":
		if len(args) != 1 {
			return "", fmt.Errorf("%s requires 1 argument", name)
		}

		if len(args[0]) == 0 {
			return "", nil
		}
		switch name {
		case "lower":
			return strings.ToLower(args[0]), nil
		case "upper":
			return strings.ToUpper(args[0]), nil
		case "sentence":
			runes := []rune(strings.ToLower(args[0]))
			return strings.ToUpper(string(runes[0])) + string(runes[1:]), nil
		case "title":
			caser := cases.Title(language.English)
			titleStr := caser.String(args[0])
			return titleStr, nil
		}
	default:
		return "", fmt.Errorf("unknown function: " + name)
	}

	return "", nil
}

func replaceTemplateToken(token, filePath string, fileCreateTime time.Time) (string, error) {
	switch token {
	case "folderpath":
		return albumNameUsingFolderPath(filePath), nil
	case "directory":
		return albumNameUsingFolderName(filePath), nil
	case "parent_directory":
		return albumNameUsingFolderName(filepath.Dir(filePath)), nil
	case "month":
		return fileCreateTime.Format("01"), nil
	case "day":
		return fileCreateTime.Format("02"), nil
	case "year":
		return fileCreateTime.Format("2006"), nil
	case "time":
		return fileCreateTime.Format("15:04:05"), nil
	case "time_en":
		return fileCreateTime.Format("03:04:05 PM"), nil
	}

	return "", fmt.Errorf("invalid token: %s", token)
}
