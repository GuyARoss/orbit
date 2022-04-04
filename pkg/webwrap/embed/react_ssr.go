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

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var nodeProcess *os.Process

func init() {
	serverStartupTasks = append(serverStartupTasks, reactSSRNodeServerStartup)
}

func reactSSRNodeServerStartup() {
	err := nodeServerInstance()
	if err != nil {
		panic(err)
	}

	d := setupDoc()

	for b, r := range wrapDocRender {
		if r.version == "reactSSR" && staticResourceMap[b] {
			sr := reactSSR(string(b), []byte("{}"), *d)

			path := fmt.Sprintf("%s%c%s", http.Dir(bundleDir), os.PathSeparator, b)
			body := append(wrapBody[b], sr.Body...)

			so := fmt.Sprintf(`<!doctype html><head>%s</head><body>%s</body></html>`, strings.Join(sr.Head, ""), strings.Join(body, ""))

			err := ioutil.WriteFile(path, []byte(so), 0644)
			if err != nil {
				fmt.Printf("error creating static resource for bundle %s\n", b)
				continue
			}

			fmt.Printf("successfully created static resource for %s\n", b)
		}
	}
}
func nodeServerInstance() error {
	if nodeProcess != nil {
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

			fmt.Println(line)
			if strings.Contains(line, "boot success") {
				booted <- true
			}
			if strings.Contains(line, "boot fail") {
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

func reactSSR(bundleKey string, data []byte, doc htmlDoc) htmlDoc {
	if nodeProcess == nil {
		fmt.Println("react ssr process has not yet boot")
		return htmlDoc{}
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

	response, err := client.Render(context.TODO(), &RenderRequest{
		BundleID: bundleKey,
		JSONData: string(data),
	})

	if err != nil {
		// @@todo return this instead of body
		doc.Body = append(doc.Body, "<div>error loading page part of the page</div>")
	}

	doc.Body = append(doc.Body, response.StaticContent)

	return doc
}
