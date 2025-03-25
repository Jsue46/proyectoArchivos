package Reportes

import (
	"Proyecto/Partitions"
	"Proyecto/Utils"
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
					err := Utils.LeerArchivo(file, &ebr, int64(ebrPosition)) // Leer EBR desde el archivo
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
