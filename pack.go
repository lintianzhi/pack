package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	bin    = kingpin.Flag("bin", "bin file").Short('b').Required().String()
	config = kingpin.Flag("config", "config file").Short('c').Required().String()
	flag   = kingpin.Flag("flag", "flag").Short('f').String()
)

func main() {

	kingpin.Parse()

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "#!/bin/bash\n\n")
	fmt.Fprintf(&buf, "bin=%s\nconfig=%s\nflag=%s\n\n", filepath.Base(*bin), filepath.Base(*config), *flag)
	fmt.Fprintf(&buf, `dir=pack$RANDOM

function extract_bin()
{
    match=$(grep --text --line-number '^PAYLOAD_BIN:$' $0 | cut -d ':' -f 1)
    matchend=$(grep --text --line-number '^PAYLOAD_BINEND:$' $0 | cut -d ':' -f 1)
    awk "NR>$match && NR<$matchend" $0 | base64 -d > $dir/$bin
    chmod `+"`"+`stat -c "%%a" $0`+"`"+` $dir/$bin
}


function extract_config()
{
    match=$(grep --text --line-number '^PAYLOAD_CONFIG:$' $0 | cut -d ':' -f 1)
    matchend=$(grep --text --line-number '^PAYLOAD_CONFIGEND:$' $0 | cut -d ':' -f 1)
    awk "NR>$match && NR<$matchend" $0 > $dir/$config;
}

# 提供修改配置的接口

mkdir $dir | exit 1;

extract_bin;
extract_config;
$dir/$bin ${flag} $dir/$config

rm -r $dir;
exit 0


`)

	binData, err := ioutil.ReadFile(*bin)
	if err != nil {
		fmt.Errorf("err read: %s %v", *bin, err)
		return
	}

	fmt.Fprintf(&buf, "PAYLOAD_BIN:\n%s\nPAYLOAD_BINEND:\n", base64.StdEncoding.EncodeToString(binData))

	configData, err := ioutil.ReadFile(*config)
	if err != nil {
		fmt.Errorf("err read: %s %v", *config, err)
		return
	}

	fmt.Fprintf(&buf, "PAYLOAD_CONFIG:\n%s\nPAYLOAD_CONFIGEND:\n", string(configData))

	buf.WriteTo(os.Stdout)
	// print bin
}
