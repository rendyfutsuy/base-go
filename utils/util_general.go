package utils

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/html"
)

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

// Replace the nth occurrence of old in s by new.
func ReplaceNth(s, old, new string, n int) string {
	i := 0
	for m := 1; m <= n; m++ {
		x := strings.Index(s[i:], old)
		if x < 0 {
			break
		}
		i += x
		if m == n {
			return s[:i] + new + s[i+len(old):]
		}
		i += len(old)
	}
	return s
}

func SafeCutSlice[T comparable](slice []T, start, end int) []T {
	if start < 0 {
		start = 0
	}
	if start > len(slice) {
		start = len(slice)
	}
	if end < start {
		end = start
	}
	if end > len(slice) {
		end = len(slice)
	}
	return slice[start:end]
}

// SafeRemoveIndexSlice safely removes an element from a slice at the specified index.
// If the index is out of bounds, it adjusts the index to be within the valid range.
//
// Parameters:
// - slice: The slice from which to remove the element.
// - index: The index of the element to remove.
//
// Returns:
// A new slice with the element at the specified index removed.
func SafeRemoveIndexSlice(slice interface{}, index int) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		panic("SafeRemoveIndexSlice: not a pointer to a slice")
	}

	v = v.Elem()
	length := v.Len()
	if index < 0 {
		index = 0
	}
	if index >= length {
		index = length
	}

	result := reflect.AppendSlice(v.Slice(0, index), v.Slice(index+1, length))
	v.Set(result)
}

