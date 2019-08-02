package main

import (
    "os"
    "fmt"
    "log"
    "strings"
    "net/http"
    "html/template"
    "path/filepath"
)

const (
    PageNotFoundString        = "Page Not Found"
    InternalServerErrorString = "Internal Server Error"
)

var (
    StaticResourceTypes = []string{"css", "img", "js", "fonts"}
    pageNotFoundError = fmt.Errorf("Page Not Found")
    webPageTemplatesInstance = &WebPageTemplates{}
)

func main() {
    os.Chdir(`D:\development\codeview\servicestage\src\code.huawei.com\2012_PaaS_CPE\go-templates\simplegowebapp`)

    // Start web server
    startWebServer()
}

func startWebServer() {
    // Register static resource handler
    for _, staticResourceType := range StaticResourceTypes {
        rsrPath := fmt.Sprintf("/static/%s/", staticResourceType)
        http.Handle(rsrPath, http.StripPrefix(rsrPath, http.FileServer(http.Dir(strings.TrimPrefix(rsrPath, "/")))))
    }

    // Register home handler
    http.HandleFunc("/", HomeHandler)
    // Load static pages
    RegisterWebPages("view/index.html", "view/404.html")

    // Start HTTP server on specified port
    err := http.ListenAndServe("0.0.0.0:8080", nil)
    if err != nil {
        log.Fatalf("start http server failed(err: %s).", err.Error())
    }
}

func HomeHandler(w http.ResponseWriter, req *http.Request) {
    path := req.URL.Path
    if path == "" || path == "/" {
        // uri not specified, redirect to home index page
        w.Header().Set("Location", "/index.html")
        w.WriteHeader(http.StatusMovedPermanently)
        return
    }

    // If uri ends has suffix like '.html', '.htm', treat it as web page request
    if req.Method == http.MethodGet && (strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".htm")) {
        WebPageHandler(w, req)
        return
    }

    // else, treat it as REST request
    RestHandler(w, req)
}

func RestHandler(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte(PageNotFoundString))
}

func WebPageHandler(w http.ResponseWriter, req *http.Request) {
    pageName := filepath.Base(req.URL.Path)
    tpl, err := Template(pageName)
    if err != nil {
        // redirect to page not found error page
        w.Header().Set("Location", "/404.html")
        w.WriteHeader(http.StatusMovedPermanently)
        return
    }

    params := struct {
        ServerType string
    }{"Golang Web APP"}

    // Render template with parameters
    err = tpl.Execute(w, params)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(InternalServerErrorString))
        return
    }
}

// Template cache
type WebPageTemplates struct {
    tpl *template.Template
}

func RegisterWebPages(filenames ...string) error {
    return webPageTemplatesInstance.LoadWebPageFiles(filenames...)
}

func Template(pageName string) (*template.Template, error) {
    return webPageTemplatesInstance.Template(pageName)
}

func (tpls *WebPageTemplates) LoadWebPageFiles(filenames ...string) error {
    tpl, err := template.ParseFiles(filenames...)
    if err != nil {
        log.Printf("load web pages failed(err: %s).", err.Error())
        return err
    }

    tpls.tpl = tpl
    return nil
}

func (tpls *WebPageTemplates) Template(pageName string) (*template.Template, error) {
    if tpls.tpl == nil {
        return nil, pageNotFoundError
    }

    tpl := tpls.tpl.Lookup(pageName)
    if tpl== nil {
        return nil, pageNotFoundError
    }

    return tpl, nil
}
