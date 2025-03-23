package Utils

import (
	"Proyecto/Partitions"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateMBRReport genera un reporte del MBR y los EBRs en formato Graphviz y lo guarda en un archivo .dot
func GenerateMBRReport(mbr Partitions.MBR, ebrs []Partitions.EBR, outputPath string, file *os.File) error {
	// Crear la carpeta de reportes si no existe
	reportsDir := filepath.Dir(outputPath)
	err := os.MkdirAll(reportsDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error al crear la carpeta de reportes: %v", err)
	}

	// Crear el archivo .dot donde se generará el reporte
	dotFilePath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".dot"
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		return fmt.Errorf("error al crear el archivo .dot de reporte: %v", err)
	}
	defer dotFile.Close()

	// Iniciar el contenido del archivo en formato Graphviz (.dot)
	graphContent := "digraph G {\n"
	graphContent += "\tnode [fillcolor=lightyellow style=filled]\n"

	// Subgrafo del MBR con sus propiedades
	graphContent += fmt.Sprintf("\tsubgraph cluster_MBR {\n\t\tcolor=lightgrey fillcolor=lightblue label=\"MBR\nTamaño: %d\nFecha Creación: %s\nDisk Signature: %d\" style=filled\n",
		mbr.MbrTamanio, string(mbr.MbrFechaCreacion[:]), mbr.MbrDskSignature)

	lastPartitionID := ""
	// Iterar sobre las 4 particiones del MBR
	for i := 0; i < 4; i++ {
		partition := mbr.MbrPartitions[i]
		if partition.PartSize > 0 { // Si la partición tiene un tamaño válido
			partitionName := strings.TrimRight(string(partition.PartName[:]), "\x00") // Obtener nombre sin caracteres nulos
			partitionID := fmt.Sprintf("PART%d", i+1)

			// Agregar información de la partición en formato Graphviz
			graphContent += fmt.Sprintf("\t\t%s [label=\"Partición %d\nStatus: %s\nType: %s\nFit: %s\nStart: %d\nSize: %d\nName: %s\" fillcolor=green shape=box style=filled]\n",
				partitionID, i+1, string(partition.PartStatus[:]), string(partition.PartType[:]), string(partition.PartFit[:]), partition.PartStart, partition.PartSize, partitionName)

			// Conectar la partición actual con la anterior para mantener el orden
			if lastPartitionID != "" {
				graphContent += fmt.Sprintf("\t\t%s -> %s [style=invis]\n", lastPartitionID, partitionID)
			}
			lastPartitionID = partitionID

			// Manejar particiones extendidas y sus EBRs
			if string(partition.PartType[:]) == "e" {
				graphContent += fmt.Sprintf("\tsubgraph cluster_EBR%d {\n\t\tcolor=black fillcolor=lightpink label=\"Partición Extendida %d\" style=dashed\n", i+1, i+1)

				// Leer los EBRs en la partición extendida
				ebrPosition := partition.PartStart
				var ebrList []Partitions.EBR
				for {
					var ebr Partitions.EBR
					err := LeerArchivo(file, &ebr, int64(ebrPosition)) // Leer EBR desde el archivo
					if err != nil {
						break
					}
					ebrList = append(ebrList, ebr)

					// Verificar si hay más EBRs
					if ebr.PartNext == -1 {
						break
					}
					ebrPosition = ebr.PartNext
				}

				lastEbrID := ""
				// Agregar EBRs al reporte
				for j, ebr := range ebrList {
					ebrName := strings.TrimRight(string(ebr.PartName[:]), "\x00")
					ebrID := fmt.Sprintf("EBR%d", j+1)
					graphContent += fmt.Sprintf("\t\t%s [label=\"EBR\nStart: %d\nSize: %d\nNext: %d\nName: %s\" fillcolor=lightpink shape=box style=filled]\n",
						ebrID, ebr.PartStart, ebr.PartSize, ebr.PartNext, ebrName)
					if lastEbrID != "" {
						graphContent += fmt.Sprintf("\t\t%s -> %s [style=invis]\n", lastEbrID, ebrID)
					}
					lastEbrID = ebrID
				}
				graphContent += "\t}\n" // Cerrar subgrafo de EBRs
			}
		}
	}
	graphContent += "\t}\n" // Cerrar el subgrafo del MBR
	graphContent += "}\n"   // Cerrar el grafo principal

	// Escribir el contenido en el archivo .dot
	_, err = dotFile.WriteString(graphContent)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo .dot: %v", err)
	}

	return nil
}

