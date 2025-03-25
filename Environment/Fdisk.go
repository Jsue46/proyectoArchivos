package Environment

import (
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
)

func Fdisk(size int, path string, name string, unit string, type_ string, fit string) {
	fmt.Println("╔═════════════════════ INICIO  FDISK  ═════════════════════════╗")
	fmt.Println("  Size:", size)
	fmt.Println("  Path:", path)
	fmt.Println("  Name:", name)
	fmt.Println("  Unit:", unit)
	fmt.Println("  Type:", type_)
	fmt.Println("  Fit:", fit)

	if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("Error: Fit must be 'bf', 'ff', or 'wf'")
		return
	}
	if size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Unit must be 'b', 'k', or 'm'")
		return
	}
	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Println("Error: Type must be 'p', 'e', or 'l'")
		return
	}
	if name == "" {
		fmt.Println("Error: Name is required")
		return
	}
	if path == "" {
		fmt.Println("Error: Path is required")
		return
	}
	// Abrir el archivo binario en la ruta proporcionada
	file, err := Utils.AbrirArchivo(path)
	if err != nil {
		fmt.Println("Error: Could not open file at path:", path)
		return
	}
	defer file.Close()

	// Leer el MBR
	var TempMBR Partitions.MBR
	if err := Utils.LeerArchivo(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: Could not read MBR from file")
		return
	}
	// Imprimir el objeto MBR
	Partitions.ImprimirMBR(TempMBR)
	fmt.Println("─────────────────────────────────────────────────────────────────")

	// Validaciones de las particiones
	var primaryCount, extendedCount, totalPartitions int
	var usedSpace int32 = 0

	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].PartSize != 0 {
			totalPartitions++
			usedSpace += TempMBR.MbrPartitions[i].PartSize

			if TempMBR.MbrPartitions[i].PartType[0] == 'p' {
				primaryCount++
			} else if TempMBR.MbrPartitions[i].PartType[0] == 'e' {
				extendedCount++
			}
		}
	}
	// Validar que no se exceda el número máximo de particiones primarias y extendidas
	if totalPartitions >= 4 {
		fmt.Println("Error: No se pueden crear más de 4 particiones primarias o extendidas en total.")
		return
	}
	// Validar que solo haya una partición extendida
	if type_ == "e" && extendedCount > 0 {
		fmt.Println("Error: Solo se permite una partición extendida por disco.")
		return
	}
	// Validar que no se creen particiones extendidas antes de las primarias
	if type_ == "l" && extendedCount == 0 {
		fmt.Println("Error: No se puede crear una partición lógica sin una partición extendida.")
		return
	}
	// Validar que el tamaño de la nueva partición no exceda el tamaño del disco
	if usedSpace+int32(size) > TempMBR.MbrTamanio {
		fmt.Println("Error: No hay suficiente espacio en el disco para crear esta partición.")
		return
	}
	// Determinar la posición de inicio de la nueva partición
	var gap int32 = int32(binary.Size(TempMBR))
	if totalPartitions > 0 {
		gap = TempMBR.MbrPartitions[totalPartitions-1].PartStart + TempMBR.MbrPartitions[totalPartitions-1].PartSize
	}

	// Encontrar una posición vacía para la nueva partición
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].PartSize == 0 {
			if type_ == "p" || type_ == "e" {
				TempMBR.MbrPartitions[i].PartSize = int32(size)
				TempMBR.MbrPartitions[i].PartStart = gap
				copy(TempMBR.MbrPartitions[i].PartName[:], name)
				copy(TempMBR.MbrPartitions[i].PartFit[:], fit)
				copy(TempMBR.MbrPartitions[i].PartStatus[:], "0")
				copy(TempMBR.MbrPartitions[i].PartType[:], type_)
				TempMBR.MbrPartitions[i].PartCorrelative = int32(totalPartitions + 1)

				if type_ == "e" {
					ebrStart := gap
					ebr := Partitions.EBR{
						PartFit:   fit[0],
						PartStart: ebrStart,
						PartSize:  0,
						PartNext:  -1,
					}
					copy(ebr.PartName[:], "")
					Utils.EscribirArchivo(file, ebr, int64(ebrStart))
				}

				break
			}
		}
	}
	// Manejar la creación de particiones lógicas dentro de una partición extendida
	if type_ == "l" {
		for i := 0; i < 4; i++ {
			if TempMBR.MbrPartitions[i].PartType[0] == 'e' {
				ebrPos := TempMBR.MbrPartitions[i].PartStart
				var ebr Partitions.EBR
				for {
					Utils.LeerArchivo(file, &ebr, int64(ebrPos))
					if ebr.PartNext == -1 {
						break
					}
					ebrPos = ebr.PartNext
				}
				// Calcular la posición de inicio de la nueva partición lógica
				newEBRPos := ebr.PartStart + ebr.PartSize
				logicalPartitionStart := newEBRPos + int32(binary.Size(ebr))
				ebr.PartNext = newEBRPos
				Utils.EscribirArchivo(file, ebr, int64(ebrPos))
				// Crear y escribir el nuevo EBR
				newEBR := Partitions.EBR{
					PartFit:   fit[0],
					PartStart: logicalPartitionStart,
					PartSize:  int32(size),
					PartNext:  -1,
				}
				copy(newEBR.PartName[:], name)
				Utils.EscribirArchivo(file, newEBR, int64(newEBRPos))
				// Imprimir el nuevo EBR creado
				fmt.Println("Nuevo EBR creado:")
				Partitions.ImrprimirEBR(newEBR)
				break
			}
		}
	}
	// Sobrescribir el MBR
	if err := Utils.EscribirArchivo(file, TempMBR, 0); err != nil {
		fmt.Println("Error: Could not write MBR to file")
		return
	}

	var TempMBR2 Partitions.MBR
	// Leer el objeto nuevamente para verificar
	if err := Utils.LeerArchivo(file, &TempMBR2, 0); err != nil {
		fmt.Println("Error: Could not read MBR from file after writing")
		return
	}

	// Imprimir el objeto MBR actualizado
	Partitions.ImprimirMBR(TempMBR2)
	fmt.Println("╚═════════════════════   FIN   MKDISK   ═════════════════════════╝")
}
