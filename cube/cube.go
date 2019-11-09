package cube

import (
	"blendcube/conf"
	"bytes"
	"fmt"
	"math"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/unlit"
	"github.com/westphae/quaternion"
)

type (
	// State キューブの状態を保持する構造体
	State struct {
		UBL gltf.Node `json:"UBL"`
		UB  gltf.Node `json:"UB"`
		UBR gltf.Node `json:"UBR"`
		UL  gltf.Node `json:"UL"`
		U   gltf.Node `json:"U"`
		UR  gltf.Node `json:"UR"`
		UFL gltf.Node `json:"UFL"`
		UF  gltf.Node `json:"UF"`
		UFR gltf.Node `json:"UFR"`
		BL  gltf.Node `json:"BL"`
		B   gltf.Node `json:"B"`
		BR  gltf.Node `json:"BR"`
		L   gltf.Node `json:"L"`
		R   gltf.Node `json:"R"`
		FL  gltf.Node `json:"FL"`
		F   gltf.Node `json:"F"`
		FR  gltf.Node `json:"FR"`
		DBL gltf.Node `json:"DBL"`
		DB  gltf.Node `json:"DB"`
		DBR gltf.Node `json:"DBR"`
		DL  gltf.Node `json:"DL"`
		D   gltf.Node `json:"D"`
		DR  gltf.Node `json:"DR"`
		DFL gltf.Node `json:"DFL"`
		DF  gltf.Node `json:"DF"`
		DFR gltf.Node `json:"DFR"`
	}

	// RotationMark 回転記号
	RotationMark int
)

const (
	// c* clockwise 90°
	// cc* counter-clockwise 90°
	// dc* clockwise 180°
	unknown RotationMark = iota
	cU
	cD
	cF
	cB
	cL
	cR
	ccU
	ccD
	ccF
	ccB
	ccL
	ccR
	dcU
	dcD
	dcF
	dcB
	dcL
	dcR
)

var (
	qU  = quaternion.FromEuler(0, math.Pi/2, 0)
	qD  = quaternion.FromEuler(0, -math.Pi/2, 0)
	qF  = quaternion.FromEuler(-math.Pi/2, 0, 0)
	qB  = quaternion.FromEuler(math.Pi/2, 0, 0)
	qL  = quaternion.FromEuler(0, 0, -math.Pi/2)
	qR  = quaternion.FromEuler(0, 0, -math.Pi/2)
	qU2 = quaternion.FromEuler(0, math.Pi, 0)
	qF2 = quaternion.FromEuler(-math.Pi, 0, 0)
	qL2 = quaternion.FromEuler(0, 0, -math.Pi)

	doc   *gltf.Document
	state *State
)

// Initial キューブの状態を生成する
func Initial() error {
	gltf.RegisterExtension(unlit.ExtUnlit, unlit.New)

	var err error
	doc, err = gltf.Open(conf.Shared.Model.FilePath)
	if err != nil {
		return err
	}

	if len(doc.Nodes) != 26 {
		return fmt.Errorf("insufficient number of nodes")
	}

	state = new(State)
	for _, node := range doc.Nodes {
		switch node.Name {
		case "01_UBL":
			state.UBL = node
		case "02_UB":
			state.UB = node
		case "03_UBR":
			state.UBR = node
		case "04_UL":
			state.UL = node
		case "05_U":
			state.U = node
		case "06_UR":
			state.UR = node
		case "07_UFL":
			state.UFL = node
		case "08_UF":
			state.UF = node
		case "09_UFR":
			state.UFR = node
		case "10_BL":
			state.BL = node
		case "11_B":
			state.B = node
		case "12_BR":
			state.BR = node
		case "13_L":
			state.L = node
		case "14_R":
			state.R = node
		case "15_FL":
			state.FL = node
		case "16_F":
			state.F = node
		case "17_FR":
			state.FR = node
		case "18_DBL":
			state.DBL = node
		case "19_DB":
			state.DB = node
		case "20_DBR":
			state.DBR = node
		case "21_DL":
			state.DL = node
		case "22_D":
			state.D = node
		case "23_DR":
			state.DR = node
		case "24_DFL":
			state.DFL = node
		case "25_DF":
			state.DF = node
		case "26_DFR":
			state.DFR = node
		default:
			return fmt.Errorf("unexpected node name: %s", node.Name)
		}
	}

	return nil
}

