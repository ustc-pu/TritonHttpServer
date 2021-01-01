package tritonhttp

import (
	"log"
	"os"
	"bufio"
	"strings"
)
/** 
	Load and parse the mime.types file 
**/
func ParseMIME(MIMEPath string) (MIMEMap map[string]string, err error) {
	// panic("todo - ParseMIME")
	log.Println("Reading from file " + MIMEPath)

	f, err := os.Open(MIMEPath)
	if err != nil {
		log.Panicln(err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Panicln(err)
		}
	}()

	MIMEMap = make(map[string]string)
	// winners := []string{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		content := []string{}
		content = strings.Fields(s.Text())
		MIMEMap[content[0]] = content[1]
	}
	err = s.Err()
	if err != nil {
		log.Panicln(err)
	}

	return MIMEMap, err

}

