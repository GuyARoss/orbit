package orbitgen

import (
	"context"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"reflect"
	"sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)
// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.6.1
// source: src/com.proto
const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)
type RenderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	BundleID string `protobuf:"bytes,1,opt,name=BundleID,proto3" json:"BundleID,omitempty"`
	JSONData string `protobuf:"bytes,2,opt,name=JSONData,proto3" json:"JSONData,omitempty"`
}
func (x *RenderRequest) Reset() {
	*x = RenderRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_com_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}
func (x *RenderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*RenderRequest) ProtoMessage() {}
func (x *RenderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_src_com_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
// Deprecated: Use RenderRequest.ProtoReflect.Descriptor instead.
func (*RenderRequest) Descriptor() ([]byte, []int) {
	return file_src_com_proto_rawDescGZIP(), []int{0}
}
func (x *RenderRequest) GetBundleID() string {
	if x != nil {
		return x.BundleID
	}
	return ""
}
func (x *RenderRequest) GetJSONData() string {
	if x != nil {
		return x.JSONData
	}
	return ""
}
type RenderResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	StaticContent string `protobuf:"bytes,1,opt,name=StaticContent,proto3" json:"StaticContent,omitempty"`
}
func (x *RenderResponse) Reset() {
	*x = RenderResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_com_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}
func (x *RenderResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*RenderResponse) ProtoMessage() {}
func (x *RenderResponse) ProtoReflect() protoreflect.Message {
	mi := &file_src_com_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
// Deprecated: Use RenderResponse.ProtoReflect.Descriptor instead.
func (*RenderResponse) Descriptor() ([]byte, []int) {
	return file_src_com_proto_rawDescGZIP(), []int{1}
}
func (x *RenderResponse) GetStaticContent() string {
	if x != nil {
		return x.StaticContent
	}
	return ""
}
var File_src_com_proto protoreflect.FileDescriptor
var file_src_com_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x72, 0x63, 0x2f, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x04, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x47, 0x0a, 0x0d, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65,
	0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x4a, 0x53, 0x4f, 0x4e, 0x44, 0x61, 0x74, 0x61, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4a, 0x53, 0x4f, 0x4e, 0x44, 0x61, 0x74, 0x61, 0x22, 0x36,
	0x0a, 0x0e, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x24, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x32, 0x46, 0x0a, 0x0d, 0x52, 0x65, 0x61, 0x63, 0x74, 0x52,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x12, 0x35, 0x0a, 0x06, 0x52, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x12, 0x13, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x52, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0c,
	0x5a, 0x0a, 0x2e, 0x2f, 0x73, 0x72, 0x63, 0x3b, 0x6d, 0x61, 0x69, 0x6e, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}
