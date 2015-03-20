package coresize

import (
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Config Config
	Router *httprouter.Router
	Files  []ImageFile
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ParseFlags() {
	s.Config.ParseFlags()
}

func (s *Server) Setup() {
	// Routes
	router := httprouter.New()

	router.NotFound = s.handleNotFound
	router.PanicHandler = s.handlePanic

	router.GET("/", s.handleIndex)
	router.GET("/filenames.json", s.handleFilePaths)
	router.GET("/i/*filename", s.handleImage)

	s.Router = router

	// Pull images files
	if s.Config.PullFrom != "" {
		// TODO Pull options
	}

	// Load images from folder
	s.loadImages()
}

func (s *Server) Run() {
	log.Printf("Listening on port %d\n", s.Config.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.Config.Port), s.Router))
}

func (s *Server) SetupAndRun() {
	s.Setup()
	s.Run()
}

func (s *Server) loadImages() {
	fileInfos, err := ioutil.ReadDir(s.Config.FolderName)
	if err != nil {
		log.Println(err.Error())
		log.Println("0 files discovered")
		return
	}

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			if s.Config.Verbose {
				log.Printf("Discovered: %s\n", fileInfo.Name())
			}

			image := ImageFile{
				Path: path.Join(s.Config.FolderName, fileInfo.Name()),
			}

			// If configured compute checksum
			if s.Config.Hash {
				err = image.ComputeHash()
				if err != nil {
					log.Printf("Error calculating checksum for file %s (%s)", image.Name(), err.Error())
					continue
				}
			}

			s.Files = append(s.Files, image)
		}
	}

	log.Printf("%d files discovered\n", len(s.Files))
}