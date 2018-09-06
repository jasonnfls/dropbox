package main

import (
    "github.com/gorilla/mux"
    "net/http"
    "time"
    "fmt"
    "log"
    "os"
    "io"
    "io/ioutil"
    "bufio"
)

const base_path = "/volume1/Dropbox/"

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/{path:.*}", uploadHandler).Methods("POST")
    r.HandleFunc("/{path:.*}", downloadHandler).Methods("GET")
    r.HandleFunc("/{path:.*}", deleteHandler).Methods("DELETE")

    srv := &http.Server{
        Handler: r,
        Addr: "127.0.0.1:1234",
        WriteTimeout: 15 * time.Second,
        ReadTimeout: 15 * time.Second,
    }

    fmt.Printf("Listening on %s...\n", srv.Addr)
    log.Fatal(srv.ListenAndServe())
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    path := mux.Vars(r)["path"]
    fullpath := base_path + path

    file, handler, err := r.FormFile("file")
    if err == nil {
        defer file.Close()

        f, err := os.OpenFile(fullpath + handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            w.WriteHeader(500)
            w.Write([]byte(err.Error()))
            return
        }
        defer f.Close()
        io.Copy(f, file)
        http.Redirect(w, r, "", 301)
    } else {
        if err := os.Mkdir(fullpath, 0700); err != nil {
            w.WriteHeader(500)
            w.Write([]byte(err.Error()))
            return
        }
        w.WriteHeader(200)
    }
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
    path := mux.Vars(r)["path"]
    fullpath := base_path + path

    fi, err := os.Stat(fullpath)
    if err != nil {
        w.WriteHeader(404)
        w.Write([]byte(path+" not found"))
        return
    }

    if fi.IsDir() {
        if err := os.RemoveAll(fullpath); err != nil {
            w.WriteHeader(500)
            w.Write([]byte(err.Error()))
            return
        }
    } else {
        if err := os.Remove(fullpath); err != nil {
            w.WriteHeader(500)
            w.Write([]byte(err.Error()))
            return
        }
    }
    w.WriteHeader(200)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
    path := mux.Vars(r)["path"]
    fullpath := base_path + path

    fi, err := os.Stat(fullpath)
    if err != nil {
        w.WriteHeader(404)
        w.Write([]byte(path+" not found"))
        return
    }

    if fi.IsDir() {
        if len(path)>0 && path[len(path)-1] != '/' {
            http.Redirect(w, r, path+"/", 301)
            return
        }
        fis, err := ioutil.ReadDir(fullpath)
        if err != nil {
            w.WriteHeader(500)
            w.Write([]byte(err.Error()))
            return
        }
        w.WriteHeader(200)

        displayPath := "/" + path


        bw := bufio.NewWriter(w)
        bw.WriteString(fmt.Sprintf(`<!DOCTYPE html><html>
            <head>
                <meta charset=utf-8>
                <meta http-equiv=X-UA-Compatible content="IE=edge">
                <meta name=viewport content="width=device-width,initial-scale=1, user-scalable=no">
                <title>Du Cheng's File Exchanger</title>
                <script>
                function Remove(file) {
                    if (window.fetch) {
                        fetch(file, {
                            method: "DELETE",
                        }).then(function(){location.reload()});
                    }
                }
                function NewFolder() {
                    if (window.fetch) {
                        var path = document.getElementById("new_folder");
                        fetch(path.value, {
                            method: "POST",
                        }).then(function(){location.reload()});
                    }
                }
                </script>
            </head>
            <body>
                <h2>%s <a href="..">↑</a></h2>
                <hr>
                <form ENCTYPE="multipart/form-data" method="post">
                    <input name="file" type="file"/><input type="submit" value="upload"/>
                </form>
                <hr>
                <ul>`, displayPath))
        for _, f := range fis {
            if f.IsDir() {
                bw.WriteString(fmt.Sprintf(`<li style="clear:both"><a href="%s/">%s/</a> <button onclick='Remove("%s/");' style="float:right">Delete</button></li>`, f.Name(), f.Name(), f.Name()))
            } else {
                bw.WriteString(fmt.Sprintf(`<li style="clear:both"><a href="%s">%s</a> <button onclick='Remove("%s");' style="float:right">Delete</button></li>`, f.Name(), f.Name(), f.Name()))
            }
        }
        bw.WriteString(`
                </ul>
                <hr>
                <input id="new_folder" type="text"/><button onclick="NewFolder();">New Folder</button>
            </body>
        </html>`)
        bw.Flush()
    } else {
        in, err := os.Open(fullpath)
        if err != nil {
            w.WriteHeader(500)
            w.Write([]byte(err.Error()))
            return
        }
        defer in.Close()

        w.Header().Set("Content-Disposition", `attachment; filename="`+fi.Name()+`"`)
        w.Header().Set("Content-Transfer-Encoding", `binary`)
        w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

        io.Copy(w, in)
    }
}
