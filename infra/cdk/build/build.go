package build

import (
	"archive/zip"
	"context"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var workingDir string

var FoldersTree []string

var lambdaEnvs = []string{
	"GOOS=linux",
	"GOARCH=arm64",
	"CGO_ENABLED=0",
}

type buildEntry struct {
	mainPath string
	outPath  string
}

func init() {
	var err error
	workingDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

// Exec builds binaries for all main.go files in the cmd folder git the given env variables.
func Exec() error {
	path := filepath.Join(workingDir, "..", "cmd")
	fileList, err := getMainFiles(path)
	if err != nil {
		return err
	}

	entries, err := getBuildEntries(fileList)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())
	group.SetLimit(5)

	for _, entry := range entries {
		current := entry

		group.Go(func() error {
			err = compileEntry(current, lambdaEnvs)
			if err != nil {
				return err
			}
			return compressBin(current)
		})
	}

	return group.Wait()
}
func compileEntry(entry buildEntry, env []string) error {
	println("Compiling: ", entry.mainPath)

	if env == nil {
		env = make([]string, 0)
	}

	cmd := "go"
	args := []string{"build", "-tags", "local", "-ldflags", "-s -w", "-o", entry.outPath, entry.mainPath}

	command := exec.Command(cmd, args...)
	command.Env = append(os.Environ(), env...)
	if err := command.Run(); err != nil {
		return err
	}
	println("Compiled: ", entry.outPath)
	return nil
}

func compressBin(entry buildEntry) error {
	zipPath := entry.outPath + ".zip"
	zipFile, err := os.Create(zipPath) // crea el archivo vacio .zip
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile) // escritor de zip
	defer zipWriter.Close()

	binary, err := os.Open(entry.outPath) // se abre el binario
	if err != nil {
		return err
	}
	defer binary.Close()

	info, err := binary.Stat() // Obtener información del archivo
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info) // Crear un encabezado para el archivo ZIP con los metadatos correctos
	if err != nil {
		return err
	}

	header.Name = filepath.Base(entry.outPath) // Asegurarse de que el nombre del archivo en el ZIP sea solo el nombre base

	header.Method = zip.Deflate // Establecer el metodo de compresión

	// Crear la entrada en el ZIP con el encabezado que incluye los metadatos
	compressedEntry, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Copiar el contenido del binario al archivo comprimido
	_, err = io.Copy(compressedEntry, binary)
	if err != nil {
		return err
	}

	return nil
}

func getMainFiles(path string) ([]string, error) {
	omit := []string{
		"cdk",
		"config",
		"node_modules",
		".serverless",
		".bin",
	}
	return Search(path, "main.go", WithOmit(omit...))
}

func getBuildEntries(fileList []string) ([]buildEntry, error) {
	entries := make([]buildEntry, 0, len(fileList))

	for _, file := range fileList {
		entries = append(entries, buildEntry{
			mainPath: file,
			outPath:  getOutPath(file),
		})
	}

	return entries, nil
}

func getOutPath(file string) string {
	sections := strings.Split(file, "cmd/")
	folderTree := strings.ReplaceAll(sections[len(sections)-1], "/main.go", "")
	FoldersTree = append(FoldersTree, folderTree)
	return filepath.Join(workingDir, ".bin", folderTree, "bootstrap")
}
