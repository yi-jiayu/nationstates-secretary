package nationstates

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestUnmarshalConsequences(t *testing.T) {
	s := `<NATION id="wilbert">
  <ISSUE id="369" choice="1">
    <OK>1</OK>
    <DESC>companies balk at paying their workers</DESC>
    <RANKINGS>
      <RANK id="4">
        <SCORE>6.29</SCORE>
        <CHANGE>1.15</CHANGE>
        <PCHANGE>22.373541</PCHANGE>
      </RANK>
      <RANK id="5">
        <SCORE>21.08</SCORE>
        <CHANGE>0.07</CHANGE>
        <PCHANGE>0.333175</PCHANGE>
      </RANK>
    </RANKINGS>
    <HEADLINES>
      <HEADLINE>Retailers Welcome Tax Cut</HEADLINE>
      <HEADLINE>Aristocrats Welcome Rising Income Inequality</HEADLINE>
      <HEADLINE>Lemonade Stand Children Accused Of Price-Fixing</HEADLINE>
      <HEADLINE>School Bans Chess As &#8220;Too Passive&#8221;</HEADLINE>
    </HEADLINES>
  </ISSUE>
</NATION>
`
	var n Nation
	err := xml.Unmarshal([]byte(s), &n)
	if err != nil {
		t.Fatal(err)
	}
	want := Consequences{
		Desc: "companies balk at paying their workers",
		Rankings: []Rank{
			{Score: 6.29, Change: 1.15, PChange: 22.37354},
			{Score: 21.08, Change: 0.07, PChange: 0.333175},
		},
		Headlines: []string{
			"Retailers Welcome Tax Cut",
			"Aristocrats Welcome Rising Income Inequality",
			"Lemonade Stand Children Accused Of Price-Fixing",
			"School Bans Chess As “Too Passive”",
		},
	}
	if got := n.Consequences; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, wanted %v", got, want)
	}
}
