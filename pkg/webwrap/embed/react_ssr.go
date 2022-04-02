package webwrap

import (
	"bufio"
	context "context"
	"fmt"
	"os"
	"os/exec"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var nodeProcess *os.Process

func init() {
	err := nodeServerInstance()
	if err != nil {
		panic(err)
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

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			f, _ := os.Create("./temp.txt")
			fmt.Fprintf(f, "then %s\n", line)
			f.Close()
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

	conn, err := grpc.Dial("0.0.0.0:30032", opts...)
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
		fmt.Println(err)
		// @@todo(guy): this should not be handled here
		return htmlDoc{
			Body: []string{"ssr failed to resolve"},
		}
	}

	doc.Body = append(doc.Body, response.StaticContent)

	return doc
}
