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
	stateDefinition struct {
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

	rotationDegree int
)

const (
	rotateUnknown rotationDegree = iota
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

	document         *gltf.Document
	binaryDocument   *gltf.Document
	definition       stateDefinition
	binaryDefinition stateDefinition
)

// Initialize initialize cube state
func Initialize() (err error) {
	gltf.RegisterExtension(unlit.ExtUnlit, unlit.New)

	document, err = gltf.Open(conf.Shared.Model.FilePath)
	if err != nil {
		return err
	}

	if definition, err = nodeToDefinition(document.Nodes); err != nil {
		return err
	}

	binaryDocument, err = gltf.Open(conf.Shared.Model.BinaryFilePath)
	if err != nil {
		return err
	}

	if binaryDefinition, err = nodeToDefinition(document.Nodes); err != nil {
		return err
	}

	return nil
}

func nodeToDefinition(nodes []gltf.Node) (stateDefinition, error) {
	var def stateDefinition

	if len(nodes) != 26 {
		return def, fmt.Errorf("insufficient number of nodes")
	}

	for _, node := range document.Nodes {
		switch node.Name {
		case "01_UBL":
			def.UBL = node
		case "02_UB":
			def.UB = node
		case "03_UBR":
			def.UBR = node
		case "04_UL":
			def.UL = node
		case "05_U":
			def.U = node
		case "06_UR":
			def.UR = node
		case "07_UFL":
			def.UFL = node
		case "08_UF":
			def.UF = node
		case "09_UFR":
			def.UFR = node
		case "10_BL":
			def.BL = node
		case "11_B":
			def.B = node
		case "12_BR":
			def.BR = node
		case "13_L":
			def.L = node
		case "14_R":
			def.R = node
		case "15_FL":
			def.FL = node
		case "16_F":
			def.F = node
		case "17_FR":
			def.FR = node
		case "18_DBL":
			def.DBL = node
		case "19_DB":
			def.DB = node
		case "20_DBR":
			def.DBR = node
		case "21_DL":
			def.DL = node
		case "22_D":
			def.D = node
		case "23_DR":
			def.DR = node
		case "24_DFL":
			def.DFL = node
		case "25_DF":
			def.DF = node
		case "26_DFR":
			def.DFR = node
		default:
			return def, fmt.Errorf("unexpected node name: %s", node.Name)
		}
	}
	return def, nil
}

