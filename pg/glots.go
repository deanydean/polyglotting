package pg

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path"
    "strings"
)

// SourceIndex contains stats for a certain type of source
type SourceIndex struct {
    SourceFileExts []string
    SourceType     string
    FileCount      int64
    LineCount      int64
}

// GlotIndices containing all source indices
type GlotIndices []*SourceIndex

func (gi GlotIndices) Len() int      { return len(gi) }
func (gi GlotIndices) Swap(i, j int) { gi[i], gi[j] = gi[j], gi[i] }
func (gi GlotIndices) Less(i, j int) bool {
    return gi[i].FileCount < gi[j].FileCount
}

// NewSourceIndex create a new source index
func NewSourceIndex(srcType string, exts ...string) *SourceIndex {
    srcIndex := new(SourceIndex)

    srcIndex.SourceFileExts = exts
    srcIndex.SourceType = srcType
    srcIndex.FileCount = 0
    srcIndex.LineCount = 0

    return srcIndex
}

// LinesInFile returns the number of lines in the provided filename
func LinesInFile(filename string) (int, error) {
    var reader, _ = os.Open(filename)
    var count = 0
    buf := make([]byte, 32*1024)
    lineSep := []byte{'\n'}

    for {
        c, err := reader.Read(buf)
        count += bytes.Count(buf[:c], lineSep)

        switch {
        case err == io.EOF:
            reader.Close()
            return count, nil

        case err != nil:
            reader.Close()
            return count, err
        }
    }
}

// IndexFile will add the information about the provided filePath into the
// provided index
func IndexFile(filePath string, index *SourceIndex) {
    index.FileCount++

    var lines, err = LinesInFile(filePath)

    if err != nil {
        fmt.Println("Failed to count lines in", filePath,
            "err:", err)
    }

    index.LineCount += int64(lines)
}

// IsFileType returns true if fileName has one of the provided exts, or false
// otherwise.
func IsFileType(fileName string, exts []string) bool {
    for e := range exts {
        if strings.HasSuffix(fileName, exts[e]) {
            return true
        }
    }
    return false
}

// FindInDir looks for source files in the provided dirname
func FindInDir(dirname string, glots []*SourceIndex) {
    // Read all files in dirname
    files, err := ioutil.ReadDir(dirname)

    if err != nil {
        fmt.Println("Failed to read", dirname, " e:", err)
        return
    }

    for _, file := range files {
        var filePath = path.Join(dirname, file.Name())

        if file.IsDir() {
            FindInDir(filePath, glots)
        } else {
            for i := range glots {
                var index = glots[i]
                if IsFileType(file.Name(), index.SourceFileExts) {
                    IndexFile(filePath, index)
                }
            }
        }
    }
}

// GetNewGlotsList containing sources indexes
func GetNewGlotsList() GlotIndices {
    // Create the global glot list....
    return GlotIndices{
        NewSourceIndex("Go", ".go"),
        NewSourceIndex("Javascript", ".js"),
        NewSourceIndex("Java", ".java", ".jsp"),
        NewSourceIndex("C/C++", ".c", ".cpp", ".cc", ".h"),
        NewSourceIndex("TCL", ".tcl", ".tcsh"),
        NewSourceIndex("Shell", ".sh"),
        NewSourceIndex("Expect", ".exp"),
        NewSourceIndex("Python", ".py"),
        NewSourceIndex("Make", "Makefile", ".mk", ".m"),
        NewSourceIndex("Groovy", ".groovy"),
        NewSourceIndex("Ruby", ".rb"),
        NewSourceIndex("C#", ".cs"),
        NewSourceIndex("WWW", ".htm", ".html", ".css"),
        NewSourceIndex("Rust", ".rs"),
        NewSourceIndex("PHP", ".php"),
        NewSourceIndex("Perl", ".pl"),
    }
}
