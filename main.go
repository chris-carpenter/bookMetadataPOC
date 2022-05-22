package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var parser = argparse.NewParser("mergeJson", "Merges JSON")
	pretty := parser.Flag("p", "pretty", &argparse.Options{Required: false, Help: "Pretty output"})
	debugLevel := parser.Selector("d", "debug-level", []string{"INFO", "DEBUG", "ERROR"}, &argparse.Options{Required: false, Help: "Logging debug level", Default: "ERROR"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}
	initPrettyLog(*debugLevel, *pretty)

	end := Index{
		Keys: Key{
			Kindle: 3,
			Audio:  "00:16:00",
		},
	}
	dirs, err := ioutil.ReadDir("./in")
	if err != nil {
		log.Error().Err(err).Msg("Failed to read './in'")
	}
	for _, dir := range dirs {
		merged := m{}
		files, err := ioutil.ReadDir(fmt.Sprintf("./in/%s", dir.Name()))
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read ./in/%s", dir.Name())
		}
		for _, file := range files {
			log.Info().Str("file.Name()", file.Name())
			jsonFileName := fmt.Sprintf("./in/%s/%s", dir.Name(), file.Name())
			jsonFile, err := os.Open(jsonFileName)
			if err != nil {
				log.Error().Err(err).Str("filename", jsonFileName).Msg("Unable to open json file")
			}
			log.Info().Str("filename", jsonFileName).Msg("Successfully opened file")

			byteValue, _ := ioutil.ReadAll(jsonFile)

			keys := Index{}
			if err := json.Unmarshal(byteValue, &keys); err != nil {
				log.Error().Err(err).Str("filename", jsonFileName).Msg("Unable to parse keys")
			}
			if keys.GreaterThan(end) {
				log.Info().Str("breakFile", file.Name()).Msg("Stopped before file.")
				break
			}
			out := m{}
			if err := json.Unmarshal(byteValue, &out); err != nil {
				log.Error().Err(err).Str("filename", jsonFileName).Msg("Failed to unmarshal json")
			}

			//log.Info().Msgf("%s - %+v", jsonFileName, out)
			if err := jsonFile.Close(); err != nil {
				log.Error().Err(err).Str("filename", jsonFileName).Msg("Failed to close file")
			}
			merged = mergeKeys(merged, out)
		}
		log.Info().Msgf("merged - %+v", merged)
	}
}

func initPrettyLog(debugLevel string, pretty bool) {

	switch debugLevel {
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	if pretty {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("| %s:", i)
		}
		output.FormatFieldValue = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		}

		log.Logger = zerolog.New(output).With().Timestamp().Logger()
	}
}

func prettyString(b string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(b), "", "    "); err != nil {
		return ""
	}
	return prettyJSON.String()
}

type m = map[string]interface{}

func mergeKeys(left, right m) m {
	for key, rightVal := range right {
		if leftVal, present := left[key]; present && (reflect.ValueOf(leftVal).Kind() == reflect.Map || reflect.ValueOf(leftVal).Kind() == reflect.Slice) {
			//then we don't want to replace it - recurse
			left[key] = mergeKeys(leftVal.(m), rightVal.(m))
		} else {
			// key not in left so we can just shove it in
			left[key] = rightVal
		}
	}
	return left
}
