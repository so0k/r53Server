package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/so0k/r53Server/version"
)

const (
	// BANNER is what is printed for help/info output.
	BANNER = `
 Server to index & view records in r53 zones.
 Version: %s
 Build: %s
`
)

var (
	provider   string
	configFile string
	interval   string

	awsAccessKey string
	awsSecretKey string

	port string

	updating bool

	vrsn bool
)

func init() {
	provider = "r53" // only supported provider
	flag.StringVar(&interval, "interval", "5m", "interval to generate new index.html's at")

	flag.StringVar(&awsAccessKey, "aws-access-key-id", "", "AWS access key")
	flag.StringVar(&awsSecretKey, "aws-secret-access-key", "", "AWS access secret")
	flag.StringVar(&configFile, "config", "config.yaml", "config file")

	flag.StringVar(&port, "p", "8080", "port for server to run on")

	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, version.VERSION, version.GITCOMMIT))
		flag.PrintDefaults()
	}

	flag.Parse()

	if vrsn {
		fmt.Printf("r53Server version %s, build %s", version.VERSION, version.GITCOMMIT)
		os.Exit(0)
	}
}

func main() {
	logrus.Infof("Reading config from %q", configFile)
	config, err := readConfig(configFile)
	if err != nil {
		logrus.Fatalf("Reading configFile %q failed: %v", configFile, err)
	}

	// create a new provider
	p, err := newProviders(provider, awsAccessKey, awsSecretKey, config)
	if err != nil {
		logrus.Fatalf("Creating new providers failed: %v", err)
	}

	// get the path to the static directory
	wd, err := os.Getwd()
	if err != nil {
		logrus.Fatalf("Getting working directory failed: %v", err)
	}
	staticDir := filepath.Join(wd, "static")

	// create the initial index
	if err := createStaticIndex(p, staticDir); err != nil {
		logrus.Fatalf("Creating initial static index failed: %v", err)
	}

	// parse the duration
	dur, err := time.ParseDuration(interval)
	if err != nil {
		logrus.Fatalf("Parsing %s as duration failed: %v", interval, err)
	}
	ticker := time.NewTicker(dur)

	go func() {
		// create more indices every X minutes based off interval
		for range ticker.C {
			if !updating {
				if err := createStaticIndex(p, staticDir); err != nil {
					logrus.Warnf("Creating static index failed: %v", err)
					updating = false
				}
			}
		}
	}()

	// create mux server
	mux := http.NewServeMux()

	// static files handler
	staticHandler := http.FileServer(http.Dir(staticDir))
	mux.Handle("/", staticHandler)

	// set up the server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	logrus.Infof("Starting server on port %q", port)
	logrus.Fatal(server.ListenAndServe())
}

type zone struct {
	Name    string
	Records []string
}

type data struct {
	LastUpdated string
	Zones       []zone
}

func createStaticIndex(providers []cloud, staticDir string) error {
	updating = true

	//logrus.Infof("Fetching records from %s", p.ZonesString())
	ctx := context.Background()
	//var cancelFn func()
	// if timeout > 0 {
	// 	ctx, cancelFn = context.WithTimeout(ctx, timeout)
	// }
	//defer cancelFn()

	var zones []zone
	for _, p := range providers {
		z, err := p.List(ctx)
		if err != nil {
			return fmt.Errorf("Listing all records in Zones failed: %v", err)
		}
		zones = append(zones, z...)
	}

	// create temp file to save template to
	logrus.Info("creating temporary file for template")
	f, err := ioutil.TempFile("", "r53Server")
	if err != nil {
		return fmt.Errorf("creating temp file failed: %v", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	// parse & execute the template
	logrus.Info("parsing and executing the template")
	templateDir := filepath.Join(staticDir, "../templates")
	lp := filepath.Join(templateDir, "layout.html")

	d := data{
		Zones:       zones,
		LastUpdated: time.Now().Local().Format(time.RFC1123),
	}
	tmpl := template.Must(template.New("").ParseFiles(lp))
	if err := tmpl.ExecuteTemplate(f, "layout", d); err != nil {
		return fmt.Errorf("execute template failed: %v", err)
	}
	f.Close()

	index := filepath.Join(staticDir, "index.html")
	logrus.Infof("renaming the temporary file %s to %s", f.Name(), index)
	if _, err := moveFile(index, f.Name()); err != nil {
		return fmt.Errorf("renaming result from %s to %s failed: %v", f.Name(), index, err)
	}
	updating = false
	return nil
}

func moveFile(dst, src string) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer df.Close()

	i, err := io.Copy(df, sf)
	if err != nil {
		return i, err
	}

	// Cleanup
	err = os.Remove(src)
	return i, err
}
