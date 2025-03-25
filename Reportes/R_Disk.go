package Reportes

import (
	"Proyecto/Partitions"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
