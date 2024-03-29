package config

import (
	"bytes"
	"os"
	"strings"
)

// SaveConfigFile writes configuration file to local file system
func SaveConfigFile(c *ConfigFile, filename string) (err error) {
	// Write configuration file by filename.
	var f *os.File
	if f, err = os.Create(filename); err != nil {
		return err
	}

	equalSign := "="
	if c.PrettyFormat {
		equalSign = " = "
	}

	buf := bytes.NewBuffer(nil)
	for _, section := range c.sectionList {
		// Write section comments.
		if len(c.GetSectionComments(section)) > 0 {
			if _, err = buf.WriteString(c.GetSectionComments(section) + LineBreak); err != nil {
				return err
			}
		}

		if section != DEFAULT_SECTION {
			// Write section name.
			if _, err = buf.WriteString("[" + section + "]" + LineBreak); err != nil {
				return err
			}
		}

		for _, key := range c.keyList[section] {
			if key != " " {
				// Write key comments.
				if len(c.GetKeyComments(section, key)) > 0 {
					if _, err = buf.WriteString(c.GetKeyComments(section, key) + LineBreak); err != nil {
						return err
					}
				}

				keyName := key
				// Check if it's auto increment.
				if keyName[0] == '#' {
					keyName = "-"
				}

				// [SWH|+]:支持键名包含等号和冒号
				if c.KeyQuote {
					if strings.Contains(keyName, `=`) || strings.Contains(keyName, `:`) {
						if strings.Contains(keyName, "`") {
							if strings.Contains(keyName, `"`) {
								keyName = `"""` + keyName + `"""`
							} else {
								keyName = `"` + keyName + `"`
							}
						} else {
							keyName = "`" + keyName + "`"
						}
					}
				}
				//[SWH|+];

				value := c.data[section][key]

				//[SWH|+]:支持值包含等号和冒号
				if c.ValueQuote {
					if strings.Contains(value, `=`) || strings.Contains(value, `:`) {
						if strings.Contains(value, "`") {
							if strings.Contains(value, `"`) {
								value = `"""` + value + `"""`
							} else {
								value = `"` + value + `"`
							}
						} else {
							value = "`" + value + "`"
						}
					}
				}
				//[SWH|+];

				// Write key and value.
				if _, err = buf.WriteString(keyName + equalSign + value + LineBreak); err != nil {
					return err
				}
			}
		}

		// Put a line between sections.
		if _, err = buf.WriteString(LineBreak); err != nil {
			return err
		}
	}

	if _, err = buf.WriteTo(f); err != nil {
		return err
	}
	return f.Close()
}
