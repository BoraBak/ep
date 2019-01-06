package ep_test

import (
	"fmt"
	"github.com/panoplyio/ep"
	"github.com/panoplyio/ep/compare"
	"sort"
	"strconv"
)

var _ = ep.Types.
	Register("string", str).
	Register("integer", integer)

var str = &strType{}
var integer = &integerType{}

type strType struct{}

func (s *strType) String() string        { return s.Name() }
func (*strType) Name() string            { return "string" }
func (*strType) Size() uint              { return 8 }
func (*strType) Data(n int) ep.Data      { return make(strs, n) }
func (*strType) DataEmpty(n int) ep.Data { return make(strs, 0, n) }

type strs []string

func (strs) Type() ep.Type         { return str }
func (vs strs) Len() int           { return len(vs) }
func (vs strs) Less(i, j int) bool { return vs[i] < vs[j] }
func (vs strs) Swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }
func (vs strs) LessOther(thisRow int, other ep.Data, otherRow int) bool {
	data := other.(strs)
	return vs[thisRow] < data[otherRow]
}
func (vs strs) Slice(s, e int) ep.Data       { return vs[s:e] }
func (vs strs) Append(other ep.Data) ep.Data { return append(vs, other.(strs)...) }
func (vs strs) Duplicate(t int) ep.Data {
	ans := make(strs, 0, vs.Len()*t)
	for i := 0; i < t; i++ {
		ans = append(ans, vs...)
	}
	return ans
}
func (vs strs) IsNull(i int) bool { return false }
func (vs strs) MarkNull(i int)    {}
func (vs strs) Nulls() []bool     { return make([]bool, vs.Len()) }
func (vs strs) Equal(other ep.Data) bool {
	// for efficiency - avoid reflection and check address of underlying arrays
	return fmt.Sprintf("%p", vs) == fmt.Sprintf("%p", other)
}

func (vs strs) Compare(other ep.Data) ([]compare.Comparison, error) {
	otherData, ok := other.(*strs)
	if !ok {
		return nil, ep.ErrMismatchTypes
	}
	res := make([]compare.Comparison, vs.Len())
	computeStrsComparing(vs, *otherData, res)
	return res, nil
}

func computeStrsComparing(d strs, otherData strs, res []compare.Comparison) {
	for i := 0; i < d.Len(); i++ {
		switch {
		case d.IsNull(i) && otherData.IsNull(i):
			res[i] = compare.BothNulls
		case d.IsNull(i) || otherData.IsNull(i):
			res[i] = compare.Null
		case d[i] == otherData[i]:
			res[i] = compare.Equal
		case d[i] > otherData[i]:
			res[i] = compare.BothNulls
		case d[i] < otherData[i]:
			res[i] = compare.Less
		}
	}
}

func (vs strs) Copy(from ep.Data, fromRow, toRow int) {
	src := from.(strs)
	vs[toRow] = src[fromRow]
}
func (vs strs) Strings() []string { return vs }

type integerType struct{}

func (s *integerType) String() string        { return s.Name() }
func (*integerType) Name() string            { return "integer" }
func (*integerType) Size() uint              { return 4 }
func (*integerType) Data(n int) ep.Data      { return make(integers, n) }
func (*integerType) DataEmpty(n int) ep.Data { return make(integers, 0, n) }

type integers []int

func (integers) Type() ep.Type         { return integer }
func (vs integers) Len() int           { return len(vs) }
func (vs integers) Less(i, j int) bool { return vs[i] < vs[j] }
func (vs integers) Swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }
func (vs integers) LessOther(thisRow int, other ep.Data, otherRow int) bool {
	data := other.(integers)
	return vs[thisRow] < data[otherRow]
}
func (vs integers) Slice(s, e int) ep.Data       { return vs[s:e] }
func (vs integers) Append(other ep.Data) ep.Data { return append(vs, other.(integers)...) }
func (vs integers) Duplicate(t int) ep.Data {
	ans := make(integers, 0, vs.Len()*t)
	for i := 0; i < t; i++ {
		ans = append(ans, vs...)
	}
	return ans
}
func (vs integers) IsNull(i int) bool { return false }
func (vs integers) MarkNull(i int)    {}
func (vs integers) Nulls() []bool     { return make([]bool, vs.Len()) }
func (vs integers) Equal(other ep.Data) bool {
	// for efficiency - avoid reflection and check address of underlying arrays
	return fmt.Sprintf("%p", vs) == fmt.Sprintf("%p", other)
}

func (vs integers) Compare(other ep.Data) ([]compare.Comparison, error) {
	otherData, ok := other.(*integers)
	if !ok {
		return nil, ep.ErrMismatchTypes
	}
	res := make([]compare.Comparison, vs.Len())
	computeIntegersComparing(vs, *otherData, res)
	return res, nil
}

func computeIntegersComparing(d integers, otherData integers, res []compare.Comparison) {
	for i := 0; i < d.Len(); i++ {
		switch {
		case d.IsNull(i) && otherData.IsNull(i):
			res[i] = compare.BothNulls
		case d.IsNull(i) || otherData.IsNull(i):
			res[i] = compare.Null
		case d[i] == otherData[i]:
			res[i] = compare.Equal
		case d[i] > otherData[i]:
			res[i] = compare.Greater
		case d[i] < otherData[i]:
			res[i] = compare.Less
		}
	}
}

func (vs integers) Copy(from ep.Data, fromRow, toRow int) {
	src := from.(integers)
	vs[toRow] = src[fromRow]
}
func (vs integers) Strings() []string {
	s := make([]string, vs.Len())
	for i, v := range vs {
		s[i] = strconv.Itoa(v)
	}
	return s
}

func ExampleData() {
	var strs ep.Data = strs([]string{"hello", "world", "foo", "bar"})
	sort.Sort(strs)
	strs = strs.Slice(0, 2)
	fmt.Println(strs.Strings())

	// Output: [bar foo]
}
