package workflow_test

import (
	"bytes"
	"testing"

	. "github.com/slabgorb/wk/workflow"
)

var exampleTSV = `
Rank	Port	Volume 2015	Volume 2014	Volume 2013	Volume 2012	Volume 2011
18	Hamburg, Germany	8.82	9.73	9.30	8.89	9.01
19	Los Angeles, U.S.A.	8.16	8.33	7.87	8.08	7.94
20*	Keihin Ports, Japan	7.52	7.85	7.81	7.85	7.64
21	Long Beach, U.S.A.	7.19	6.82	6.73	6.05	6.06
22	Laem Chabang, Thailand	6.82	6.58	6.04	5.93	5.73
23	New York-New Jersey, U.S.A.	6.37	5.77	5.47	5.53	5.50
24	Yingkou, China	5.92	5.77	5.30	4.85	4.03
25	Bremen/Bremerhaven, Germany	5.48	5.78	5.84	6.13	5.92
26	Ho Chi Minh, Vietnam	5.31	6.39	5.96	5.19	4.53
27	Tanjung Priok, Jakarta, Indonesia	5.20	5.77	5.47	5.53	5.50
`

func TestReadTSV(t *testing.T) {
	buf := bytes.NewBufferString(exampleTSV)
	data, err := ReadTSV(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(data) != 10 {
		t.Errorf("expected %d rows, got %d", 10, len(data))
	}
	if len(data[0]) != 7 {
		t.Errorf("expected %d columns, got %d", 7, len(data[0]))
	}
}