var (
	file_src_com_proto_rawDescOnce sync.Once
	file_src_com_proto_rawDescData = file_src_com_proto_rawDesc
)
func file_src_com_proto_rawDescGZIP() []byte {
	file_src_com_proto_rawDescOnce.Do(func() {
		file_src_com_proto_rawDescData = protoimpl.X.CompressGZIP(file_src_com_proto_rawDescData)
	})
	return file_src_com_proto_rawDescData
}
var file_src_com_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_src_com_proto_goTypes = []interface{}{
	(*RenderRequest)(nil),  // 0: main.RenderRequest
	(*RenderResponse)(nil), // 1: main.RenderResponse
}
var file_src_com_proto_depIdxs = []int32{
	0, // 0: main.ReactRenderer.Render:input_type -> main.RenderRequest
	1, // 1: main.ReactRenderer.Render:output_type -> main.RenderResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}
func init() { file_src_com_proto_init() }
func file_src_com_proto_init() {
	if File_src_com_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_src_com_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RenderRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_src_com_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RenderResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_src_com_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_src_com_proto_goTypes,
		DependencyIndexes: file_src_com_proto_depIdxs,
		MessageInfos:      file_src_com_proto_msgTypes,
	}.Build()
	File_src_com_proto = out.File
	file_src_com_proto_rawDesc = nil
	file_src_com_proto_goTypes = nil
	file_src_com_proto_depIdxs = nil
}
// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// This is a compile-time assertion to ensure that this generated file
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7
// ReactRendererClient is the client API for ReactRenderer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReactRendererClient interface {
	Render(ctx context.Context, in *RenderRequest, opts ...grpc.CallOption) (*RenderResponse, error)
}
type reactRendererClient struct {
	cc grpc.ClientConnInterface
}
func NewReactRendererClient(cc grpc.ClientConnInterface) ReactRendererClient {
	return &reactRendererClient{cc}
}
func (c *reactRendererClient) Render(ctx context.Context, in *RenderRequest, opts ...grpc.CallOption) (*RenderResponse, error) {
	out := new(RenderResponse)
	err := c.cc.Invoke(ctx, "/main.ReactRenderer/Render", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
// ReactRendererServer is the server API for ReactRenderer service.
// All implementations must embed UnimplementedReactRendererServer
// for forward compatibility
type ReactRendererServer interface {
	Render(context.Context, *RenderRequest) (*RenderResponse, error)
	mustEmbedUnimplementedReactRendererServer()
}
// UnimplementedReactRendererServer must be embedded to have forward compatible implementations.
type UnimplementedReactRendererServer struct {
}
func (UnimplementedReactRendererServer) Render(context.Context, *RenderRequest) (*RenderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Render not implemented")
}
func (UnimplementedReactRendererServer) mustEmbedUnimplementedReactRendererServer() {}
// UnsafeReactRendererServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReactRendererServer will
// result in compilation errors.
type UnsafeReactRendererServer interface {
	mustEmbedUnimplementedReactRendererServer()
}
func RegisterReactRendererServer(s grpc.ServiceRegistrar, srv ReactRendererServer) {
	s.RegisterService(&ReactRenderer_ServiceDesc, srv)
}
func _ReactRenderer_Render_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReactRendererServer).Render(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.ReactRenderer/Render",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReactRendererServer).Render(ctx, req.(*RenderRequest))
	}
	return interceptor(ctx, in, info, handler)
}
// ReactRenderer_ServiceDesc is the grpc.ServiceDesc for ReactRenderer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReactRenderer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.ReactRenderer",
	HandlerType: (*ReactRendererServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Render",
			Handler:    _ReactRenderer_Render_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "src/com.proto",
}
func serverRenderInnerHTML(bundleKey string, data []byte) string {
	if nodeProcess == nil {
		fmt.Println("react ssr process has not yet boot")
		return ""
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial("0.0.0.0:3024", opts...)
	if err != nil {
		return ""
	}
	defer conn.Close()
	client := NewReactRendererClient(conn)
	response, err := client.Render(context.TODO(), &RenderRequest{
		BundleID: bundleKey,
		JSONData: string(data),
	})
	if err != nil {
		return ""
	}
	return response.StaticContent
}
func reactHydrate(ctx context.Context, bundleKey string, data []byte, doc *htmlDoc) (*htmlDoc, context.Context) {
	innerServerHTML := serverRenderInnerHTML(bundleKey, data)
	if v := ctx.Value(OrbitManifest); v == nil {
		doc.Head = append(doc.Head, fmt.Sprintf(`<script id="orbit_manifest" type="application/json">%s</script>`, data))
		ctx = context.WithValue(ctx, OrbitManifest, true)
	}
	doc.Body = append(doc.Body, fmt.Sprintf(`<script class="orbit_bk" src="/p/%s.js"></script>`, bundleKey))
	copy := doc.Body
	// the doc body is adjusted +1 indices to insert the react frame at the front of the list
	// this is due to react requiring the div id to exist before the necessary javascript is loaded in
	doc.Body = make([]string, len(doc.Body)+1)
	doc.Body[0] = fmt.Sprintf(`<div id="%s_react_frame">%s</div>`, bundleKey, innerServerHTML)
	for i, c := range copy {
		doc.Body[i+1] = c
	}
	return doc, ctx
}
var nodeProcess *os.Process
func Close() error {
	// this is a hack, node process does not get terminated with the nodeProcess.Kill
	// method, but I found that if I use ctrl^c in the terminal, it closes it correctly
	if nodeProcess == nil {
		return nil
	}
	err := syscall.Kill(nodeProcess.Pid, syscall.SIGSTOP)
	nodeProcess = nil
	return err
}
// TODO: phase out the init stuff, prefer this to be autogenerated.
func init() {
	serverStartupTasks = append(serverStartupTasks, StartupTaskReactSSR(bundleDir, wrapDocRender, staticResourceMap, make(map[PageRender]string), *setupDoc()))
}
func StartupTaskReactSSR(
	outDir string,
	pages map[PageRender]*DocumentRenderer,
	staticMap map[PageRender]bool,
	nameMap map[PageRender]string,
	doc htmlDoc,
) func() {
	return func() {
		err := startNodeServer()
		if err != nil {
			panic(err)
		}
		for renderKey := range pages {
			if !staticMap[renderKey] {
				continue
			}
			sr, _ := reactSSR(context.Background(), string(renderKey), []byte("{}"), &doc)
			pathName := string(renderKey)
			if nameMap[renderKey] != "" {
				pathName = nameMap[renderKey]
			}
			path := fmt.Sprintf("%s%c%s", http.Dir(outDir), os.PathSeparator, pathName)
			body := append(pageDependencies[renderKey], sr.Body...)
			so := fmt.Sprintf(`<!doctype html><head>%s</head><body>%s</body></html>`, strings.Join(sr.Head, ""), strings.Join(body, ""))
			err := ioutil.WriteFile(path, []byte(so), 0644)
			if err != nil {
				fmt.Printf("error creating static resource for bundle %s => %s\n", renderKey, err)
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
	// TODO(stab) verify that babel node & grpc are both installed.
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
var staticResourceMap = map[PageRender]bool{
	ExampleTwoPage: true,
	ExamplePage: false,
}
var serverStartupTasks = []func(){}
type RenderFunction func(context.Context, string, []byte, *htmlDoc) (*htmlDoc, context.Context)
var wrapDocRender = map[PageRender]*DocumentRenderer{
	ExampleTwoPage: {fn: reactHydrate, version: "reactHydrate"},
	ExamplePage: {fn: reactHydrate, version: "reactHydrate"},
}

type DocumentRenderer struct {
	fn RenderFunction
	version string
}
var bundleDir string = ".orbit/dist"

var publicDir string = "./public/index.html"
var hotReloadPort int = 0
type PageRender string

const ( 
	// orbit:page .//pages/example2.jsx
	ExampleTwoPage PageRender = "fe9faa2750e8559c8c213c2c25c4ce73"
	// orbit:page .//pages/example.jsx
	ExamplePage PageRender = "496a05464c3f5aa89e1d8bed7afe59d4"
)

var pageDependencies = map[PageRender][]string{
	ExampleTwoPage: {`<script src="/p/fc38086145547d465be97fec2e412a16.js"></script>`,
`<script src="/p/a63649d90703a7b09f22aed8d310be5b.js"></script>`,
},
	ExamplePage: {`<script src="/p/fc38086145547d465be97fec2e412a16.js"></script>`,
`<script src="/p/a63649d90703a7b09f22aed8d310be5b.js"></script>`,
},
}

	
type HydrationCtxKey string

const (
	OrbitManifest HydrationCtxKey = "orbitManifest"
)

type BundleMode int32

const (
	DevBundleMode  BundleMode = 0
	ProdBundleMode BundleMode = 1
)

var CurrentDevMode BundleMode = ProdBundleMode