func GenerateEBRReport(mbr Partitions.MBR, outputPath string, file *os.File) error {
	reportsDir := filepath.Dir(outputPath)
	err := os.MkdirAll(reportsDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error al crear la carpeta de reportes: %v", err)
	}

	dotFilePath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".dot"
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		return fmt.Errorf("error al crear el archivo .dot de reporte: %v", err)
	}
	defer dotFile.Close()

	graphContent := "digraph G {\n"
	graphContent += "\tnode [fillcolor=lightblue style=filled]\n"

	for i := 0; i < 4; i++ {
		partition := mbr.MbrPartitions[i]
		if string(partition.PartType[:]) == "e" {
			ebrPosition := partition.PartStart
			var ebrList []Partitions.EBR
			for {
				var ebr Partitions.EBR
				err := LeerArchivo(file, &ebr, int64(ebrPosition))
				if err != nil {
					break
				}
				ebrList = append(ebrList, ebr)

				if ebr.PartNext == -1 {
					break
				}
				ebrPosition = ebr.PartNext
			}

			for j, ebr := range ebrList {
				ebrName := strings.TrimRight(string(ebr.PartName[:]), "\x00")
				ebrID := fmt.Sprintf("EBR%d", j+1)
				graphContent += fmt.Sprintf("\t%s [label=\"EBR\nStart: %d\nSize: %d\nNext: %d\nName: %s\" shape=box fillcolor=lightpink]\n",
					ebrID, ebr.PartStart, ebr.PartSize, ebr.PartNext, ebrName)
			}
		}
	}

	graphContent += "}\n"

	_, err = dotFile.WriteString(graphContent)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo .dot: %v", err)
	}

	return nil
}

// Funcion para generar el reporte DISK en formato .dot
func GenerateDiskReport(mbr Partitions.MBR, ebrs []Partitions.EBR, outputPath string, file *os.File, totalDiskSize int32) error {
	// Obtener la ruta del directorio donde se guardará el reporte
	reportsDir := filepath.Dir(outputPath)
	// Crear la carpeta si no existe
	err := os.MkdirAll(reportsDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error al crear la carpeta de reportes: %v", err)
	}

	// Crear el archivo .dot donde se generará el reporte
	dotFilePath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".dot"
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		return fmt.Errorf("error al crear el archivo .dot de reporte: %v", err)
	}
	defer dotFile.Close()

	// Iniciar el contenido del archivo en formato Graphviz (.dot)
	graphContent := "digraph G {\n"
	graphContent += "\tnode [shape=none];\n"
	graphContent += "\tgraph [splines=false];\n"
	graphContent += "\tsubgraph cluster_disk {\n"
	graphContent += "\t\tlabel=\"Disco1.dsk\";\n"
	graphContent += "\t\tstyle=rounded;\n"
	graphContent += "\t\tcolor=black;\n"

	// Iniciar tabla para las particiones
	graphContent += "\t\ttable [label=<\n\t\t\t<TABLE BORDER=\"1\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"10\">\n"
	graphContent += "\t\t\t<TR>\n"
	graphContent += "\t\t\t<TD>MBR (159 bytes)</TD>\n"

	// Variables para el porcentaje y espacio libre
	var usedSpace int32 = 159 // Tamaño del MBR en bytes
	var freeSpace int32 = totalDiskSize - usedSpace

	// Iterar sobre las 4 particiones posibles del MBR
	for i := 0; i < 4; i++ {
		partition := mbr.MbrPartitions[i]
		if partition.PartSize > 0 { // Si la partición tiene un tamaño válido
			percentage := float64(partition.PartSize) / float64(totalDiskSize) * 100
			partitionName := strings.TrimRight(string(partition.PartName[:]), "\x00") // Limpiar el nombre de la partición

			if string(partition.PartType[:]) == "p" { // Partición primaria
				graphContent += fmt.Sprintf("\t\t\t<TD>Primaria<br/>%s<br/>%.2f%% del disco</TD>\n", partitionName, percentage)
				usedSpace += partition.PartSize
			} else if string(partition.PartType[:]) == "e" { // Partición extendida
				graphContent += "\t\t\t<TD>\n"
				graphContent += "\t\t\t\t<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"
				graphContent += "\t\t\t\t<TR><TD COLSPAN=\"5\">Extendida</TD></TR>\n"

				// Leer los EBRs y agregar las particiones lógicas
				graphContent += "\t\t\t\t<TR>\n"
				for _, ebr := range ebrs {
					logicalPercentage := float64(ebr.PartSize) / float64(totalDiskSize) * 100
					graphContent += fmt.Sprintf("\t\t\t\t<TD>EBR (32 bytes)</TD>\n\t\t\t\t<TD>Lógica<br/>%.2f%% del disco</TD>\n", logicalPercentage)
					usedSpace += ebr.PartSize + 32 // Añadir el tamaño de la partición lógica y el EBR
				}
				graphContent += "\t\t\t\t</TR>\n"
				graphContent += "\t\t\t\t</TABLE>\n"
				graphContent += "\t\t\t</TD>\n"
			}
		}
	}

	// Recalcular el espacio libre restante
	freeSpace = totalDiskSize - usedSpace
	freePercentage := float64(freeSpace) / float64(totalDiskSize) * 100

	// Agregar el espacio libre al reporte
	graphContent += fmt.Sprintf("\t\t\t<TD>Libre<br/>%.2f%% del disco</TD>\n", freePercentage)
	graphContent += "\t\t\t</TR>\n"
	graphContent += "\t\t\t</TABLE>\n>];\n"
	graphContent += "\t}\n"
	graphContent += "}\n"

	// Escribir el contenido en el archivo .dot
	_, err = dotFile.WriteString(graphContent)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo .dot: %v", err)
	}

	return nil
}