// SafeRemoveValueSlice removes the first element from a slice that satifies the value.
//
// Parameters:
// - slice: The slice from which to remove the element.
// - value: The value of the element to remove.
//
// Returns:
// A new slice with the first element that satifies the value removed.
func SafeRemoveValueSlice[T comparable](s []T, val T) []T {
	for i, v := range s {
		if v == val {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s // value not found
}

// InArray checks if a given value is present in a slice.
// Notes that value and slice should have the same type
// It returns true if the value is found, otherwise it returns false.
//
// Parameters:
//   - val: The value to search for in the array.
//   - array: The slice to search within.
//
// Returns:
//   - ok: A boolean indicating whether the value was found in the array.
func InArray[T comparable](val T, array []T) (ok bool) {
	return slices.Contains(array, val)
}

// MinMax returns the min and max from a slice of ordered values.
// returns min, max
func MinMax[T Ordered](slice ...T) (T, T) {
	var def T
	if len(slice) == 0 {
		return def, def
	}

	min, max := slice[0], slice[0]
	for _, v := range slice[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

// CompareSlices: Compare Two Value of Slices with the same variable type, ignoring the order of the value in slice
func CompareSlices[T comparable](arr []T, arr2 []T) (ok bool) {
	if len(arr) != len(arr2) {
		return false
	}

	// Create maps to count occurrences of each string in both slices
	count1 := make(map[T]int)
	count2 := make(map[T]int)

	// Count frequencies of elements in the first slice
	for _, v := range arr {
		count1[v]++
	}

	// Count frequencies of elements in the second slice
	for _, v := range arr2 {
		count2[v]++
	}

	// Compare both maps
	for key, value := range count1 {
		if count2[key] != value {
			return false
		}
	}

	// Ensure there are no extra elements in count2 that aren't in count1
	for key := range count2 {
		if count1[key] != count2[key] {
			return false
		}
	}

	return true
}

// AlphabetIndex : Create Index of String based on Alphabet
func AlphabetIndex() []string {
	alphabetIndex := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	return alphabetIndex
}

func GetDifferences(array1, array2 []string) []string {
	// Buat map untuk menyimpan elemen-elemen dari array1
	// Ini digunakan untuk efisiensi dalam pencarian
	elementMap := make(map[string]bool)

	// Tambahkan elemen-elemen dari array1 ke map
	for _, elem := range array1 {
		elementMap[elem] = true
	}

	// Buat array untuk menyimpan perbedaan
	differences := []string{}

	// Iterasi melalui array2 dan cek setiap elemen
	// Jika tidak ada dalam map, itu adalah perbedaan
	for _, elem := range array2 {
		if _, exists := elementMap[elem]; !exists {
			differences = append(differences, elem)
		}
	}

	return differences
}

// IsNumericString memeriksa apakah string hanya berisi angka
func IsNumericString(input string) bool {
	// Gunakan ekspresi reguler untuk memeriksa hanya terdiri dari angka (0-9)
	numericRegex := regexp.MustCompile("^[0-9]+$")
	return numericRegex.MatchString(input)
}

// CheckStringValidity checks if the input string is valid.
// A valid string contains only digits and dots, and does not have a dot at the beginning or end.
//
// Parameters:
//
//	input (string): The string to be checked.
//
// Returns:
//
//	error: An error if the string is invalid, otherwise nil.
func CheckStringValidity(input string) error {
	if len(input) == 0 {
		return errors.New("String should not be empty")
	}

	for _, char := range input {
		// Mengecek apakah karakter adalah huruf atau bukan tanda titik
		if !unicode.IsDigit(char) && char != '.' {
			return errors.New("String should only contain digits and dots")
		}
	}

	// Pastikan bahwa tanda titik tidak berada di awal atau akhir string
	if input[0] == '.' || input[len(input)-1] == '.' {
		return errors.New("Dot should not be at the beginning or end of the string")
	}

	return nil
}

// CheckPercentValue checks if the given input is a valid percentage value.
// A valid percentage value is greater than 0 and less than or equal to 100.
//
// Parameters:
//
//	input (float64): The percentage value to check.
//
// Returns:
//
//	bool: true if the input is a valid percentage value, false otherwise.
func CheckPercentValue(input float64) bool {
	// MAIN VARIABLE
	res := true

	if input == 0 {
		res = false
	}

	if input > 100 {
		res = false
	}

	return res
}

// ToFixed formats a floating-point number to a specified number of decimal places.
// The number of decimal places is specified by the `comma` parameter, which should be a string representation of an integer.
// The function returns the formatted number as a float64.
//
// Parameters:
//   - comma: A string representing the number of decimal places to format the float to.
//   - v: The float64 value to be formatted.
//
// Returns:
//   - A float64 value formatted to the specified number of decimal places.
func ToFixed(comma string, v float64) float64 {
	// FLOAT TO STRING FORMATING DIGIT DECIMAL
	res := fmt.Sprintf("%."+comma+"f", v)
	// STRING TO FLOAT64
	v, _ = strconv.ParseFloat(res, 64)

	return v
}

// DecimalCheck checks if the number of decimal places in a given float64 value
// does not exceed the specified maximum number of decimal places.
//
// Parameters:
//   - max: The maximum number of decimal places allowed.
//   - value: The float64 value to be checked.
//
// Returns:
//   - bool: true if the number of decimal places is less than or equal to max,
//     false otherwise.
func DecimalCheck(max int, value float64) bool {
	stringFloat := strconv.FormatFloat(value, 'f', -1, 64)
	data := strings.Split(stringFloat, ".")
	if len(data) > 2 {
		return false
	}

	if len(data) == 1 {
		return true
	}

	return len(data[1]) <= max
}

func RoundingFloat64(x, y float64) float64 {
	precision := math.Pow(10, y)
	absX := math.Abs(x)

	if math.Mod(absX*precision, 1) >= 0.5 {
		if x < 0 {
			return -math.Ceil(absX*precision) / precision
		}
		return math.Ceil(x*precision) / precision
	}

	if x < 0 {
		return -math.Floor(absX*precision) / precision
	}

	return math.Floor(x*precision) / precision
}

func RoundingFloat64Padded(x, y float64) string {
	val := RoundingFloat64(x, y)

	// Use formatting to ensure zero-padding based on y
	format := fmt.Sprintf("%%.%df", int(y))
	return fmt.Sprintf(format, val)
}

func StartOfDay(t time.Time) time.Time {
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return date
}

func EndOfDay(t time.Time) time.Time {
	date := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())

	return date
}

func StartOfMonth(t time.Time) *time.Time {
	date := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)

	return &date
}

func EndOfMonth(t time.Time) *time.Time {
	parse := StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
	date := time.Date(t.Year(), parse.Month(), parse.Day(), 23, 59, 59, 59, t.Location())

	return &date
}

func StartOfYear(t time.Time) *time.Time {
	date := time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())

	return &date
}

func EndOfYear(t time.Time) *time.Time {
	parse := StartOfYear(t).AddDate(1, 0, 0).Add(-time.Nanosecond)
	date := time.Date(t.Year(), parse.Month(), parse.Day(), 23, 59, 59, 59, t.Location())

	return &date
}

// StringToUUIDSlice converts a comma-separated string of UUIDs into a slice of uuid.UUID
func StringToUUIDSlice(s string) ([]uuid.UUID, error) {
	var uuids []uuid.UUID
	if s == "" {
		return uuids, nil
	}

	parts := strings.Split(s, ",")
	for _, part := range parts {
		id, err := StringToUUID(part)
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, id)
	}

	return uuids, nil
}

