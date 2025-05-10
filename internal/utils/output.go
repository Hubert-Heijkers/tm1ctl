package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
)

func printPrettyJSON(data any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ") // 2-space indentation
	return enc.Encode(data)
}

func stringify(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case float64, bool, int:
		return fmt.Sprint(val)
	default:
		b, _ := json.Marshal(val)
		return string(b)
	}
}

func printArrayTable(list []any) error {
	if len(list) == 0 {
		fmt.Println("No data.")
		return nil
	}

	// Collect headers from the first object
	first := list[0].(map[string]any)
	headers := make([]string, 0, len(first))
	for k := range first {
		headers = append(headers, k)
	}

	// Softing the 'headers' for consistency, idealy they'd use the predefined order from the CSDL and, perhaps later, whatever is specified in a select option
	sort.Strings(headers)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	for _, item := range list {
		row := make([]string, len(headers))
		obj := item.(map[string]any)
		for i, key := range headers {
			row[i] = stringify(obj[key])
		}
		table.Append(row)
	}

	table.Render()
	return nil
}

func printTable(data any) error {
	switch val := data.(type) {
	case []any:
		return printArrayTable(val)
	case map[string]any:
		// Wrap the object as a one-item array
		return printArrayTable([]any{val})
	default:
		return fmt.Errorf("unsupported data type: %s", reflect.TypeOf(data))
	}
}

func Output(data any) error {
	switch viper.GetString("output-format") {
	case "table":
		return printTable(data)
	case "json":
		return printPrettyJSON(data)
	}
	return fmt.Errorf("invalid output format specified: %s", viper.GetString("output-format"))
}

func OutputEntity(data any) error {
	obj, ok := data.(map[string]any)
	if !ok {
		return errors.New("expected object at top level")
	}

	// Remove the @odata.context control information
	delete(obj, "@odata.context")

	return Output(obj)
}

func OutputCollection(data any) error {
	obj, ok := data.(map[string]any)
	if !ok {
		return errors.New("expected object at top level to extract 'value'")
	}

	selected, exists := obj["value"]
	if !exists {
		return errors.New("'value' not found in response")
	}

	return Output(selected)
}
