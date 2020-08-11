package grpcjson_test

import (
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes/any"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/grpc_json"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/grpcjson"

	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

var _ = Describe("GrpcJson", func() {
	var (
		initParams plugins.InitParams

		expectedFilter = []plugins.StagedHttpFilter{
			{
				HttpFilter: &envoyhttp.HttpFilter{
					Name: wellknown.GRPCJSONTranscoder,
					ConfigType: &envoyhttp.HttpFilter_TypedConfig{
						TypedConfig: &any.Any{
							TypeUrl: "type.googleapis.com/envoy.extensions.filters.http.grpc_json_transcoder.v3.GrpcJsonTranscoder",
							Value:   []byte{26, 0, 34, 3, 1, 2, 3}, //don't hardcode
						},
					},
				},
				Stage: plugins.AfterStage(plugins.AuthZStage),
			},
		}
	)

	It("should add filter and translate all fields", func() {
		hl := &v1.HttpListener{
			Options: &v1.HttpListenerOptions{
				GrpcJsonTranscoder: &grpc_json.GrpcJsonTranscoder{
					DescriptorSet:                &grpc_json.GrpcJsonTranscoder_ProtoDescriptorBin{ProtoDescriptorBin: []byte{1, 2, 3}},
					Services:                     nil,
					PrintOptions:                 nil,
					MatchIncomingRequestRoute:    false,
					IgnoredQueryParameters:       nil,
					AutoMapping:                  false,
					IgnoreUnknownQueryParameters: false,
					ConvertGrpcStatus:            false,
					XXX_NoUnkeyedLiteral:         struct{}{},
					XXX_unrecognized:             nil,
					XXX_sizecache:                0,
				},
			},
		}

		p := grpcjson.NewPlugin()
		p.Init(initParams)
		f, err := p.HttpFilters(plugins.Params{}, hl)
		Expect(err).NotTo(HaveOccurred())
		Expect(f).NotTo(BeNil())
		Expect(f).To(Equal(expectedFilter))
	})

})
