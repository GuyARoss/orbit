package webwrap

import (
	"bufio"
	context "context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var nodeProcess *os.Process

func Close() error {
	// this is a hack, node process does not get terminated with the nodeProcess.Kill
	// method, but I found that if I use ctrl^c in the terminal, it closes it correctly
	return syscall.Kill(nodeProcess.Pid, syscall.SIGSTOP)
}

func init() {
	serverStartupTasks = append(serverStartupTasks, StartupTaskReactSSR(bundleDir, wrapDocRender, staticResourceMap, make(map[PageRender]string)))
}

func StartupTaskReactSSR(outDir string, pages map[PageRender]*DocumentRenderer, staticMap map[PageRender]bool, nameMap map[PageRender]string) func() {
	return func() {
		err := startNodeServer()
		if err != nil {
			panic(err)
		}

		d := setupDoc()
		for b, r := range pages {
			if r.version != "reactSSR" || !staticMap[b] {
				continue
			}

			sr, _ := reactSSR(context.Background(), string(b), []byte("{}"), d)

			pathName := string(b)
			if nameMap[b] != "" {
				pathName = nameMap[b]
			}

			path := fmt.Sprintf("%s%c%s", http.Dir(outDir), os.PathSeparator, pathName)
			body := append(pageDependencies[b], sr.Body...)

			so := fmt.Sprintf(`<!doctype html><head>%s</head><body>%s</body></html>`, strings.Join(sr.Head, ""), strings.Join(body, ""))

			err := ioutil.WriteFile(path, []byte(so), 0644)
			if err != nil {
				fmt.Printf("error creating static resource for bundle %s => %s\n", b, err)
				continue
			}
		}
	}
}

func startNodeServer() error {
	if nodeProcess != nil {
		// TODO: already started
		return nil
	}

	cmd := exec.Command("./node_modules/.bin/babel-node", ".orbit/base/pages/react_ssr.js", "--presets", "@babel/react,@babel/preset-env")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	nodeProcess = cmd.Process
	booted := make(chan bool)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "boot success") {
				fmt.Println(line)
				booted <- true
			}

			if strings.Contains(line, "boot fail") {
				fmt.Println(line)
				booted <- false
			}
		}
	}()

	go func() {
		_, err := nodeProcess.Wait()

		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		return err
	}

	<-booted
	return nil
}

func reactSSR(ctx context.Context, bundleKey string, data []byte, doc *htmlDoc) (*htmlDoc, context.Context) {
	if nodeProcess == nil {
		fmt.Println("react ssr process has not yet boot")
		return doc, ctx
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial("0.0.0.0:3024", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := NewReactRendererClient(conn)

	response, err := client.Render(ctx, &RenderRequest{
		BundleID: bundleKey,
		JSONData: string(data),
	})

	if err != nil {
		// TODO: return error & body
		doc.Body = append(doc.Body, "<div>error loading page part of the page</div>")
	}

	doc.Body = append(doc.Body, response.StaticContent)
	return doc, ctx
}
