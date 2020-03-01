package cube

import (
	"bytes"
	"fmt"
	"math"

	"github.com/biohuns/blendcube/config"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/ext/unlit"
	"github.com/westphae/quaternion"
)

type (
	Definition struct {
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

	Degree int

	Rotation [4]float64
)

const (
	rotateUnknown Degree = iota
	rotateClockwiseU
	rotateClockwiseD
	rotateClockwiseF
	rotateClockwiseB
	rotateClockwiseL
	rotateClockwiseR
	rotateClockwiseU2
	rotateClockwiseD2
	rotateClockwiseF2
	rotateClockwiseB2
	rotateClockwiseL2
	rotateClockwiseR2
	rotateCounterClockwiseU
	rotateCounterClockwiseD
	rotateCounterClockwiseF
	rotateCounterClockwiseB
	rotateCounterClockwiseL
	rotateCounterClockwiseR

	degree90  = math.Pi / 2
	degree180 = math.Pi
)

var (
	quaternionUnknown = quaternion.FromEuler(0, 0, 0)
	quaternionU       = quaternion.FromEuler(0, degree90, 0)
	quaternionD       = quaternion.FromEuler(0, -degree90, 0)
	quaternionF       = quaternion.FromEuler(-degree90, 0, 0)
	quaternionB       = quaternion.FromEuler(degree90, 0, 0)
	quaternionL       = quaternion.FromEuler(0, 0, -degree90)
	quaternionR       = quaternion.FromEuler(0, 0, degree90)
	quaternionU2      = quaternion.FromEuler(0, degree180, 0)
	quaternionF2      = quaternion.FromEuler(-degree180, 0, 0)
	quaternionL2      = quaternion.FromEuler(0, 0, -degree180)

	document      *gltf.Document
	binDocument   *gltf.Document
	definition    *Definition
	binDefinition *Definition
)

func Initialize() (err error) {
	gltf.RegisterExtension(unlit.ExtensionName, unlit.Unmarshal)

	document, err = gltf.Open(config.Shared.Model.FilePath)
	if err != nil {
		return err
	}

	if definition, err = nodeToDefinition(document.Nodes); err != nil {
		return err
	}

	binDocument, err = gltf.Open(config.Shared.Model.BinaryFilePath)
	if err != nil {
		return err
	}

	if binDefinition, err = nodeToDefinition(document.Nodes); err != nil {
		return err
	}

	return nil
}

func Generate(algorithm []string, isBinary bool, isUnlit bool) ([]byte, error) {
	var (
		doc gltf.Document
		def Definition
	)

	if isBinary {
		doc = *binDocument
		def = *binDefinition
	} else {
		doc = *document
		def = *definition
	}

	if isUnlit {
		doc.ExtensionsUsed = []string{unlit.ExtensionName}
		doc.ExtensionsRequired = []string{unlit.ExtensionName}

		materials := make([]*gltf.Material, 0)
		for _, m := range doc.Materials {
			m.Extensions = map[string]interface{}{
				unlit.ExtensionName: struct{}{},
			}
			materials = append(materials, m)
		}
		doc.Materials = materials
	}

	for _, deg := range parseAlg(algorithm) {
		def = rotate(def, deg)
	}

	doc.Nodes = []*gltf.Node{
		&def.UBL, &def.UB, &def.UBR,
		&def.UL, &def.U, &def.UR,
		&def.UFL, &def.UF, &def.UFR,
		&def.BL, &def.B, &def.BR,
		&def.L, &def.R,
		&def.FL, &def.F, &def.FR,
		&def.DBL, &def.DB, &def.DBR,
		&def.DL, &def.D, &def.DR,
		&def.DFL, &def.DF, &def.DFR,
	}

	buffer := new(bytes.Buffer)
	e := gltf.NewEncoder(buffer)
	e.AsBinary = isBinary
	if err := e.Encode(&doc); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

func nodeToDefinition(nodes []*gltf.Node) (*Definition, error) {
	def := new(Definition)

	if len(nodes) != 26 {
		return nil, fmt.Errorf("insufficient number of nodes")
	}

	for _, node := range nodes {
		switch node.Name {
		case "01_UBL":
			def.UBL = *node
		case "02_UB":
			def.UB = *node
		case "03_UBR":
			def.UBR = *node
		case "04_UL":
			def.UL = *node
		case "05_U":
			def.U = *node
		case "06_UR":
			def.UR = *node
		case "07_UFL":
			def.UFL = *node
		case "08_UF":
			def.UF = *node
		case "09_UFR":
			def.UFR = *node
		case "10_BL":
			def.BL = *node
		case "11_B":
			def.B = *node
		case "12_BR":
			def.BR = *node
		case "13_L":
			def.L = *node
		case "14_R":
			def.R = *node
		case "15_FL":
			def.FL = *node
		case "16_F":
			def.F = *node
		case "17_FR":
			def.FR = *node
		case "18_DBL":
			def.DBL = *node
		case "19_DB":
			def.DB = *node
		case "20_DBR":
			def.DBR = *node
		case "21_DL":
			def.DL = *node
		case "22_D":
			def.D = *node
		case "23_DR":
			def.DR = *node
		case "24_DFL":
			def.DFL = *node
		case "25_DF":
			def.DF = *node
		case "26_DFR":
			def.DFR = *node
		default:
			return nil, fmt.Errorf("unexpected node name: %s", node.Name)
		}
	}
	return def, nil
}

func parseAlg(alg []string) []Degree {
	var deg []Degree
	for _, a := range alg {
		switch a {
		case "U":
			deg = append(deg, rotateClockwiseU)
		case "D":
			deg = append(deg, rotateClockwiseD)
		case "F":
			deg = append(deg, rotateClockwiseF)
		case "B":
			deg = append(deg, rotateClockwiseB)
		case "L":
			deg = append(deg, rotateClockwiseL)
		case "R":
			deg = append(deg, rotateClockwiseR)
		case "U2":
			deg = append(deg, rotateClockwiseU2)
		case "D2":
			deg = append(deg, rotateClockwiseD2)
		case "F2":
			deg = append(deg, rotateClockwiseF2)
		case "B2":
			deg = append(deg, rotateClockwiseB2)
		case "L2":
			deg = append(deg, rotateClockwiseL2)
		case "R2":
			deg = append(deg, rotateClockwiseR2)
		case "U'":
			deg = append(deg, rotateCounterClockwiseU)
		case "D'":
			deg = append(deg, rotateCounterClockwiseD)
		case "F'":
			deg = append(deg, rotateCounterClockwiseF)
		case "B'":
			deg = append(deg, rotateCounterClockwiseB)
		case "L'":
			deg = append(deg, rotateCounterClockwiseL)
		case "R'":
			deg = append(deg, rotateCounterClockwiseR)
		default:
			deg = append(deg, rotateUnknown)
		}
	}
	return deg
}

func rotate(def Definition, deg Degree) Definition {
	switch deg {
	case rotateClockwiseU, rotateClockwiseU2, rotateCounterClockwiseU:
		def.UBL, def.UB, def.UBR, def.UL, def.U, def.UR, def.UFL, def.UF, def.UFR =
			move(deg, def.UBL, def.UB, def.UBR, def.UL, def.U, def.UR, def.UFL, def.UF, def.UFR)
	case rotateClockwiseD, rotateClockwiseD2, rotateCounterClockwiseD:
		def.DFL, def.DF, def.DFR, def.DL, def.D, def.DR, def.DBL, def.DB, def.DBR =
			move(deg, def.DFL, def.DF, def.DFR, def.DL, def.D, def.DR, def.DBL, def.DB, def.DBR)
	case rotateClockwiseF, rotateClockwiseF2, rotateCounterClockwiseF:
		def.UFL, def.UF, def.UFR, def.FL, def.F, def.FR, def.DFL, def.DF, def.DFR =
			move(deg, def.UFL, def.UF, def.UFR, def.FL, def.F, def.FR, def.DFL, def.DF, def.DFR)
	case rotateClockwiseB, rotateClockwiseB2, rotateCounterClockwiseB:
		def.UBR, def.UB, def.UBL, def.BR, def.B, def.BL, def.DBR, def.DB, def.DBL =
			move(deg, def.UBR, def.UB, def.UBL, def.BR, def.B, def.BL, def.DBR, def.DB, def.DBL)
	case rotateClockwiseL, rotateClockwiseL2, rotateCounterClockwiseL:
		def.UBL, def.UL, def.UFL, def.BL, def.L, def.FL, def.DBL, def.DL, def.DFL =
			move(deg, def.UBL, def.UL, def.UFL, def.BL, def.L, def.FL, def.DBL, def.DL, def.DFL)
	case rotateClockwiseR, rotateClockwiseR2, rotateCounterClockwiseR:
		def.UFR, def.UR, def.UBR, def.FR, def.R, def.BR, def.DFR, def.DR, def.DBR =
			move(deg, def.UFR, def.UR, def.UBR, def.FR, def.R, def.BR, def.DFR, def.DR, def.DBR)
	}
	return def
}

func move(deg Degree,
	n1 gltf.Node, n2 gltf.Node, n3 gltf.Node,
	n4 gltf.Node, n5 gltf.Node, n6 gltf.Node,
	n7 gltf.Node, n8 gltf.Node, n9 gltf.Node,
) (
	gltf.Node, gltf.Node, gltf.Node,
	gltf.Node, gltf.Node, gltf.Node,
	gltf.Node, gltf.Node, gltf.Node,
) {
	n1.Rotation = prod(deg, n1.RotationOrDefault())
	n2.Rotation = prod(deg, n2.RotationOrDefault())
	n3.Rotation = prod(deg, n3.RotationOrDefault())
	n4.Rotation = prod(deg, n4.RotationOrDefault())
	n5.Rotation = prod(deg, n5.RotationOrDefault())
	n6.Rotation = prod(deg, n6.RotationOrDefault())
	n7.Rotation = prod(deg, n7.RotationOrDefault())
	n8.Rotation = prod(deg, n8.RotationOrDefault())
	n9.Rotation = prod(deg, n9.RotationOrDefault())

	switch deg {
	case rotateClockwiseU, rotateClockwiseD, rotateClockwiseF,
		rotateClockwiseB, rotateClockwiseL, rotateClockwiseR:
		return n7, n4, n1, n8, n5, n2, n9, n6, n3
	case rotateClockwiseU2, rotateClockwiseD2, rotateClockwiseF2,
		rotateClockwiseB2, rotateClockwiseL2, rotateClockwiseR2:
		return n9, n8, n7, n6, n5, n4, n3, n2, n1
	case rotateCounterClockwiseU, rotateCounterClockwiseD, rotateCounterClockwiseF,
		rotateCounterClockwiseB, rotateCounterClockwiseL, rotateCounterClockwiseR:
		return n3, n6, n9, n2, n5, n8, n1, n4, n7
	default:
		return n1, n2, n3, n4, n5, n6, n7, n8, n9
	}
}

func prod(deg Degree, rot Rotation) Rotation {
	return quaternionToRotation(
		quaternion.Prod(
			rotationToQuaternion(rot),
			getRotateQuaternion(deg),
		),
	)
}

func rotationToQuaternion(rot Rotation) quaternion.Quaternion {
	return quaternion.New(rot[0], rot[1], rot[2], rot[3])
}

func quaternionToRotation(q quaternion.Quaternion) Rotation {
	return Rotation{q.W, q.X, q.Y, q.Z}
}

func getRotateQuaternion(deg Degree) quaternion.Quaternion {
	switch deg {
	case rotateClockwiseU, rotateCounterClockwiseD:
		return quaternionU
	case rotateClockwiseD, rotateCounterClockwiseU:
		return quaternionD
	case rotateClockwiseF, rotateCounterClockwiseB:
		return quaternionF
	case rotateClockwiseB, rotateCounterClockwiseF:
		return quaternionB
	case rotateClockwiseL, rotateCounterClockwiseR:
		return quaternionL
	case rotateClockwiseR, rotateCounterClockwiseL:
		return quaternionR
	case rotateClockwiseU2, rotateClockwiseD2:
		return quaternionU2
	case rotateClockwiseF2, rotateClockwiseB2:
		return quaternionF2
	case rotateClockwiseL2, rotateClockwiseR2:
		return quaternionL2
	default:
		return quaternionUnknown
	}
}