// Generate ノードを生成する
func Generate(
	algorithm []string,
	isBinary bool,
	isUnlit bool,
) ([]byte, error) {
	d := *doc
	s := *state

	if isUnlit {
		d.ExtensionsUsed = []string{unlit.ExtUnlit}
		d.ExtensionsRequired = []string{unlit.ExtUnlit}

		var materials []gltf.Material
		for _, m := range d.Materials {
			m.Extensions = map[string]interface{}{
				unlit.ExtUnlit: struct{}{},
			}
			materials = append(materials, m)
		}
		d.Materials = materials
	}

	for _, a := range algorithm {
		s.Rotate(parseAlg(a))
	}

	d.Nodes = []gltf.Node{
		s.UBL, s.UB, s.UBR, s.UL, s.U, s.UR, s.UFL, s.UF, s.UFR,
		s.BL, s.B, s.BR, s.L, s.R, s.FL, s.F, s.FR,
		s.DBL, s.DB, s.DBR, s.DL, s.D, s.DR, s.DFL, s.DF, s.DFR,
	}

	buffer := new(bytes.Buffer)
	e := gltf.NewEncoder(buffer)
	e.AsBinary = isBinary
	if err := e.Encode(&d); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Rotate キューブを回転する
func (s *State) Rotate(mark RotationMark) {
	switch mark {
	case cU, ccU, dcU:
		s.UBL, s.UB, s.UBR, s.UL, s.U, s.UR, s.UFL, s.UF, s.UFR =
			move(mark, s.UBL, s.UB, s.UBR, s.UL, s.U, s.UR, s.UFL, s.UF, s.UFR)
	case cD, ccD, dcD:
		s.DFL, s.DF, s.DFR, s.DL, s.D, s.DR, s.DBL, s.DB, s.DBR =
			move(mark, s.DFL, s.DF, s.DFR, s.DL, s.D, s.DR, s.DBL, s.DB, s.DBR)
	case cF, ccF, dcF:
		s.UFL, s.UF, s.UFR, s.FL, s.F, s.FR, s.DFL, s.DF, s.DFR =
			move(mark, s.UFL, s.UF, s.UFR, s.FL, s.F, s.FR, s.DFL, s.DF, s.DFR)
	case cB, ccB, dcB:
		s.UBR, s.UB, s.UBL, s.BR, s.B, s.BL, s.DBR, s.DB, s.DBL =
			move(mark, s.UBR, s.UB, s.UBL, s.BR, s.B, s.BL, s.DBR, s.DB, s.DBL)
	case cL, ccL, dcL:
		s.UBL, s.UL, s.UFL, s.BL, s.L, s.FL, s.DBL, s.DL, s.DFL =
			move(mark, s.UBL, s.UL, s.UFL, s.BL, s.L, s.FL, s.DBL, s.DL, s.DFL)
	case cR, ccR, dcR:
		s.UFR, s.UR, s.UBR, s.FR, s.R, s.BR, s.DFR, s.DR, s.DBR =
			move(mark, s.UFR, s.UR, s.UBR, s.FR, s.R, s.BR, s.DFR, s.DR, s.DBR)
	}
}

func parseAlg(alg string) RotationMark {
	if len(alg) != 1 && len(alg) != 2 {
		return unknown
	}

	var mark RotationMark
	switch alg[0] {
	case 'U':
		mark = cU
	case 'D':
		mark = cD
	case 'F':
		mark = cF
	case 'B':
		mark = cB
	case 'L':
		mark = cL
	case 'R':
		mark = cR
	default:
		return unknown
	}

	if len(alg) == 2 {
		if alg[1] == '\'' {
			mark += 6
		} else if alg[1] == '2' {
			mark += 12
		} else {
			return unknown
		}
	}

	return mark
}

func move(mark RotationMark,
	n1 gltf.Node, n2 gltf.Node, n3 gltf.Node,
	n4 gltf.Node, n5 gltf.Node, n6 gltf.Node,
	n7 gltf.Node, n8 gltf.Node, n9 gltf.Node,
) (
	gltf.Node, gltf.Node, gltf.Node,
	gltf.Node, gltf.Node, gltf.Node,
	gltf.Node, gltf.Node, gltf.Node,
) {
	n1.Rotation = rotate(mark, n1.RotationOrDefault())
	n2.Rotation = rotate(mark, n2.RotationOrDefault())
	n3.Rotation = rotate(mark, n3.RotationOrDefault())
	n4.Rotation = rotate(mark, n4.RotationOrDefault())
	n5.Rotation = rotate(mark, n5.RotationOrDefault())
	n6.Rotation = rotate(mark, n6.RotationOrDefault())
	n7.Rotation = rotate(mark, n7.RotationOrDefault())
	n8.Rotation = rotate(mark, n8.RotationOrDefault())
	n9.Rotation = rotate(mark, n9.RotationOrDefault())

	switch mark {
	case cU, cD, cF, cB, cL, cR:
		return n7, n4, n1,
			n8, n5, n2,
			n9, n6, n3
	case ccU, ccD, ccF, ccB, ccL, ccR:
		return n3, n6, n9,
			n2, n5, n8,
			n1, n4, n7
	case dcU, dcD, dcF, dcB, dcL, dcR:
		return n9, n8, n7,
			n6, n5, n4,
			n3, n2, n1
	default:
		return n1, n2, n3, n4, n5, n6, n7, n8, n9
	}
}

func rotate(mark RotationMark, rotation [4]float64) [4]float64 {
	return quaternionToArray(
		quaternion.Prod(
			arrayToQuaternion(rotation),
			rotateQuaternion(mark),
		),
	)
}

func arrayToQuaternion(a [4]float64) quaternion.Quaternion {
	return quaternion.New(a[0], a[1], a[2], a[3])
}

func quaternionToArray(q quaternion.Quaternion) [4]float64 {
	return [4]float64{q.W, q.X, q.Y, q.Z}
}

func rotateQuaternion(mark RotationMark) quaternion.Quaternion {
	switch mark {
	case cU, ccD:
		return qU
	case cD, ccU:
		return qD
	case cF, ccB:
		return qF
	case cB, ccF:
		return qB
	case cL, ccR:
		return qL
	case cR, ccL:
		return qR
	case dcU, dcD:
		return qU2
	case dcF, dcB:
		return qF2
	case dcL, dcR:
		return qL2
	default:
		return quaternion.Quaternion{}
	}
}