// Haversine calculates the distance between two points on the Earth's surface
// given their latitude and longitude in degrees. It returns the distance in meters.
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000               // Earth radius in meters
	Radian1 := lat1 * math.Pi / 180 // convert to radians
	Radian2 := lat2 * math.Pi / 180
	LatDiffInRadian := (lat2 - lat1) * math.Pi / 180
	LonDiffInRadian := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(LatDiffInRadian/2)*math.Sin(LatDiffInRadian/2) + math.Cos(Radian1)*math.Cos(Radian2)*math.Sin(LonDiffInRadian/2)*math.Sin(LonDiffInRadian/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// Safely Get Value from a Pointer, return the zero value from the variable type if pointer is nil
// example zero value: String: "", int: 0
func GetPointerValue[T any](p *T) T {
	var zeroValue T // default value for each type
	if p == nil {
		return zeroValue
	}

	return *p
}

// Generate pointer of a value
func GetPointer[T any](p T) *T {
	// Handle cases where T itself can be nil (e.g., interface, slice, map, func, chan)
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Slice ||
		v.Kind() == reflect.Map || v.Kind() == reflect.Chan ||
		v.Kind() == reflect.Func || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
	}

	return &p
}

// Safely Get Value from a string Pointer, return the zero value from the variable type if pointer is nil
// example zero value: String: ""
func GetStringPointerValueNil(str *string) interface{} {
	if (str != nil && *str == "") || str == nil {
		return nil
	}

	return *str
}

// Safely Get Pointer String from a uuid var
func GetPointerStringUUID(id uuid.UUID) (res *string) {
	str := id.String()

	res = &str

	return res
}

// Reformat Numbers to Formatted Currency using default format
func ToFormattedCurrency(data interface{}) string {
	format := "#,###.##"
	switch val := data.(type) {
	case float64:
		if math.IsNaN(val) {
			return humanize.FormatFloat(format, 0)
		}

		return humanize.FormatFloat(format, val)
	case int32:
		return humanize.FormatInteger(format, int(val))
	case int64:
		return humanize.FormatInteger(format, int(val))
	case int:
		return humanize.FormatInteger(format, int(val))
	default:
		return fmt.Sprint(data)
	}
}

// ToTitleCase converts a string to Title Case.
func ToTitleCase(str string) string {
	words := strings.Fields(str)
	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

// ToSnakeCase converts a string to snake_case.
func ToSnakeCase(str string) string {
	// Replace all uppercase letters with the lowercase letter and prepend an underscore
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	str = re.ReplaceAllString(str, `${1}_${2}`)

	// Convert the whole string to lowercase
	return strings.ToLower(str)
}

func DuplicateCount(list []string) map[string]int {

	duplicate_frequency := make(map[string]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map

		if _, exist := duplicate_frequency[item]; exist {
			duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
		} else {
			duplicate_frequency[item] = 1 // else start counting from 1
		}
	}
	return duplicate_frequency
}

// Print all Field and Value in a Struct
// Used for debugging purpose
func PrintStructFields(s interface{}) {
	// Get the value of the struct
	v := reflect.ValueOf(s)

	// Ensure that the provided argument is a struct
	if v.Kind() == reflect.Struct {
		// Iterate through all the fields of the struct
		for i := 0; i < v.NumField(); i++ {
			// Get the field's name and value
			field := v.Type().Field(i)
			fieldValue := v.Field(i)

			// Check if the field is a pointer and handle it
			if fieldValue.Kind() == reflect.Ptr {
				// Check if the pointer is nil
				if !fieldValue.IsNil() {
					// Dereference the pointer and print the value
					fmt.Printf("%s: %v\n", field.Name, fieldValue.Elem())
				} else {
					// Print a message if the pointer is nil
					fmt.Printf("%s: nil\n", field.Name)
				}
			} else {
				// Print the value directly if it's not a pointer
				fmt.Printf("%s: %v\n", field.Name, fieldValue.Interface())
			}
		}
	} else {
		fmt.Println("Provided argument is not a struct")
	}
}

func ChunkSlice[T any](slice []T, size int) [][]T {
	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// RandomNumber generate a number between two integer describe by two params of the func
func RandomNumber(low, hi int) int {
	return low + rand.IntN(hi-low)
}

// sanitizeFloat takes a float64 value and returns it unchanged if it's a valid number.
// If the value is NaN (Not a Number) or an infinity (positive or negative), it returns 0.
// This helps ensure that the returned value is a finite and valid float64.
func SanitizeFloat(value float64) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0
	}
	return value
}

func IsEndorseTrx(documentNumber string, parentId uuid.UUID) bool {
	return EndorsedFacultativeRegex.MatchString(documentNumber) && parentId != uuid.Nil
}

func IsReplacedTrx(documentNumber string, parentId uuid.UUID) bool {
	return regexp.MustCompile(`-RL\d+$`).MatchString(documentNumber) && parentId != uuid.Nil
}