// Generate generate cube data
func Generate(
	algorithm []string,
	isBinary bool,
	isUnlit bool,
) ([]byte, error) {
	if document == nil {
		return nil, fmt.Errorf("glTF document is not found")
	}

	var (
		doc gltf.Document
		def stateDefinition
	)

	if isBinary {
		doc = *document
		def = definition
	} else {
		doc = *binaryDocument
		def = binaryDefinition
	}

	if isUnlit {
		doc.ExtensionsUsed = []string{unlit.ExtUnlit}
		doc.ExtensionsRequired = []string{unlit.ExtUnlit}

		var materials []gltf.Material
		for _, m := range doc.Materials {
			m.Extensions = map[string]interface{}{
				unlit.ExtUnlit: struct{}{},
			}
			materials = append(materials, m)
		}
		doc.Materials = materials
	}

	for _, a := range algorithm {
		def = rotate(def, parseAlg(a))
	}

	doc.Nodes = []gltf.Node{
		def.UBL,
		def.UB,
		def.UBR,
		def.UL,
		def.U,
		def.UR,
		def.UFL,
		def.UF,
		def.UFR,
		def.BL,
		def.B,
		def.BR,
		def.L,
		def.R,
		def.FL,
		def.F,
		def.FR,
		def.DBL,
		def.DB,
		def.DBR,
		def.DL,
		def.D,
		def.DR,
		def.DFL,
		def.DF,
		def.DFR,
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

func parseAlg(alg string) rotationDegree {
	switch alg {
	case "U":
		return rotateClockwiseU
	case "D":
		return rotateClockwiseD
	case "F":
		return rotateClockwiseF
	case "B":
		return rotateClockwiseB
	case "L":
		return rotateClockwiseL
	case "R":
		return rotateClockwiseR
	case "U2":
		return rotateClockwiseU2
	case "D2":
		return rotateClockwiseD2
	case "F2":
		return rotateClockwiseF2
	case "B2":
		return rotateClockwiseB2
	case "L2":
		return rotateClockwiseL2
	case "R2":
		return rotateClockwiseR2
	case "U'":
		return rotateCounterClockwiseU
	case "D'":
		return rotateCounterClockwiseD
	case "F'":
		return rotateCounterClockwiseF
	case "B'":
		return rotateCounterClockwiseB
	case "L'":
		return rotateCounterClockwiseL
	case "R'":
		return rotateCounterClockwiseR
	default:
		return rotateUnknown
	}
}

func rotate(sd stateDefinition, rd rotationDegree) stateDefinition {
	switch rd {
	case rotateClockwiseU, rotateClockwiseU2, rotateCounterClockwiseU:
		sd.UBL, sd.UB, sd.UBR, sd.UL, sd.U, sd.UR, sd.UFL, sd.UF, sd.UFR =
			move(rd, sd.UBL, sd.UB, sd.UBR, sd.UL, sd.U, sd.UR, sd.UFL, sd.UF, sd.UFR)
	case rotateClockwiseD, rotateClockwiseD2, rotateCounterClockwiseD:
		sd.DFL, sd.DF, sd.DFR, sd.DL, sd.D, sd.DR, sd.DBL, sd.DB, sd.DBR =
			move(rd, sd.DFL, sd.DF, sd.DFR, sd.DL, sd.D, sd.DR, sd.DBL, sd.DB, sd.DBR)
	case rotateClockwiseF, rotateClockwiseF2, rotateCounterClockwiseF:
		sd.UFL, sd.UF, sd.UFR, sd.FL, sd.F, sd.FR, sd.DFL, sd.DF, sd.DFR =
			move(rd, sd.UFL, sd.UF, sd.UFR, sd.FL, sd.F, sd.FR, sd.DFL, sd.DF, sd.DFR)
	case rotateClockwiseB, rotateClockwiseB2, rotateCounterClockwiseB:
		sd.UBR, sd.UB, sd.UBL, sd.BR, sd.B, sd.BL, sd.DBR, sd.DB, sd.DBL =
			move(rd, sd.UBR, sd.UB, sd.UBL, sd.BR, sd.B, sd.BL, sd.DBR, sd.DB, sd.DBL)
	case rotateClockwiseL, rotateClockwiseL2, rotateCounterClockwiseL:
		sd.UBL, sd.UL, sd.UFL, sd.BL, sd.L, sd.FL, sd.DBL, sd.DL, sd.DFL =
			move(rd, sd.UBL, sd.UL, sd.UFL, sd.BL, sd.L, sd.FL, sd.DBL, sd.DL, sd.DFL)
	case rotateClockwiseR, rotateClockwiseR2, rotateCounterClockwiseR:
		sd.UFR, sd.UR, sd.UBR, sd.FR, sd.R, sd.BR, sd.DFR, sd.DR, sd.DBR =
			move(rd, sd.UFR, sd.UR, sd.UBR, sd.FR, sd.R, sd.BR, sd.DFR, sd.DR, sd.DBR)
	}
	return sd
}

func move(rd rotationDegree,
	n1 gltf.Node, n2 gltf.Node, n3 gltf.Node,
	n4 gltf.Node, n5 gltf.Node, n6 gltf.Node,
	n7 gltf.Node, n8 gltf.Node, n9 gltf.Node,
) (
	gltf.Node, gltf.Node, gltf.Node,
	gltf.Node, gltf.Node, gltf.Node,
	gltf.Node, gltf.Node, gltf.Node,
) {
	n1.Rotation = prod(rd, n1.RotationOrDefault())
	n2.Rotation = prod(rd, n2.RotationOrDefault())
	n3.Rotation = prod(rd, n3.RotationOrDefault())
	n4.Rotation = prod(rd, n4.RotationOrDefault())
	n5.Rotation = prod(rd, n5.RotationOrDefault())
	n6.Rotation = prod(rd, n6.RotationOrDefault())
	n7.Rotation = prod(rd, n7.RotationOrDefault())
	n8.Rotation = prod(rd, n8.RotationOrDefault())
	n9.Rotation = prod(rd, n9.RotationOrDefault())

	switch rd {
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

func prod(rd rotationDegree, rotation [4]float64) [4]float64 {
	return quaternionToArray(
		quaternion.Prod(
			arrayToQuaternion(rotation),
			rotateQuaternion(rd),
		),
	)
}

func arrayToQuaternion(a [4]float64) quaternion.Quaternion {
	return quaternion.New(a[0], a[1], a[2], a[3])
}

func quaternionToArray(q quaternion.Quaternion) [4]float64 {
	return [4]float64{q.W, q.X, q.Y, q.Z}
}

func rotateQuaternion(rd rotationDegree) quaternion.Quaternion {
	switch rd {
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
