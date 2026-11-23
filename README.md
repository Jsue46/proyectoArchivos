# Virtual File System Engine — EXT2 / EXT3 Implementation in Go

Academic project developed for the **File Management and Implementation** course at the Universidad de San Carlos de Guatemala (First Semester 2025).

This repository contains a full virtual file system implementation written in Go. The engine simulates physical disk storage, partitioning (MBR and EBR), and ext2/ext3 file system structures including superblocks, inodes, data blocks, journaling, user and permission management, and graphical report generation.

---

## Overview

The system simulates low-level binary disk management by reading and writing structured binary data to raw disk image files (`.dsk`). It exposes an execution shell that accepts commands from terminal sessions or batch input files, processing operations for formatting disks, managing directories and files, modifying permissions, and outputting Graphviz diagrams for visual inspection.

---

## Architecture and Core Modules

```
[ Command Shell / Input Script Parser ]
                  │
                  v
       [ Analyzer & Dispatcher ]
                  │
   +--------------+--------------+
   │              │              │
   v              v              v
[ Disk & MBR ] [ Inode / Superblock ] [ Permissions & Users ]
   │              │              │
   +--------------+--------------+
                  │
                  v
[ Virtual File System Manager (EXT2/EXT3) ]
  ├── Superblock allocation & state tracking
  ├── Bitmap management for Inodes and Data Blocks
  ├── Journaling mechanisms (EXT3 transaction log)
  ├── File pointers (Direct, Indirect, Double indirect)
  └── User/Group permission validation (rwx)
                  │
                  v
[ Binary Disk Image (.dsk) & Graphviz DOT Reports ]
```

---

## Key Features

- Raw binary file manipulation for simulated hard drives
- Primary, extended, and logical partition management via Master Boot Record (MBR) and Extended Boot Record (EBR)
- Native EXT2 and EXT3 file system initialization (`mkfs`)
- Inode allocation, block management (folder blocks, file blocks, pointer blocks)
- Multi-user authentication, group privileges, and UNIX-style permission checks (`chmod`, `chown`)
- Journaling module for transaction tracking and recovery
- Automated report generation using Graphviz (disk layouts, inode trees, file block maps)

---

## Technologies

- **Programming Language:** Go (Golang)
- **Visualization:** Graphviz (DOT format)
- **Interface:** CLI and REST API server with web front-end
- **Storage:** Binary `.dsk` files

---

## Repository Structure

```
proyectoArchivos/
├── Analyzer/                 # Lexical/syntactic command parsers
├── FSystem/                  # EXT2/EXT3 filesystem structures and logic
├── Partitions/               # MBR and EBR binary partition handlers
├── FPermissions/             # User and group authorization models
├── Reportes/                 # Graphviz report generators
├── server/                   # Optional HTTP API server
├── frontend/                 # Web visualization interface
├── main.go                   # CLI entry point
└── README.md
```

---

## Academic Context

| Field       | Detail                                                 |
|-------------|--------------------------------------------------------|
| Course      | File Management and Implementation (MIA)               |
| Institution | Universidad de San Carlos de Guatemala                 |
| Semester    | First Semester 2025                                    |
| Author      | Estuardo Josue Vaquiax Reyes (202201336)               |