// check if the document number just contain `FAC` letter
// or the length of the document number is 15
func IsDocumentNumberFACParent(documentNumber string) bool {
	if len(documentNumber) == 15 {
		return true
	}

	return regexp.MustCompile(`^FAC`).MatchString(documentNumber)
}

// SafeDivide safely divides a by b (a / b), returning 0 if b is 0 or if a panic occurs.
func SafeDivide(a, b float64) (result float64) {
	defer func() {
		if r := recover(); r != nil {
			result = 0 // or any default value you want to return
		}
	}()
	if b == 0 {
		return 0
	}
	return a / b
}

// safeNumber returns the value if it is not NaN, otherwise returns 0
func SafeNumber(value float64) float64 {

	var (
		floatZeroThreshold float64
	)

	// Get the float zero threshold from the config
	valFloatZeroThreshold := ConfigVars.String("FLOAT_ZERO_THRESHOLD")
	if valFloatZeroThreshold == "" {
		floatZeroThreshold = 6 // Default to 6 decimal places if not set
	} else {
		// Convert the string to float64
		var err error
		floatZeroThreshold, err = strconv.ParseFloat(valFloatZeroThreshold, 64)
		if err != nil {
			Logger.Error(fmt.Sprintf("FixZero - failed to parse FLOAT_ZERO_THRESHOLD: %s, err: %v", valFloatZeroThreshold, err))
			floatZeroThreshold = 6 // Default to 6 decimal places if parsing fails
		}
	}

	// If the value is NaN, return 0.0
	if math.IsNaN(value) {
		Logger.Error("You Got NaN Value")
		return 0.0
	}

	// If the value is infinite, return 0.0
	if math.IsInf(value, 0) {
		return 0.0
	}

	if math.Abs(value) < math.Pow(10, -floatZeroThreshold) {
		// If the absolute value of f is less than the threshold, return 0.0
		// This is to avoid returning very small numbers that are close to zero.
		return 0.0
	}

	return value
}

// check if the document number is a risk document: FAC...-R...
func IsDocumentNumberRisk(documentNumber string) bool {
	// Matches: FAC followed by digits, then -R followed by digits, and nothing else
	matched, _ := regexp.MatchString(`^FAC\d+-R\d+$`, documentNumber)
	return matched
}

// IsDocumentNumberRiskExtend checks if the document number is an extended risk document.
// It matches patterns like: FAC...-R... followed by other characters.
var documentNumberRiskExtendRegex = regexp.MustCompile(`^FAC\d+-R\d+.*$`)

func IsDocumentNumberRiskExtend(documentNumber string) bool {
	// Matches: FAC followed by digits, then -R followed by digits, and can be followed by other characters.
	return documentNumberRiskExtendRegex.MatchString(documentNumber)
}

func IsDocumentNumberEndorse(documentNumber string) bool {
	// Matches: FAC followed by digits, then -E followed by digits, and nothing else
	matched, _ := regexp.MatchString(`^FAC\d+.*-E\d+$`, documentNumber)
	return matched
}

// indexComma finds the index of the first comma in a tag string (e.g., "name,omitempty")
func indexComma(tag string) int {
	for i := 0; i < len(tag); i++ {
		if tag[i] == ',' {
			return i
		}
	}
	return -1
}

func GetHostName(c echo.Context) (url string) {
	// Get scheme: "http" or "https"
	scheme := "http"
	if c.Request().TLS != nil {
		scheme = "https"
	}

	host := c.Request().Host

	return scheme + "://" + host
}

func RemoveHTMLTags(str string) string {
	doc, err := html.Parse(strings.NewReader(str))
	if err != nil {
		return str
	}
	var text strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return RemoveNonASCII(text.String())
}

func RemoveNonASCII(s string) string {
	// keep only ASCII characters
	result := make([]rune, 0, len(s))
	for _, r := range s {
		if r <= unicode.MaxASCII {
			result = append(result, r)
		}
	}
	return string(result)
}

// findMissingColumns checks which elements from `required` are missing in `actual`.
func FindMissingColumns(actual []string, required []string) []string {
	// Create a set of the actual columns for efficient lookup.
	// Using struct{} is a Go idiom for a set where only the key matters.
	actualColumnSet := make(map[string]struct{}, len(actual))
	for _, col := range actual {
		actualColumnSet[col] = struct{}{}
	}

	var missing []string
	// Iterate through the required columns and check if they exist in the set.
	for _, reqCol := range required {
		if _, exists := actualColumnSet[reqCol]; !exists {
			missing = append(missing, reqCol)
		}
	}

	return missing
}
